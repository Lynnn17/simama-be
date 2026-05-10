package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"lms-be/shared/model"
)

type Logbook struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	StudentID      uuid.UUID  `db:"student_id" json:"studentId"`
	LogDate        time.Time  `db:"log_date" json:"logDate"`
	Activities     string     `db:"activities" json:"activities"`
	Blockers       string     `db:"blockers" json:"blockers"`
	PlanTomorrow   string     `db:"plan_tomorrow" json:"planTomorrow"`
	EvidenceURL    *string    `db:"evidence_url" json:"evidenceUrl"`
	ProgressStatus string     `db:"progress_status" json:"progressStatus"`
	SubmittedAt    *time.Time `db:"submitted_at" json:"submittedAt"`
}

type LogbookDTO struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	StudentID      uuid.UUID  `db:"student_id" json:"studentId"`
	StudentName    string     `db:"student_name" json:"studentName"`
	MentorID       *uuid.UUID `db:"mentor_id" json:"mentorId"`
	MentorName     *string    `db:"mentor_name" json:"mentorName"`
	LogDate        time.Time  `db:"log_date" json:"logDate"`
	Activities     string     `db:"activities" json:"activities"`
	Blockers       string     `db:"blockers" json:"blockers"`
	PlanTomorrow   string     `db:"plan_tomorrow" json:"planTomorrow"`
	EvidenceURL    *string    `db:"evidence_url" json:"evidenceUrl"`
	ProgressStatus string     `db:"progress_status" json:"progressStatus"`
	SubmittedAt    *time.Time `db:"submitted_at" json:"submittedAt"`
}

type RequestLogbookFormat struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	StudentID      uuid.UUID      `json:"-"`
	LogDate        model.JSONDate `db:"log_date" json:"logDate" validate:"required"`
	Activities     string         `db:"activities" json:"activities" validate:"required"`
	Blockers       string         `db:"blockers" json:"blockers" validate:"required"`
	PlanTomorrow   string         `db:"plan_tomorrow" json:"planTomorrow" validate:"required"`
	EvidenceURL    string         `db:"evidence_url" json:"evidenceUrl" validate:"required,url"`
	ProgressStatus string         `db:"progress_status" json:"progressStatus" validate:"required,oneof=in_progress done blocked"`
}

type RequestLogbookListFormat struct {
	PageSize       int    `json:"pageSize"`
	PageNumber     int    `json:"pageNumber"`
	Search         string `json:"search"`
	ProgressStatus string `json:"progressStatus"`
	Date           string `json:"date"`
}

var ColumnMapLogbook = map[string]interface{}{
	"id":             "id",
	"studentId":      "student_id",
	"studentName":    "student_name",
	"mentorId":       "mentor_id",
	"mentorName":     "mentor_name",
	"logDate":        "log_date",
	"activities":     "activities",
	"blockers":       "blockers",
	"planTomorrow":   "plan_tomorrow",
	"evidenceUrl":    "evidence_url",
	"progressStatus": "progress_status",
	"submittedAt":    "submitted_at",
}

func (l *Logbook) NewLogbookFormat(reqFormat RequestLogbookFormat) (newLogbook Logbook, err error) {
	now := time.Now()
	logbookID := reqFormat.ID
	if logbookID == uuid.Nil {
		logbookID, _ = uuid.NewV4()
	}

	progressStatus := reqFormat.ProgressStatus
	if progressStatus == "" {
		progressStatus = "in_progress"
	}

	newLogbook = Logbook{
		ID:             logbookID,
		StudentID:      reqFormat.StudentID,
		LogDate:        reqFormat.LogDate.Time(),
		Activities:     reqFormat.Activities,
		Blockers:       reqFormat.Blockers,
		PlanTomorrow:   reqFormat.PlanTomorrow,
		EvidenceURL:    &reqFormat.EvidenceURL,
		ProgressStatus: progressStatus,
		SubmittedAt:    &now,
	}
	return
}

func (l *Logbook) UpdateFormat(reqFormat RequestLogbookFormat) (newLogbook Logbook, err error) {
	newLogbook = *l
	newLogbook.Activities = reqFormat.Activities
	newLogbook.Blockers = reqFormat.Blockers
	newLogbook.PlanTomorrow = reqFormat.PlanTomorrow
	newLogbook.EvidenceURL = &reqFormat.EvidenceURL
	if reqFormat.ProgressStatus != "" {
		newLogbook.ProgressStatus = reqFormat.ProgressStatus
	}
	if !reqFormat.LogDate.Time().IsZero() {
		newLogbook.LogDate = reqFormat.LogDate.Time()
	}
	return
}
