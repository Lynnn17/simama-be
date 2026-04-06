package internship

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type RequestTaskFileFormat struct {
	ID         uuid.UUID `db:"id" json:"id"`
	TaskID     uuid.UUID `db:"task_id" json:"taskId" validate:"required"`
	FileURL    string    `db:"file_url" json:"fileUrl" validate:"required"`
	UploadedBy uuid.UUID `json:"-"`
}

func (t *TaskFile) NewTaskFileFormat(reqFormat RequestTaskFileFormat) (newTaskFile TaskFile, err error) {
	now := time.Now()
	fileID := reqFormat.ID
	if fileID == uuid.Nil {
		fileID, _ = uuid.NewV4()
	}

	newTaskFile = TaskFile{
		ID:         fileID,
		TaskID:     reqFormat.TaskID,
		FileURL:    reqFormat.FileURL,
		UploadedBy: reqFormat.UploadedBy,
		CreatedAt:  &now,
	}
	return
}
