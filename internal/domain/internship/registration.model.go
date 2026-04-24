package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Registration struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     *uuid.UUID `db:"user_id" json:"userId"`
	FullName   string     `db:"full_name" json:"fullName"`
	University string     `db:"university" json:"university"`
	Major      string     `db:"major" json:"major"`
	Semester   string     `db:"semester" json:"semester"`
	Phone      string     `db:"phone" json:"phone"`
	Email      string     `db:"email" json:"email"`
	Period     string     `db:"period" json:"period"`
	CVFilePath string     `db:"cv_file_path" json:"cvFilePath"`
	Status     string     `db:"status" json:"status"`
	AppliedAt  *time.Time `db:"applied_at" json:"appliedAt"`
	ReviewedAt *time.Time `db:"reviewed_at" json:"reviewedAt"`
	ReviewedBy *uuid.UUID `db:"reviewed_by" json:"reviewedBy"`
}

type RegistrationDTO struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     *uuid.UUID `db:"user_id" json:"userId"`
	FullName   string     `db:"full_name" json:"fullName"`
	University string     `db:"university" json:"university"`
	Major      string     `db:"major" json:"major"`
	Semester   string     `db:"semester" json:"semester"`
	Phone      string     `db:"phone" json:"phone"`
	Email      string     `db:"email" json:"email"`
	Period     string     `db:"period" json:"period"`
	CVFilePath string     `db:"cv_file_path" json:"cvFilePath"`
	Status     string     `db:"status" json:"status"`
	AppliedAt  *time.Time `db:"applied_at" json:"appliedAt"`
	ReviewedAt *time.Time `db:"reviewed_at" json:"reviewedAt"`
	ReviewedBy *uuid.UUID `db:"reviewed_by" json:"reviewedBy"`
}

type RequestRegistrationFormat struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	UserID     *uuid.UUID `json:"-"`
	FullName   string     `db:"full_name" json:"fullName" validate:"required"`
	University string     `db:"university" json:"university" validate:"required"`
	Major      string     `db:"major" json:"major" validate:"required"`
	Semester   string     `db:"semester" json:"semester" validate:"required"`
	Phone      string     `db:"phone" json:"phone" validate:"required"`
	Email      string     `db:"email" json:"email" validate:"required,email"`
	Period     string     `db:"period" json:"period" validate:"required"`
	CVFilePath string     `db:"cv_file_path" json:"cvFilePath" validate:"required"`
}

type RequestRegistrationListFormat struct {
	PageNumber int    `json:"pageNumber" validate:"gte=1"`
	PageSize   int    `json:"pageSize" validate:"gte=1"`
	Status     string `json:"status"`
}

type RequestUpdateRegistrationStatusFormat struct {
	Status string     `db:"status" json:"status" validate:"required,oneof=accepted rejected"`
	UserID *uuid.UUID `json:"-"`
}

var ColumnMapRegistration = map[string]interface{}{
	"id":         "id",
	"userId":     "user_id",
	"fullName":   "full_name",
	"university": "university",
	"major":      "major",
	"semester":   "semester",
	"phone":      "phone",
	"email":      "email",
	"period":     "period",
	"cvFilePath": "cv_file_path",
	"status":     "status",
	"appliedAt":  "applied_at",
	"reviewedAt": "reviewed_at",
	"reviewedBy": "reviewed_by",
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
		FullName:   reqFormat.FullName,
		University: reqFormat.University,
		Major:      reqFormat.Major,
		Semester:   reqFormat.Semester,
		Phone:      reqFormat.Phone,
		Email:      reqFormat.Email,
		Period:     reqFormat.Period,
		CVFilePath: reqFormat.CVFilePath,
		Status:     "pending",
		AppliedAt:  &now,
		ReviewedAt: nil,
		ReviewedBy: nil,
	}
	return
}

func (r *Registration) UpdateStatusFormat(reqFormat RequestUpdateRegistrationStatusFormat) (newRegistration Registration, err error) {
	now := time.Now()
	newRegistration = *r
	newRegistration.Status = reqFormat.Status
	newRegistration.ReviewedAt = &now
	newRegistration.ReviewedBy = reqFormat.UserID
	return
}
