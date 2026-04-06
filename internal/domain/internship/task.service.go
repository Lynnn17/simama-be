package internship

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"lms-be/infras"
)

type TaskService interface {
	Create(ctx context.Context, req RequestTaskFormat) (newTask Task, err error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []TaskDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []TaskDTO, err error)
	SubmitTaskFile(ctx context.Context, req RequestSubmitTaskFileFormat) (newTaskFile TaskFile, err error)
	GradeTask(ctx context.Context, id uuid.UUID, req RequestGradeTaskFormat) (newTask Task, err error)
}

type TaskServiceImpl struct {
	TaskRepository             TaskRepository
	TaskFileRepository         TaskFileRepository
	MentorAssignmentRepository MentorAssignmentRepository
	TxManager                  *infras.TxManager
}

func ProvideTaskServiceImpl(taskRepository TaskRepository, taskFileRepository TaskFileRepository, mentorAssignmentRepository MentorAssignmentRepository, txManager *infras.TxManager) *TaskServiceImpl {
	s := new(TaskServiceImpl)
	s.TaskRepository = taskRepository
	s.TaskFileRepository = taskFileRepository
	s.MentorAssignmentRepository = mentorAssignmentRepository
	s.TxManager = txManager
	return s
}

func (s *TaskServiceImpl) Create(ctx context.Context, req RequestTaskFormat) (newTask Task, err error) {
	if req.MentorID == uuid.Nil || req.StudentID == uuid.Nil {
		return Task{}, errors.New("mentor id and student id are required")
	}
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Description) == "" || req.Deadline.IsZero() {
		return Task{}, errors.New("title, description, and deadline are required")
	}
	if exist, err := s.MentorAssignmentRepository.ExistByStudentID(ctx, req.StudentID); err != nil || !exist {
		if err != nil {
			return Task{}, err
		}
		return Task{}, errors.New("mentor assignment not found")
	}
	if req.MentorID != (uuid.UUID{}) {
		assignment, err := s.MentorAssignmentRepository.ResolveMentorByStudentID(ctx, req.StudentID)
		if err != nil || assignment.ID == uuid.Nil {
			return Task{}, errors.New("mentor assignment not found")
		}
		if assignment.MentorID != req.MentorID {
			return Task{}, errors.New("access denied")
		}
	}
	newTask, _ = newTask.NewTaskFormat(req)
	err = s.TaskRepository.Create(ctx, &newTask)
	if err != nil {
		return Task{}, err
	}
	return newTask, nil
}

func (s *TaskServiceImpl) GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []TaskDTO, err error) {
	if studentID == uuid.Nil {
		return nil, errors.New("student id is required")
	}
	return s.TaskRepository.GetByStudentID(ctx, studentID)
}

func (s *TaskServiceImpl) GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []TaskDTO, err error) {
	if mentorID == uuid.Nil {
		return nil, errors.New("mentor id is required")
	}
	return s.TaskRepository.GetByMentorID(ctx, mentorID)
}

func (s *TaskServiceImpl) SubmitTaskFile(ctx context.Context, req RequestSubmitTaskFileFormat) (newTaskFile TaskFile, err error) {
	if req.TaskID == uuid.Nil || req.UploadedBy == uuid.Nil || strings.TrimSpace(req.FileURL) == "" {
		return TaskFile{}, errors.New("task id, uploaded by, and file url are required")
	}

	task, err := s.TaskRepository.ResolveByID(ctx, req.TaskID)
	if err != nil || task.ID == uuid.Nil {
		return TaskFile{}, errors.New("task not found")
	}
	if task.StudentID != req.UploadedBy {
		return TaskFile{}, errors.New("access denied")
	}

	newTaskFile, _ = newTaskFile.NewTaskFileFormat(RequestTaskFileFormat{
		TaskID:     req.TaskID,
		FileURL:    req.FileURL,
		UploadedBy: req.UploadedBy,
	})

	err = s.TxManager.WithTx(ctx, func(tx *sqlx.Tx) error {
		fileRepo := s.TaskFileRepository.(*TaskFileRepositoryPostgreSQL)
		taskRepo := s.TaskRepository.(*TaskRepositoryPostgreSQL)
		if err := fileRepo.CreateTx(ctx, tx, &newTaskFile); err != nil {
			return err
		}
		task.Status = "submitted"
		if err := taskRepo.UpdateSubmittedTx(ctx, tx, task); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TaskFile{}, err
	}
	return newTaskFile, nil
}

func (s *TaskServiceImpl) GradeTask(ctx context.Context, id uuid.UUID, req RequestGradeTaskFormat) (newTask Task, err error) {
	if id == uuid.Nil {
		return Task{}, errors.New("task id is required")
	}
	if req.UserID == uuid.Nil {
		return Task{}, errors.New("user id is required")
	}
	task, err := s.TaskRepository.ResolveByID(ctx, id)
	if err != nil || task.ID == uuid.Nil {
		return Task{}, errors.New("task not found")
	}
	assignment, err := s.MentorAssignmentRepository.ResolveMentorByStudentID(ctx, task.StudentID)
	if err != nil || assignment.ID == uuid.Nil {
		return Task{}, errors.New("mentor assignment not found")
	}
	if assignment.MentorID != req.UserID {
		return Task{}, errors.New("access denied")
	}
	newTask, _ = task.UpdateGradeFormat(req)
	newTask.ID = task.ID
	err = s.TaskRepository.UpdateGrade(ctx, newTask)
	if err != nil {
		return Task{}, err
	}
	return newTask, nil
}
