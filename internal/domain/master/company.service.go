package master

import (
	"bytes"
	"context"
	"errors"
	"lms-be/configs"
	"lms-be/infras"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

type CompanyService interface {
	Create(ctx context.Context, req RequestCompanyFormat) (newCompany Company, err error)
	Update(ctx context.Context, id int, req RequestCompanyFormat) (newCompany Company, err error)
	GetAll(ctx context.Context, req model.StandardRequest) (data []CompanyDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (Company Company, err error)
	DeleteByID(ctx context.Context, id int, userId uuid.UUID) (err error)
	PreviewFromExcel(ctx context.Context, fileBytes []byte) (result PreviewCompanyResult, err error)
	ImportFromPreview(ctx context.Context, req ImportFromPreviewCompanyRequest, userID uuid.UUID) (result model.ImportResult, err error)
}

type CompanyServiceImpl struct {
	CompanyRepository CompanyRepository
	TxManager         *infras.TxManager
	Config            *configs.Config
}

func ProvideCompanyServiceImpl(repository CompanyRepository, txManager *infras.TxManager) *CompanyServiceImpl {
	s := new(CompanyServiceImpl)
	s.CompanyRepository = repository
	s.TxManager = txManager
	return s
}

func (s *CompanyServiceImpl) GetAll(ctx context.Context, req model.StandardRequest) (data []CompanyDTO, err error) {
	return s.CompanyRepository.GetAll(ctx, req)
}

func (s *CompanyServiceImpl) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	return s.CompanyRepository.ResolveAll(ctx, req)
}

func (s *CompanyServiceImpl) Create(ctx context.Context, req RequestCompanyFormat) (newCompany Company, err error) {
	existName, err := s.CompanyRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:     "name",
		Value:     req.Name,
		ExcludeID: 0,
	})
	if err != nil {
		return Company{}, err
	}

	if existName {
		return Company{}, errors.New("company name already exists")
	}

	newCompanyFormat, _ := newCompany.NewCompanyFormat(req)
	err = s.CompanyRepository.Create(ctx, &newCompanyFormat)
	if err != nil {
		return Company{}, err
	}

	return newCompanyFormat, nil
}

func (s *CompanyServiceImpl) Update(ctx context.Context, id int, req RequestCompanyFormat) (newCompany Company, err error) {
	existing, err := s.CompanyRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == 0 {
		return Company{}, errors.New("company not found")
	}

	existName, err := s.CompanyRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:     "name",
		Value:     req.Name,
		ExcludeID: id,
	})
	if err != nil {
		return Company{}, err
	}

	if existName {
		return Company{}, errors.New("company name already exists")
	}

	req.ID = id
	newCompany, _ = newCompany.NewCompanyFormat(req)
	newCompany.ID = id
	err = s.CompanyRepository.Update(ctx, newCompany)
	if err != nil {
		return Company{}, err
	}

	return newCompany, nil
}

func (s *CompanyServiceImpl) ResolveByID(ctx context.Context, id int) (Company Company, err error) {
	data, err := s.CompanyRepository.ResolveByID(ctx, id)
	if err != nil || data.ID == 0 {
		return Company, errors.New("company not found")
	}

	return data, nil
}

func (s *CompanyServiceImpl) DeleteByID(ctx context.Context, id int, userId uuid.UUID) error {
	company, err := s.CompanyRepository.ResolveByID(ctx, id)
	if err != nil || company.ID == 0 {
		return errors.New("company not found")
	}

	company.SoftDelete(userId)
	err = s.CompanyRepository.Update(ctx, company)
	if err != nil {
		return errors.New("failed to delete company")
	}

	return nil
}

func (r *CompanyServiceImpl) PreviewFromExcel(ctx context.Context, fileBytes []byte) (result PreviewCompanyResult, err error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return result, errors.New("failed to open excel file " + err.Error())
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return result, errors.New("failed to get excel rows " + err.Error())
	}

	if len(rows) < 2 {
		return result, errors.New("file is empty or has no data rows")
	}

	result.TotalRows = len(rows) - 1
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1

		name := strings.TrimSpace(row[0])

		if name == "" {
			result.Errors = append(result.Errors, model.ImportRowError{
				Row:     rowNum,
				Message: "Name is required",
			})
			result.ErrorCount++
			continue
		}

		isRegisteredPartner := false
		if len(row) > 1 && strings.TrimSpace(row[1]) != "" {
			isRegisteredPartnerStr := strings.TrimSpace(row[1])
			if isRegisteredPartnerStr == "true" || isRegisteredPartnerStr == "1" || isRegisteredPartnerStr == "yes" {
				isRegisteredPartner = true
			}
		}

		var picContact *string
		if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
			picContactVal := strings.TrimSpace(row[2])
			picContact = &picContactVal
		}

		exists, _ := r.CompanyRepository.ExistByField(ctx, model.ExistByFieldParams{
			Field:     "name",
			Value:     name,
			ExcludeID: 0,
		})

		result.Data = append(result.Data, PreviewCompanyRow{
			Row:                 rowNum,
			Name:                name,
			IsRegisteredPartner: isRegisteredPartner,
			PICContact:          picContact,
			Exists:              exists,
		})
		result.ValidCount++
	}

	return result, nil
}

func (s *CompanyServiceImpl) ImportFromPreview(ctx context.Context, req ImportFromPreviewCompanyRequest, userID uuid.UUID) (result model.ImportResult, err error) {
	if len(req.Data) == 0 {
		return result, errors.New("no data to import")
	}

	mode := req.Mode
	if mode == "" {
		mode = "insert"
	}

	result.TotalRows = len(req.Data)
	var validData []Company
	now := time.Now()

	for _, row := range req.Data {
		if mode == "insert" && row.Exists {
			result.Errors = append(result.Errors, model.ImportRowError{
				Row:     row.Row,
				Message: "Company with name '" + row.Name + "' already exists (skipped)",
			})
			result.SkipCount++
			continue
		}

		isRegisteredPartner := row.IsRegisteredPartner
		company := Company{
			Name:                row.Name,
			IsRegisteredPartner: isRegisteredPartner,
			PICContact:          row.PICContact,
			CreatedAt:           &now,
			CreatedBy:           &userID,
		}

		validData = append(validData, company)
	}

	if len(validData) > 0 {
		err = s.TxManager.WithTx(ctx, func(tx *sqlx.Tx) error {
			if mode == "upsert" {
				return s.CompanyRepository.UpsertBatchTx(ctx, tx, validData)
			}
			return s.CompanyRepository.CreateBatchTx(ctx, tx, validData)
		})

		if err != nil {
			return result, errors.New("failed to save data: " + err.Error())
		}

		result.SuccessCount = len(validData)
	}

	return result, nil
}
