package master

import (
	"bytes"
	"context"
	"errors"
	"lms-be/configs"
	"lms-be/infras"
	"lms-be/shared/model"
	"lms-be/shared/pagination"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/xuri/excelize/v2"
)

type PersonnelService interface {
	Create(ctx context.Context, req RequestPersonnelFormat) (newPersonnel Personnel, err error)
	Update(ctx context.Context, id int, req RequestPersonnelFormat) (newPersonnel Personnel, err error)
	GetAll(ctx context.Context, req model.StandardRequest) (data []PersonnelDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (Personnel Personnel, err error)
	DeleteByID(ctx context.Context, id int, userId uuid.UUID) (err error)
	PreviewFromExcel(ctx context.Context, fileBytes []byte) (result PreviewPersonnelResult, err error)
	ImportFromPreview(ctx context.Context, req ImportFromPreviewPersonnelRequest, userID uuid.UUID) (result model.ImportResult, err error)
}

type PersonnelServiceImpl struct {
	PersonnelRepository PersonnelRepository
	TxManager          *infras.TxManager
	Config             *configs.Config
}

func ProvidePersonnelServiceImpl(repository PersonnelRepository, txManager *infras.TxManager) *PersonnelServiceImpl {
	s := new(PersonnelServiceImpl)
	s.PersonnelRepository = repository
	s.TxManager = txManager
	return s
}

func (s *PersonnelServiceImpl) GetAll(ctx context.Context, req model.StandardRequest) (data []PersonnelDTO, err error) {
	return s.PersonnelRepository.GetAll(ctx, req)
}

func (s *PersonnelServiceImpl) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	return s.PersonnelRepository.ResolveAll(ctx, req)
}

func (s *PersonnelServiceImpl) Create(ctx context.Context, req RequestPersonnelFormat) (newPersonnel Personnel, err error) {
	// Validate company_id if provided
	if req.CompanyID != nil && *req.CompanyID > 0 {
		existCompany, err := s.PersonnelRepository.ExistCompany(ctx, *req.CompanyID)
		if err != nil {
			return Personnel{}, err
		}
		if !existCompany {
			return Personnel{}, errors.New("company not found")
		}
	}

	// Validate department_id if provided
	if req.DepartmentID != nil && *req.DepartmentID > 0 {
		existDepartment, err := s.PersonnelRepository.ExistDepartment(ctx, *req.DepartmentID)
		if err != nil {
			return Personnel{}, err
		}
		if !existDepartment {
			return Personnel{}, errors.New("department not found")
		}
	}

	newPersonnelFormat, _ := newPersonnel.NewPersonnelFormat(req)
	err = s.PersonnelRepository.Create(ctx, &newPersonnelFormat)
	if err != nil {
		return Personnel{}, err
	}

	return newPersonnelFormat, nil
}

func (s *PersonnelServiceImpl) Update(ctx context.Context, id int, req RequestPersonnelFormat) (newPersonnel Personnel, err error) {
	existing, err := s.PersonnelRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == 0 {
		return Personnel{}, errors.New("personnel not found")
	}

	// Validate company_id if provided
	if req.CompanyID != nil && *req.CompanyID > 0 {
		existCompany, err := s.PersonnelRepository.ExistCompany(ctx, *req.CompanyID)
		if err != nil {
			return Personnel{}, err
		}
		if !existCompany {
			return Personnel{}, errors.New("company not found")
		}
	}

	// Validate department_id if provided
	if req.DepartmentID != nil && *req.DepartmentID > 0 {
		existDepartment, err := s.PersonnelRepository.ExistDepartment(ctx, *req.DepartmentID)
		if err != nil {
			return Personnel{}, err
		}
		if !existDepartment {
			return Personnel{}, errors.New("department not found")
		}
	}

	req.ID = id
	newPersonnel, _ = newPersonnel.NewPersonnelFormat(req)
	newPersonnel.ID = id
	err = s.PersonnelRepository.Update(ctx, newPersonnel)
	if err != nil {
		return Personnel{}, err
	}

	return newPersonnel, nil
}

func (s *PersonnelServiceImpl) ResolveByID(ctx context.Context, id int) (Personnel Personnel, err error) {
	data, err := s.PersonnelRepository.ResolveByID(ctx, id)
	if err != nil || data.ID == 0 {
		return Personnel, errors.New("personnel not found")
	}

	return data, nil
}

func (s *PersonnelServiceImpl) DeleteByID(ctx context.Context, id int, userId uuid.UUID) error {
	personnel, err := s.PersonnelRepository.ResolveByID(ctx, id)
	if err != nil || personnel.ID == 0 {
		return errors.New("personnel not found")
	}

	personnel.SoftDelete(userId)
	err = s.PersonnelRepository.Update(ctx, personnel)
	if err != nil {
		return errors.New("failed to delete personnel")
	}

	return nil
}

func (r *PersonnelServiceImpl) PreviewFromExcel(ctx context.Context, fileBytes []byte) (result PreviewPersonnelResult, err error) {
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
			result.ErrorCount++
			result.Errors = append(result.Errors, model.ImportRowError{
				Row:     rowNum,
				Message: "name is required",
			})
			continue
		}

		var personnelType *string
		if len(row) > 1 && strings.TrimSpace(row[1]) != "" {
			pt := strings.TrimSpace(row[1])
			personnelType = &pt
		}

		var departmentID *int
		if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
			deptID, err := strconv.Atoi(strings.TrimSpace(row[2]))
			if err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, model.ImportRowError{
					Row:     rowNum,
					Message: "invalid department_id format",
				})
				continue
			}
			departmentID = &deptID
		}

		var companyID *int
		if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
			compID, err := strconv.Atoi(strings.TrimSpace(row[3]))
			if err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, model.ImportRowError{
					Row:     rowNum,
					Message: "invalid company_id format",
				})
				continue
			}
			companyID = &compID
		}

		var photoFaceTemplate *int
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			pft, err := strconv.Atoi(strings.TrimSpace(row[4]))
			if err != nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, model.ImportRowError{
					Row:     rowNum,
					Message: "invalid photo_face_template format",
				})
				continue
			}
			photoFaceTemplate = &pft
		}

		isActive := true
		if len(row) > 5 && strings.TrimSpace(row[5]) != "" {
			isActiveBool := strings.ToLower(strings.TrimSpace(row[5]))
			if isActiveBool == "false" || isActiveBool == "0" || isActiveBool == "no" {
				isActive = false
			}
		}

		var departmentName *string
		if departmentID != nil && *departmentID > 0 {
			deptName, _ := r.PersonnelRepository.GetDepartmentNameByID(ctx, *departmentID)
			if deptName != "" {
				departmentName = &deptName
			}
		}

		var companyName *string
		if companyID != nil && *companyID > 0 {
			compName, _ := r.PersonnelRepository.GetCompanyNameByID(ctx, *companyID)
			if compName != "" {
				companyName = &compName
			}
		}

		previewRow := PreviewPersonnelRow{
			Row:               rowNum,
			Name:              name,
			PersonnelType:      personnelType,
			DepartmentID:      departmentID,
			DepartmentName:    departmentName,
			CompanyID:         companyID,
			CompanyName:       companyName,
			PhotoFaceTemplate: photoFaceTemplate,
			IsActive:          isActive,
			Exists:            false,
		}

		// Check if already exists
		existPersonnel, _ := r.PersonnelRepository.ExistByField(ctx, model.ExistByFieldParams{
			Field:     "name",
			Value:     name,
			ExcludeID: 0,
		})
		if existPersonnel {
			previewRow.Exists = true
		}

		result.ValidCount++
		result.Data = append(result.Data, previewRow)
	}

	return result, nil
}

func (s *PersonnelServiceImpl) ImportFromPreview(ctx context.Context, req ImportFromPreviewPersonnelRequest, userID uuid.UUID) (result model.ImportResult, err error) {
	if len(req.Data) == 0 {
		return model.ImportResult{}, errors.New("no data to import")
	}

	mode := req.Mode
	if mode == "" {
		mode = "insert"
	}

	result.TotalRows = len(req.Data)
	var validData []Personnel
	now := time.Now()

	for _, row := range req.Data {
		// Validate company_id if provided
		if row.CompanyID != nil && *row.CompanyID > 0 {
			existCompany, err := s.PersonnelRepository.ExistCompany(ctx, *row.CompanyID)
			if err != nil || !existCompany {
				result.ErrorCount++
				continue
			}
		}

		// Validate department_id if provided
		if row.DepartmentID != nil && *row.DepartmentID > 0 {
			existDepartment, err := s.PersonnelRepository.ExistDepartment(ctx, *row.DepartmentID)
			if err != nil || !existDepartment {
				result.ErrorCount++
				continue
			}
		}

		personnel := Personnel{
			Name:              row.Name,
			PersonnelType:      row.PersonnelType,
			DepartmentID:      row.DepartmentID,
			CompanyID:         row.CompanyID,
			PhotoFaceTemplate: row.PhotoFaceTemplate,
			IsActive:          row.IsActive,
			CreatedAt:         &now,
			CreatedBy:         &userID,
		}
		validData = append(validData, personnel)
	}

	if len(validData) > 0 {
		err = s.TxManager.WithTx(ctx, func(tx *sqlx.Tx) error {
			if mode == "upsert" {
				return s.PersonnelRepository.UpsertBatchTx(ctx, tx, validData)
			}
			return s.PersonnelRepository.CreateBatchTx(ctx, tx, validData)
		})

		if err == nil {
			result.SuccessCount = len(validData)
		}
	}

	return result, nil
}
