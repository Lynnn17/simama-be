package internship

import (
	"github.com/gofrs/uuid/v5"
)

type MentorAssignment struct {
	ID        uuid.UUID `db:"id" json:"id"`
	MentorID  uuid.UUID `db:"mentor_id" json:"mentorId"`
	StudentID uuid.UUID `db:"student_id" json:"studentId"`
}

type MentorAssignmentDTO struct {
	ID          uuid.UUID `db:"id" json:"id"`
	MentorID    uuid.UUID `db:"mentor_id" json:"mentorId"`
	MentorName  string    `db:"mentor_name" json:"mentorName"`
	StudentID   uuid.UUID `db:"student_id" json:"studentId"`
	StudentName string    `db:"student_name" json:"studentName"`
}

type RequestMentorAssignmentFormat struct {
	ID        uuid.UUID `db:"id" json:"id"`
	MentorID  uuid.UUID `db:"mentor_id" json:"mentorId" validate:"required"`
	StudentID uuid.UUID `db:"student_id" json:"studentId" validate:"required"`
}

var ColumnMapMentorAssignment = map[string]interface{}{
	"id":          "id",
	"mentorId":    "mentor_id",
	"mentorName":  "mentor_name",
	"studentId":   "student_id",
	"studentName": "student_name",
}

func (m *MentorAssignment) NewMentorAssignmentFormat(reqFormat RequestMentorAssignmentFormat) (newMentorAssignment MentorAssignment, err error) {
	assignmentID := reqFormat.ID
	if assignmentID == uuid.Nil {
		assignmentID, _ = uuid.NewV4()
	}

	newMentorAssignment = MentorAssignment{
		ID:        assignmentID,
		MentorID:  reqFormat.MentorID,
		StudentID: reqFormat.StudentID,
	}
	return
}
