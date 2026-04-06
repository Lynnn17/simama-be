package master

import (
	"bytes"
	"context"
	"encoding/json"
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

type DepartmentService interface {
	Create(ctx context.Context, req RequestDepartmentFormat) (newDepartment Department, err error)
	Update(ctx context.Context, id int, req RequestDepartmentFormat) (newDepartment Department, err error)
	GetAll(ctx context.Context, req model.StandardRequest) (data []DepartmentDTO, err error)
	ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error)
	ResolveByID(ctx context.Context, id int) (Department Department, err error)
	DeleteByID(ctx context.Context, id int, userId uuid.UUID) (err error)
	PreviewFromExcel(ctx context.Context, fileBytes []byte) (result PreviewDepartmentResult, err error)
	ImportFromPreview(ctx context.Context, req ImportFromPreviewDepartmentRequest, userID uuid.UUID) (result model.ImportResult, err error)
}

type DepartmentServiceImpl struct {
	DepartmentRepository DepartmentRepository
	TxManager            *infras.TxManager
	Config               *configs.Config
}

func ProvideDepartmentServiceImpl(repository DepartmentRepository, txManager *infras.TxManager) *DepartmentServiceImpl {
	s := new(DepartmentServiceImpl)
	s.DepartmentRepository = repository
	s.TxManager = txManager
	return s
}

func (s *DepartmentServiceImpl) GetAll(ctx context.Context, req model.StandardRequest) (data []DepartmentDTO, err error) {
	return s.DepartmentRepository.GetAll(ctx, req)
}

func (s *DepartmentServiceImpl) ResolveAll(ctx context.Context, req model.StandardRequest) (data pagination.Response, err error) {
	return s.DepartmentRepository.ResolveAll(ctx, req)
}

func (s *DepartmentServiceImpl) Create(ctx context.Context, req RequestDepartmentFormat) (newDepartment Department, err error) {
	existName, err := s.DepartmentRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:     "name",
		Value:     req.Name,
		ExcludeID: 0,
	})
	if err != nil {
		return Department{}, err
	}

	if existName {
		return Department{}, errors.New("department name already exists")
	}

	newDepartmentFormat, _ := newDepartment.NewDepartmentFormat(req)
	err = s.DepartmentRepository.Create(ctx, &newDepartmentFormat)
	if err != nil {
		return Department{}, err
	}

	return newDepartmentFormat, nil
}

func (s *DepartmentServiceImpl) Update(ctx context.Context, id int, req RequestDepartmentFormat) (newDepartment Department, err error) {
	existing, err := s.DepartmentRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == 0 {
		return Department{}, errors.New("department not found")
	}

	existName, err := s.DepartmentRepository.ExistByField(ctx, model.ExistByFieldParams{
		Field:     "name",
		Value:     req.Name,
		ExcludeID: id,
	})
	if err != nil {
		return Department{}, err
	}

	if existName {
		return Department{}, errors.New("department name already exists")
	}

	req.ID = id
	newDepartment, _ = newDepartment.NewDepartmentFormat(req)
	newDepartment.ID = id
	err = s.DepartmentRepository.Update(ctx, newDepartment)
	if err != nil {
		return Department{}, err
	}

	return newDepartment, nil
}

func (s *DepartmentServiceImpl) ResolveByID(ctx context.Context, id int) (Department Department, err error) {
	data, err := s.DepartmentRepository.ResolveByID(ctx, id)
	if err != nil || data.ID == 0 {
		return Department, errors.New("department not found")
	}

	return data, nil
}

func (s *DepartmentServiceImpl) DeleteByID(ctx context.Context, id int, userId uuid.UUID) error {
	department, err := s.DepartmentRepository.ResolveByID(ctx, id)
	if err != nil || department.ID == 0 {
		return errors.New("department not found")
	}

	department.SoftDelete(userId)
	err = s.DepartmentRepository.Update(ctx, department)
	if err != nil {
		return errors.New("failed to delete department")
	}

	return nil
}

func (r *DepartmentServiceImpl) PreviewFromExcel(ctx context.Context, fileBytes []byte) (result PreviewDepartmentResult, err error) {
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

		var mapCoordinateJSON *MapCoordinateJSON
		if len(row) > 1 && strings.TrimSpace(row[1]) != "" {
			mapCoordinateStr := strings.TrimSpace(row[1])
			var mapData MapCoordinateJSON
			err := json.Unmarshal([]byte(mapCoordinateStr), &mapData)
			if err == nil {
				mapCoordinateJSON = &mapData
			}
		}

		exists, _ := r.DepartmentRepository.ExistByField(ctx, model.ExistByFieldParams{
			Field:     "name",
			Value:     name,
			ExcludeID: 0,
		})

		result.Data = append(result.Data, PreviewDepartmentRow{
			Row:               rowNum,
			Name:              name,
			MapCoordinateJSON: mapCoordinateJSON,
			Exists:            exists,
		})
		result.ValidCount++
	}

	return result, nil
}

func (s *DepartmentServiceImpl) ImportFromPreview(ctx context.Context, req ImportFromPreviewDepartmentRequest, userID uuid.UUID) (result model.ImportResult, err error) {
	if len(req.Data) == 0 {
		return result, errors.New("no data to import")
	}

	mode := req.Mode
	if mode == "" {
		mode = "insert"
	}

	result.TotalRows = len(req.Data)
	var validData []Department
	now := time.Now()

	for _, row := range req.Data {
		if mode == "insert" && row.Exists {
			result.Errors = append(result.Errors, model.ImportRowError{
				Row:     row.Row,
				Message: "Department with name '" + row.Name + "' already exists (skipped)",
			})
			result.SkipCount++
			continue
		}

		department := Department{
			Name:              row.Name,
			MapCoordinateJSON: row.MapCoordinateJSON,
			CreatedAt:         &now,
			CreatedBy:         &userID,
		}

		validData = append(validData, department)
	}

	if len(validData) > 0 {
		err = s.TxManager.WithTx(ctx, func(tx *sqlx.Tx) error {
			if mode == "upsert" {
				return s.DepartmentRepository.UpsertBatchTx(ctx, tx, validData)
			}
			return s.DepartmentRepository.CreateBatchTx(ctx, tx, validData)
		})

		if err != nil {
			return result, errors.New("failed to save data: " + err.Error())
		}

		result.SuccessCount = len(validData)
	}

	return result, nil
}
