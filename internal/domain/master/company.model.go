package master

import (
	"lms-be/shared/model"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Company struct {
	ID                  int        `db:"id" json:"id"`
	Name                string     `db:"name" json:"name"`
	IsRegisteredPartner bool       `db:"is_registered_partner" json:"isRegisteredPartner"`
	PICContact          *string    `db:"pic_contact" json:"picContact"`
	CreatedAt           *time.Time `db:"created_at" json:"createdAt"`
	CreatedBy           *uuid.UUID `db:"created_by" json:"createdBy"`
	UpdatedAt           *time.Time `db:"updated_at" json:"updatedAt"`
	UpdatedBy           *uuid.UUID `db:"updated_by" json:"updatedBy"`
	DeletedAt           *time.Time `db:"deleted_at" json:"deletedAt"`
	IsDeleted           bool       `db:"is_deleted" json:"isDeleted"`
}

type CompanyDTO struct {
	ID                  int     `db:"id" json:"id"`
	Name                string  `db:"name" json:"name"`
	IsRegisteredPartner bool    `db:"is_registered_partner" json:"isRegisteredPartner"`
	PICContact          *string `db:"pic_contact" json:"picContact"`
}

type RequestCompanyFormat struct {
	ID                  int       `db:"id" json:"id"`
	Name                string    `db:"name" json:"name"`
	IsRegisteredPartner bool      `db:"is_registered_partner" json:"isRegisteredPartner"`
	PICContact          *string   `db:"pic_contact" json:"picContact"`
	UserID              uuid.UUID `json:"-"`
}

var ColumnMapCompany = map[string]interface{}{
	"id":                  "id",
	"name":                "name",
	"isRegisteredPartner": "is_registered_partner",
	"picContact":          "pic_contact",
	"createdAt":           "created_at",
}

func (c *Company) NewCompanyFormat(reqFormat RequestCompanyFormat) (newCompany Company, err error) {
	var now = time.Now()
	if reqFormat.ID == 0 {
		newCompany = Company{
			Name:                reqFormat.Name,
			IsRegisteredPartner: reqFormat.IsRegisteredPartner,
			PICContact:          reqFormat.PICContact,
			CreatedAt:           &now,
			CreatedBy:           &reqFormat.UserID,
		}
	} else {
		newCompany = Company{
			ID:                  reqFormat.ID,
			Name:                reqFormat.Name,
			IsRegisteredPartner: reqFormat.IsRegisteredPartner,
			PICContact:          reqFormat.PICContact,
			UpdatedAt:           &now,
			UpdatedBy:           &reqFormat.UserID,
		}
	}
	return
}

func (c *Company) SoftDelete(userID uuid.UUID) {
	now := time.Now()
	c.IsDeleted = true
	c.UpdatedAt = &now
	c.UpdatedBy = &userID
	c.DeletedAt = &now
}

type PreviewCompanyRow struct {
	Row                 int     `json:"row"`
	Name                string  `json:"name"`
	IsRegisteredPartner bool    `json:"isRegisteredPartner"`
	PICContact          *string `json:"picContact"`
	Exists              bool    `json:"exist"`
}

type PreviewCompanyResult struct {
	TotalRows  int                    `db:"total_rows" json:"totalRows"`
	ValidCount int                    `db:"valid_count" json:"validCount"`
	ErrorCount int                    `db:"error_count" json:"errorCount"`
	Data       []PreviewCompanyRow    `json:"data"`
	Errors     []model.ImportRowError `json:"errors,omitempty"`
}

type ImportFromPreviewCompanyRequest struct {
	Mode string              `json:"mode"`
	Data []PreviewCompanyRow `json:"data"`
}
