package master

import (
	"lms-be/shared/model"
	"time"

	"github.com/gofrs/uuid/v5"
)

// AcademicYear represents the m_academic_year table
type AcademicYear struct {
	ID            int        `db:"id" json:"id"`
	Code          string     `db:"code" json:"code"`
	Name          string     `db:"name" json:"name"`
	InstitutionID int        `db:"institution_id" json:"institutionId"`
	SemesterType  *string    `db:"semester_type" json:"semesterType"`
	IsActive      *bool      `db:"is_active" json:"isActive"`
	CreatedAt     *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy     *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy     *uuid.UUID `db:"updated_by" json:"updatedBy"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted     bool       `db:"is_deleted" json:"isDeleted"`
}

// AcademicYearDTO represents the m_academic_year table with joined data
type AcademicYearDTO struct {
	ID              int        `db:"id" json:"id"`
	Code            string     `db:"code" json:"code"`
	Name            string     `db:"name" json:"name"`
	InstitutionID   int        `db:"institution_id" json:"institutionId"`
	SemesterType    *string    `db:"semester_type" json:"semesterType"`
	IsActive        *bool      `db:"is_active" json:"isActive"`
	CreatedAt       *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy       *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy       *uuid.UUID `db:"updated_by" json:"updatedBy"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted       bool       `db:"is_deleted" json:"isDeleted"`
	InstitutionName *string    `db:"institution_name" json:"institutionName"`
}

// RequestAcademicYearFormat is swagger request format
type RequestAcademicYearFormat struct {
	ID            int       `db:"id" json:"id"`
	Code          string    `db:"code" json:"code"`
	Name          string    `db:"name" json:"name"`
	InstitutionID int       `db:"institution_id" json:"institutionId"`
	SemesterType  *string   `db:"semester_type" json:"semesterType"`
	IsActive      *bool     `db:"is_active" json:"isActive"`
	UserID        uuid.UUID `json:"-"`
}

// ColumnMapAcademicYear is alias from JSON to DB for sorting in FE
var ColumnMapAcademicYear = map[string]interface{}{
	"id":              "ay.id",
	"code":            "ay.code",
	"name":            "ay.name",
	"institutionId":   "ay.institution_id",
	"institutionName": "i.name",
	"semesterType":    "ay.semester_type",
	"isActive":        "ay.is_active",
	"createdAt":       "ay.created_at",
	"createdBy":       "ay.created_by",
	"updatedAt":       "ay.updated_at",
	"updatedBy":       "ay.updated_by",
	"deletedAt":       "ay.deleted_at",
	"isDeleted":       "ay.is_deleted",
}

// NewAcademicYearFormat creates or updates an academic year
func (ay *AcademicYear) NewAcademicYearFormat(reqFormat RequestAcademicYearFormat) (newAcademicYear AcademicYear, err error) {
	var now = time.Now()
	if reqFormat.ID == 0 {
		newAcademicYear = AcademicYear{
			Code:          reqFormat.Code,
			Name:          reqFormat.Name,
			InstitutionID: reqFormat.InstitutionID,
			SemesterType:  reqFormat.SemesterType,
			IsActive:      reqFormat.IsActive,
			CreatedAt:     &now,
			CreatedBy:     &reqFormat.UserID,
		}
	} else {
		newAcademicYear = AcademicYear{
			ID:            reqFormat.ID,
			Code:          reqFormat.Code,
			Name:          reqFormat.Name,
			InstitutionID: reqFormat.InstitutionID,
			SemesterType:  reqFormat.SemesterType,
			IsActive:      reqFormat.IsActive,
			UpdatedAt:     &now,
			UpdatedBy:     &reqFormat.UserID,
		}
	}
	return
}

// SoftDelete marks an AcademicYear as deleted
func (ay *AcademicYear) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	ay.IsDeleted = true
	ay.UpdatedAt = &now
	ay.UpdatedBy = &userID
	ay.DeletedAt = &now
}

// PreviewAcademicYearResult represents the preview result before import
type PreviewAcademicYearResult struct {
	TotalRows  int                      `json:"totalRows"`
	ValidCount int                      `json:"validCount"`
	ErrorCount int                      `json:"errorCount"`
	Data       []PreviewAcademicYearRow `json:"data"`
	Errors     []model.ImportRowError   `json:"errors,omitempty"`
}

// PreviewAcademicYearRow represents a single row in preview
type PreviewAcademicYearRow struct {
	Row          int     `json:"row"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	SemesterType *string `json:"semesterType,omitempty"`
	IsActive     bool    `json:"isActive"`
	Exists       bool    `json:"exists"` // true if code already exists in DB
}

// ImportFromPreviewRequest represents the request to import from preview data
type ImportFromPreviewRequest struct {
	InstitutionID int                      `json:"institutionId"`
	Mode          string                   `json:"mode"` // insert or upsert
	Data          []PreviewAcademicYearRow `json:"data"`
}
