package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Logbook struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	StudentID    uuid.UUID  `db:"student_id" json:"studentId"`
	LogDate      time.Time  `db:"log_date" json:"logDate"`
	Activities   string     `db:"activities" json:"activities"`
	Blockers     string     `db:"blockers" json:"blockers"`
	PlanTomorrow string     `db:"plan_tomorrow" json:"planTomorrow"`
	EvidenceURL  *string    `db:"evidence_url" json:"evidenceUrl"`
	Status       string     `db:"status" json:"status"`
	Notes        *string    `db:"notes" json:"notes"`
	SubmittedAt  *time.Time `db:"submitted_at" json:"submittedAt"`
	ReviewedAt   *time.Time `db:"reviewed_at" json:"reviewedAt"`
	ReviewedBy   *uuid.UUID `db:"reviewed_by" json:"reviewedBy"`
}

type LogbookDTO struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	StudentID    uuid.UUID  `db:"student_id" json:"studentId"`
	StudentName  string     `db:"student_name" json:"studentName"`
	MentorID     *uuid.UUID `db:"mentor_id" json:"mentorId"`
	MentorName   *string    `db:"mentor_name" json:"mentorName"`
	LogDate      time.Time  `db:"log_date" json:"logDate"`
	Activities   string     `db:"activities" json:"activities"`
	Blockers     string     `db:"blockers" json:"blockers"`
	PlanTomorrow string     `db:"plan_tomorrow" json:"planTomorrow"`
	EvidenceURL  *string    `db:"evidence_url" json:"evidenceUrl"`
	Status       string     `db:"status" json:"status"`
	Notes        *string    `db:"notes" json:"notes"`
	SubmittedAt  *time.Time `db:"submitted_at" json:"submittedAt"`
	ReviewedAt   *time.Time `db:"reviewed_at" json:"reviewedAt"`
	ReviewedBy   *uuid.UUID `db:"reviewed_by" json:"reviewedBy"`
}

type RequestLogbookFormat struct {
	ID           uuid.UUID `db:"id" json:"id"`
	StudentID    uuid.UUID `json:"-"`
	LogDate      time.Time `db:"log_date" json:"logDate" validate:"required"`
	Activities   string    `db:"activities" json:"activities" validate:"required"`
	Blockers     string    `db:"blockers" json:"blockers" validate:"required"`
	PlanTomorrow string    `db:"plan_tomorrow" json:"planTomorrow" validate:"required"`
	EvidenceURL  *string   `db:"evidence_url" json:"evidenceUrl"`
}

type RequestUpdateLogbookStatusFormat struct {
	Status string    `db:"status" json:"status" validate:"required,oneof=approved rejected"`
	Notes  *string   `db:"notes" json:"notes"`
	UserID uuid.UUID `json:"-"`
}

type RequestLogbookListFormat struct {
	PageSize   int    `json:"pageSize"`
	PageNumber int    `json:"pageNumber"`
	Search     string `json:"search"`
	Status     string `json:"status"`
}

var ColumnMapLogbook = map[string]interface{}{
	"id":           "id",
	"studentId":    "student_id",
	"studentName":  "student_name",
	"mentorId":     "mentor_id",
	"mentorName":   "mentor_name",
	"logDate":      "log_date",
	"activities":   "activities",
	"blockers":     "blockers",
	"planTomorrow": "plan_tomorrow",
	"evidenceUrl":  "evidence_url",
	"status":       "status",
	"submittedAt":  "submitted_at",
	"reviewedAt":   "reviewed_at",
	"reviewedBy":   "reviewed_by",
}

func (l *Logbook) NewLogbookFormat(reqFormat RequestLogbookFormat) (newLogbook Logbook, err error) {
	now := time.Now()
	logbookID := reqFormat.ID
	if logbookID == uuid.Nil {
		logbookID, _ = uuid.NewV4()
	}

	newLogbook = Logbook{
		ID:           logbookID,
		StudentID:    reqFormat.StudentID,
		LogDate:      reqFormat.LogDate,
		Activities:   reqFormat.Activities,
		Blockers:     reqFormat.Blockers,
		PlanTomorrow: reqFormat.PlanTomorrow,
		EvidenceURL:  reqFormat.EvidenceURL,
		Status:       "pending",
		SubmittedAt:  &now,
	}
	return
}

func (l *Logbook) UpdateStatusFormat(reqFormat RequestUpdateLogbookStatusFormat) (newLogbook Logbook, err error) {
	now := time.Now()
	newLogbook = *l
	newLogbook.Status = reqFormat.Status
	newLogbook.Notes = reqFormat.Notes
	newLogbook.ReviewedAt = &now
	newLogbook.ReviewedBy = &reqFormat.UserID
	return
}

func (l *Logbook) UpdateFormat(reqFormat RequestLogbookFormat) (newLogbook Logbook, err error) {
	newLogbook = *l
	newLogbook.Activities = reqFormat.Activities
	newLogbook.Blockers = reqFormat.Blockers
	newLogbook.PlanTomorrow = reqFormat.PlanTomorrow
	newLogbook.EvidenceURL = reqFormat.EvidenceURL
	if !reqFormat.LogDate.IsZero() {
		newLogbook.LogDate = reqFormat.LogDate
	}
	return
}
