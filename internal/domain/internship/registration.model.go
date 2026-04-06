package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Registration struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"userId"`
	University string     `db:"university" json:"university"`
	Major      string     `db:"major" json:"major"`
	CVFilePath string     `db:"cv_file_path" json:"cvFilePath"`
	Status     string     `db:"status" json:"status"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
}

type RegistrationDTO struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     uuid.UUID  `db:"user_id" json:"userId"`
	University string     `db:"university" json:"university"`
	Major      string     `db:"major" json:"major"`
	CVFilePath string     `db:"cv_file_path" json:"cvFilePath"`
	Status     string     `db:"status" json:"status"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
}

type RequestRegistrationFormat struct {
	ID         uuid.UUID `db:"id" json:"id"`
	UserID     uuid.UUID `json:"-"`
	University string    `db:"university" json:"university" validate:"required"`
	Major      string    `db:"major" json:"major" validate:"required"`
	CVFilePath string    `db:"cv_file_path" json:"cvFilePath" validate:"required"`
}

var ColumnMapRegistration = map[string]interface{}{
	"id":         "id",
	"userId":     "user_id",
	"university": "university",
	"major":      "major",
	"cvFilePath": "cv_file_path",
	"status":     "status",
	"createdAt":  "created_at",
}

func (r *Registration) NewRegistrationFormat(reqFormat RequestRegistrationFormat) (newRegistration Registration, err error) {
	now := time.Now()
	regID := reqFormat.ID
	if regID == uuid.Nil {
		regID, _ = uuid.NewV4()
	}

	newRegistration = Registration{
		ID:         regID,
		UserID:     reqFormat.UserID,
		University: reqFormat.University,
		Major:      reqFormat.Major,
		CVFilePath: reqFormat.CVFilePath,
		Status:     "pending",
		CreatedAt:  &now,
	}
	return
}
