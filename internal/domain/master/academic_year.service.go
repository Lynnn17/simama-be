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

type AcademicYearService interface {
	Create(ctx context.Context, req RequestAcademicYearFormat) (newAcademicYear AcademicYear, err error)
	Update(ctx context.Context, id int, req RequestAcademicYearFormat) (newAcademicYear AcademicYear, err error)
	GetAll(ctx context.Context, req model.StandardRequest) (data []AcademicYearDTO, err error)
	ResolveAll(ctx context.Context, request model.StandardRequest) (academicYears pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (data AcademicYear, err error)
	Delete(ctx context.Context, id int, userId uuid.UUID) error
	PreviewFromExcel(ctx context.Context, fileBytes []byte, institutionID int) (result PreviewAcademicYearResult, err error)
	ImportFromPreview(ctx context.Context, req ImportFromPreviewRequest, userID uuid.UUID) (result model.ImportResult, err error)
}

type AcademicYearServiceImpl struct {
	AcademicYearRepository AcademicYearRepository
	TxManager              *infras.TxManager
	Config                 *configs.Config
}

func ProvideAcademicYearServiceImpl(repository AcademicYearRepository, txManager *infras.TxManager) *AcademicYearServiceImpl {
	s := new(AcademicYearServiceImpl)
	s.AcademicYearRepository = repository
	s.TxManager = txManager
	return s
}

func (s *AcademicYearServiceImpl) GetAll(ctx context.Context, req model.StandardRequest) (data []AcademicYearDTO, err error) {
	return s.AcademicYearRepository.GetAll(ctx, req)
}

func (s *AcademicYearServiceImpl) ResolveAll(ctx context.Context, request model.StandardRequest) (academicYears pagination.Response, err error) {
	return s.AcademicYearRepository.ResolveAll(ctx, request)
}

func (s *AcademicYearServiceImpl) Create(ctx context.Context, req RequestAcademicYearFormat) (newAcademicYear AcademicYear, err error) {
	existName, err := s.AcademicYearRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:         "name",
		Value:         req.Name,
		ExcludeID:     0,
		InstitutionID: req.InstitutionID,
	})
	if err != nil {
		return AcademicYear{}, err
	}
	if existName {
		return AcademicYear{}, errors.New("academic year with this name already exists")
	}

	existCode, err := s.AcademicYearRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:         "code",
		Value:         req.Code,
		ExcludeID:     0,
		InstitutionID: req.InstitutionID,
	})
	if err != nil {
		return AcademicYear{}, err
	}
	if existCode {
		return AcademicYear{}, errors.New("academic year with this code already exists")
	}

	newAcademicYear, _ = newAcademicYear.NewAcademicYearFormat(req)
	err = s.AcademicYearRepository.Create(ctx, &newAcademicYear)
	if err != nil {
		return AcademicYear{}, err
	}
	return newAcademicYear, nil
}

func (s *AcademicYearServiceImpl) Update(ctx context.Context, id int, req RequestAcademicYearFormat) (newAcademicYear AcademicYear, err error) {
	existing, err := s.AcademicYearRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == 0 {
		return AcademicYear{}, errors.New("academic year not found")
	}

	existName, err := s.AcademicYearRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:         "name",
		Value:         req.Name,
		ExcludeID:     id,
		InstitutionID: req.InstitutionID,
	})
	if err != nil {
		return AcademicYear{}, err
	}
	if existName {
		return AcademicYear{}, errors.New("academic year name already exists")
	}

	existCode, err := s.AcademicYearRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:         "code",
		Value:         req.Code,
		ExcludeID:     id,
		InstitutionID: req.InstitutionID,
	})
	if err != nil {
		return AcademicYear{}, err
	}
	if existCode {
		return AcademicYear{}, errors.New("academic year code already exists")
	}

	req.ID = id
	newAcademicYear, _ = newAcademicYear.NewAcademicYearFormat(req)
	newAcademicYear.ID = id
	err = s.AcademicYearRepository.Update(ctx, newAcademicYear)
	if err != nil {
		return AcademicYear{}, err
	}
	return newAcademicYear, nil
}

func (s *AcademicYearServiceImpl) ResolveByID(ctx context.Context, id int) (data AcademicYear, err error) {
	data, err = s.AcademicYearRepository.ResolveByID(ctx, id)
	if err != nil || data.ID == 0 {
		return AcademicYear{}, errors.New("academic year not found")
	}
	return data, nil
}

func (s *AcademicYearServiceImpl) Delete(ctx context.Context, id int, userId uuid.UUID) error {
	academicYear, err := s.AcademicYearRepository.ResolveByID(ctx, id)
	if err != nil || academicYear.ID == 0 {
		return errors.New("academic year not found")
	}

	academicYear.SoftDelete(userId)
	err = s.AcademicYearRepository.Update(ctx, academicYear)
	if err != nil {
		return errors.New("failed to delete academic year")
	}
	return nil
}

func (s *AcademicYearServiceImpl) PreviewFromExcel(ctx context.Context, fileBytes []byte, institutionID int) (result PreviewAcademicYearResult, err error) {
	f, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return result, errors.New("failed to open excel file: " + err.Error())
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return result, errors.New("failed to read sheet: " + err.Error())
	}

	if len(rows) < 2 {
		return result, errors.New("file is empty or has no data rows")
	}

	result.TotalRows = len(rows) - 1
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1

		code := strings.TrimSpace(row[0])
		name := ""
		if len(row) > 1 {
			name = strings.TrimSpace(row[1])
		}

		// Validate required fields
		if code == "" || name == "" {
			result.Errors = append(result.Errors, model.ImportRowError{
				Row:     rowNum,
				Message: "Code and Name are required",
			})
			result.ErrorCount++
			continue
		}

		// Parse optional fields
		var semesterType *string
		if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
			st := strings.TrimSpace(row[2])
			semesterType = &st
		}

		isActive := false
		if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
			isActiveStr := strings.ToLower(strings.TrimSpace(row[3]))
			if isActiveStr == "true" || isActiveStr == "1" {
				isActive = true
			}
		}

		// Check if code exists
		exists, _ := s.AcademicYearRepository.ExistByField(ctx, model.ExistByFieldParams{
			Field:         "code",
			Value:         code,
			ExcludeID:     0,
			InstitutionID: institutionID,
		})

		result.Data = append(result.Data, PreviewAcademicYearRow{
			Row:          rowNum,
			Code:         code,
			Name:         name,
			SemesterType: semesterType,
			IsActive:     isActive,
			Exists:       exists,
		})
		result.ValidCount++
	}

	return result, nil
}

func (s *AcademicYearServiceImpl) ImportFromPreview(ctx context.Context, req ImportFromPreviewRequest, userID uuid.UUID) (result model.ImportResult, err error) {
	if len(req.Data) == 0 {
		return result, errors.New("no data to import")
	}

	// Default mode is "insert"
	mode := req.Mode
	if mode == "" {
		mode = "insert"
	}

	result.TotalRows = len(req.Data)
	var validData []AcademicYear
	now := time.Now()

	for _, row := range req.Data {
		// For INSERT mode: skip existing data
		if mode == "insert" && row.Exists {
			result.Errors = append(result.Errors, model.ImportRowError{
				Row:     row.Row,
				Message: "Academic year with code '" + row.Code + "' already exists (skipped)",
			})
			result.SkipCount++
			continue
		}

		isActive := row.IsActive
		academicYear := AcademicYear{
			Code:          row.Code,
			Name:          row.Name,
			InstitutionID: req.InstitutionID,
			SemesterType:  row.SemesterType,
			IsActive:      &isActive,
			CreatedAt:     &now,
			CreatedBy:     &userID,
		}

		validData = append(validData, academicYear)
	}

	// Execute in transaction
	if len(validData) > 0 {
		err = s.TxManager.WithTx(ctx, func(tx *sqlx.Tx) error {
			if mode == "upsert" {
				return s.AcademicYearRepository.UpsertBatchTx(ctx, tx, validData)
			}
			return s.AcademicYearRepository.CreateBatchTx(ctx, tx, validData)
		})

		if err != nil {
			return result, errors.New("failed to save data: " + err.Error())
		}
		result.SuccessCount = len(validData)
	}

	return result, nil
}
