package master

import (
	"database/sql/driver"
	"encoding/json"
	"lms-be/shared/model"
	"time"

	"github.com/gofrs/uuid/v5"
)

type MapCoordinateJSON map[string]interface{}

func (m MapCoordinateJSON) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *MapCoordinateJSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &m)
}

type Department struct {
	ID                int                `db:"id" json:"id"`
	Name              string             `db:"name" json:"name"`
	MapCoordinateJSON *MapCoordinateJSON `db:"map_coordinate_json" json:"mapCoordinateJson"`
	CreatedAt         *time.Time         `db:"created_at" json:"createdAt"`
	CreatedBy         *uuid.UUID         `db:"created_by" json:"createdBy"`
	UpdatedAt         *time.Time         `db:"updated_at" json:"updatedAt"`
	UpdatedBy         *uuid.UUID         `db:"updated_by" json:"updatedBy"`
	DeletedAt         *time.Time         `db:"deleted_at" json:"deletedAt"`
	IsDeleted         bool               `db:"is_deleted" json:"isDeleted"`
}

type DepartmentDTO struct {
	ID                int                `db:"id" json:"id"`
	Name              string             `db:"name" json:"name"`
	MapCoordinateJSON *MapCoordinateJSON `db:"map_coordinate_json" json:"mapCoordinateJson"`
}

type RequestDepartmentFormat struct {
	ID                int                `db:"id" json:"id"`
	Name              string             `db:"name" json:"name"`
	MapCoordinateJSON *MapCoordinateJSON `db:"map_coordinate_json" json:"mapCoordinateJson"`
	UserID            uuid.UUID          `json:"-"`
}

var ColumnMapDepartment = map[string]interface{}{
	"id":                "id",
	"name":              "name",
	"mapCoordinateJson": "map_coordinate_json",
	"createdAt":         "created_at",
}

func (d *Department) NewDepartmentFormat(reqFormat RequestDepartmentFormat) (newDepartment Department, err error) {
	var now = time.Now()
	if reqFormat.ID == 0 {
		newDepartment = Department{
			Name:              reqFormat.Name,
			MapCoordinateJSON: reqFormat.MapCoordinateJSON,
			CreatedAt:         &now,
			CreatedBy:         &reqFormat.UserID,
		}
	} else {
		newDepartment = Department{
			ID:                reqFormat.ID,
			Name:              reqFormat.Name,
			MapCoordinateJSON: reqFormat.MapCoordinateJSON,
			UpdatedAt:         &now,
			UpdatedBy:         &reqFormat.UserID,
		}
	}
	return
}

func (d *Department) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	d.IsDeleted = true
	d.UpdatedAt = &now
	d.UpdatedBy = &userID
	d.DeletedAt = &now
}

type PreviewDepartmentRow struct {
	Row               int                `json:"row"`
	Name              string             `json:"name"`
	MapCoordinateJSON *MapCoordinateJSON `json:"mapCoordinateJson"`
	Exists            bool               `json:"exist"`
}

type PreviewDepartmentResult struct {
	TotalRows  int                    `db:"total_rows" json:"totalRows"`
	ValidCount int                    `db:"valid_count" json:"validCount"`
	ErrorCount int                    `db:"error_count" json:"errorCount"`
	Data       []PreviewDepartmentRow `json:"data"`
	Errors     []model.ImportRowError `json:"errors,omitempty"`
}

type ImportFromPreviewDepartmentRequest struct {
	Mode string                 `json:"mode"`
	Data []PreviewDepartmentRow `json:"data"`
}
