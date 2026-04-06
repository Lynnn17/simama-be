package master

import (
	"lms-be/shared/model"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Personnel struct {
	ID                int        `db:"id" json:"id"`
	Name              string     `db:"name" json:"name"`
	PersonnelType     *string    `db:"personnel_type" json:"personnelType"`
	DepartmentID      *int       `db:"department_id" json:"departmentId"`
	CompanyID         *int       `db:"company_id" json:"companyId"`
	PhotoFaceTemplate *int       `db:"photo_face_template" json:"photoFaceTemplate"`
	IsActive          bool       `db:"is_active" json:"isActive"`
	CreatedAt         *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy         *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt         *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy         *uuid.UUID `db:"updated_by" json:"updatedBy"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted         bool       `db:"is_deleted" json:"isDeleted"`
}

type PersonnelDTO struct {
	ID                int     `db:"id" json:"id"`
	Name              string  `db:"name" json:"name"`
	PersonnelType     *string `db:"personnel_type" json:"personnelType"`
	DepartmentID      *int    `db:"department_id" json:"departmentId"`
	DepartmentName    *string `db:"department_name" json:"departmentName"`
	CompanyID         *int    `db:"company_id" json:"companyId"`
	CompanyName       *string `db:"company_name" json:"companyName"`
	PhotoFaceTemplate *int    `db:"photo_face_template" json:"photoFaceTemplate"`
	IsActive          bool    `db:"is_active" json:"isActive"`
}

type RequestPersonnelFormat struct {
	ID                int       `db:"id" json:"id"`
	Name              string    `db:"name" json:"name"`
	PersonnelType     *string   `db:"personnel_type" json:"personnelType"`
	DepartmentID      *int      `db:"department_id" json:"departmentId"`
	CompanyID         *int      `db:"company_id" json:"companyId"`
	PhotoFaceTemplate *int      `db:"photo_face_template" json:"photoFaceTemplate"`
	IsActive          bool      `db:"is_active" json:"isActive"`
	UserID            uuid.UUID `json:"-"`
}

var ColumnMapPersonnel = map[string]interface{}{
	"id":                "p.id",
	"name":              "p.name",
	"personnelType":     "p.personnel_type",
	"departmentId":      "p.department_id",
	"departmentName":    "d.name",
	"companyId":         "p.company_id",
	"companyName":       "c.name",
	"photoFaceTemplate": "p.photo_face_template",
	"isActive":          "p.is_active",
	"createdAt":         "p.created_at",
}

func (p *Personnel) NewPersonnelFormat(reqFormat RequestPersonnelFormat) (newPersonnel Personnel, err error) {
	var now = time.Now()
	isActive := true
	if reqFormat.ID == 0 {
		newPersonnel = Personnel{
			Name:              reqFormat.Name,
			PersonnelType:     reqFormat.PersonnelType,
			DepartmentID:      reqFormat.DepartmentID,
			CompanyID:         reqFormat.CompanyID,
			PhotoFaceTemplate: reqFormat.PhotoFaceTemplate,
			IsActive:          isActive,
			CreatedAt:         &now,
			CreatedBy:         &reqFormat.UserID,
		}
	} else {
		newPersonnel = Personnel{
			ID:                reqFormat.ID,
			Name:              reqFormat.Name,
			PersonnelType:     reqFormat.PersonnelType,
			DepartmentID:      reqFormat.DepartmentID,
			CompanyID:         reqFormat.CompanyID,
			PhotoFaceTemplate: reqFormat.PhotoFaceTemplate,
			IsActive:          reqFormat.IsActive,
			UpdatedAt:         &now,
			UpdatedBy:         &reqFormat.UserID,
		}
	}
	return
}

func (p *Personnel) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	p.IsDeleted = true
	p.UpdatedAt = &now
	p.UpdatedBy = &userID
	p.DeletedAt = &now
}

type PreviewPersonnelRow struct {
	Row               int     `json:"row"`
	Name              string  `json:"name"`
	PersonnelType     *string `json:"personnelType"`
	DepartmentID      *int    `json:"departmentId"`
	DepartmentName    *string `json:"departmentName"`
	CompanyID         *int    `json:"companyId"`
	CompanyName       *string `json:"companyName"`
	PhotoFaceTemplate *int    `json:"photoFaceTemplate"`
	IsActive          bool    `json:"isActive"`
	Exists            bool    `json:"exist"`
}

type PreviewPersonnelResult struct {
	TotalRows  int                    `db:"total_rows" json:"totalRows"`
	ValidCount int                    `db:"valid_count" json:"validCount"`
	ErrorCount int                    `db:"error_count" json:"errorCount"`
	Data       []PreviewPersonnelRow  `json:"data"`
	Errors     []model.ImportRowError `json:"errors,omitempty"`
}

type ImportFromPreviewPersonnelRequest struct {
	Mode string                `json:"mode"`
	Data []PreviewPersonnelRow `json:"data"`
}
