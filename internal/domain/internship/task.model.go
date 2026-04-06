package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Task struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	MentorID    uuid.UUID  `db:"mentor_id" json:"mentorId"`
	StudentID   uuid.UUID  `db:"student_id" json:"studentId"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	Deadline    time.Time  `db:"deadline" json:"deadline"`
	Status      string     `db:"status" json:"status"`
	Grade       *int       `db:"grade" json:"grade"`
	Feedback    *string    `db:"feedback" json:"feedback"`
	CreatedAt   *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt"`
}

type TaskDTO struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	MentorID      uuid.UUID  `db:"mentor_id" json:"mentorId"`
	MentorName    string     `db:"mentor_name" json:"mentorName"`
	StudentID     uuid.UUID  `db:"student_id" json:"studentId"`
	StudentName   string     `db:"student_name" json:"studentName"`
	Title         string     `db:"title" json:"title"`
	Description   string     `db:"description" json:"description"`
	Deadline      time.Time  `db:"deadline" json:"deadline"`
	Status        string     `db:"status" json:"status"`
	Grade         *int       `db:"grade" json:"grade"`
	Feedback      *string    `db:"feedback" json:"feedback"`
	CreatedAt     *time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updatedAt"`
	LatestFileID  *uuid.UUID `db:"latest_file_id" json:"latestFileId"`
	LatestFileURL *string    `db:"latest_file_url" json:"latestFileUrl"`
}

type TaskFile struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	TaskID     uuid.UUID  `db:"task_id" json:"taskId"`
	FileURL    string     `db:"file_url" json:"fileUrl"`
	UploadedBy uuid.UUID  `db:"uploaded_by" json:"uploadedBy"`
	CreatedAt  *time.Time `db:"created_at" json:"createdAt"`
}

type RequestTaskFormat struct {
	ID          uuid.UUID `db:"id" json:"id"`
	MentorID    uuid.UUID `db:"mentor_id" json:"mentorId" validate:"required"`
	StudentID   uuid.UUID `db:"student_id" json:"studentId" validate:"required"`
	Title       string    `db:"title" json:"title" validate:"required"`
	Description string    `db:"description" json:"description" validate:"required"`
	Deadline    time.Time `db:"deadline" json:"deadline" validate:"required"`
}

type RequestSubmitTaskFileFormat struct {
	TaskID     uuid.UUID `db:"task_id" json:"taskId" validate:"required"`
	FileURL    string    `db:"file_url" json:"fileUrl" validate:"required"`
	UploadedBy uuid.UUID `json:"-"`
}

type RequestGradeTaskFormat struct {
	Grade    int       `db:"grade" json:"grade" validate:"required"`
	Feedback *string   `db:"feedback" json:"feedback"`
	UserID   uuid.UUID `json:"-"`
}

var ColumnMapTask = map[string]interface{}{
	"id":          "id",
	"mentorId":    "mentor_id",
	"mentorName":  "mentor_name",
	"studentId":   "student_id",
	"studentName": "student_name",
	"title":       "title",
	"deadline":    "deadline",
	"status":      "status",
	"createdAt":   "created_at",
}

func (t *Task) NewTaskFormat(reqFormat RequestTaskFormat) (newTask Task, err error) {
	now := time.Now()
	taskID := reqFormat.ID
	if taskID == uuid.Nil {
		taskID, _ = uuid.NewV4()
	}

	newTask = Task{
		ID:          taskID,
		MentorID:    reqFormat.MentorID,
		StudentID:   reqFormat.StudentID,
		Title:       reqFormat.Title,
		Description: reqFormat.Description,
		Deadline:    reqFormat.Deadline,
		Status:      "assigned",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	return
}

func (t *Task) UpdateGradeFormat(reqFormat RequestGradeTaskFormat) (newTask Task, err error) {
	now := time.Now()
	newTask = *t
	newTask.Status = "graded"
	newTask.Grade = &reqFormat.Grade
	newTask.Feedback = reqFormat.Feedback
	newTask.UpdatedAt = &now
	return
}

func (t *Task) MarkSubmitted() (newTask Task, err error) {
	now := time.Now()
	newTask = *t
	newTask.Status = "submitted"
	newTask.UpdatedAt = &now
	return
}
