package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type MentorAssignment struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	MentorID   uuid.UUID  `db:"mentor_id" json:"mentorId"`
	StudentID  uuid.UUID  `db:"student_id" json:"studentId"`
	AssignedBy uuid.UUID  `db:"assigned_by" json:"assignedBy"`
	AssignedAt *time.Time `db:"assigned_at" json:"assignedAt"`
	IsActive   bool       `db:"is_active" json:"isActive"`
}

type MentorAssignmentDTO struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	MentorID    uuid.UUID  `db:"mentor_id" json:"mentorId"`
	MentorName  string     `db:"mentor_name" json:"mentorName"`
	StudentID   uuid.UUID  `db:"student_id" json:"studentId"`
	StudentName string     `db:"student_name" json:"studentName"`
	AssignedBy  uuid.UUID  `db:"assigned_by" json:"assignedBy"`
	AssignedAt  *time.Time `db:"assigned_at" json:"assignedAt"`
	IsActive    bool       `db:"is_active" json:"isActive"`
}

type RequestMentorAssignmentFormat struct {
	ID         uuid.UUID `db:"id" json:"id"`
	MentorID   uuid.UUID `db:"mentor_id" json:"mentorId" validate:"required"`
	StudentID  uuid.UUID `db:"student_id" json:"studentId" validate:"required"`
	AssignedBy uuid.UUID `json:"-"`
}

type RequestMentorAssignmentListFormat struct {
	PageSize   int        `json:"pageSize"`
	PageNumber int        `json:"pageNumber"`
	MentorID   *uuid.UUID `json:"mentorId"`
	StudentID  *uuid.UUID `json:"studentId"`
	IsActive   *bool      `json:"isActive"`
	Search     string     `json:"search"`
}

var ColumnMapMentorAssignment = map[string]interface{}{
	"id":          "id",
	"mentorId":    "mentor_id",
	"mentorName":  "mentor_name",
	"studentId":   "student_id",
	"studentName": "student_name",
	"assignedBy":  "assigned_by",
	"assignedAt":  "assigned_at",
	"isActive":    "is_active",
}

func (m *MentorAssignment) NewMentorAssignmentFormat(reqFormat RequestMentorAssignmentFormat) (newMentorAssignment MentorAssignment, err error) {
	assignmentID := reqFormat.ID
	if assignmentID == uuid.Nil {
		assignmentID, _ = uuid.NewV4()
	}

	newMentorAssignment = MentorAssignment{
		ID:         assignmentID,
		MentorID:   reqFormat.MentorID,
		StudentID:  reqFormat.StudentID,
		AssignedBy: reqFormat.AssignedBy,
		IsActive:   true,
	}
	now := time.Now()
	newMentorAssignment.AssignedAt = &now
	return
}
