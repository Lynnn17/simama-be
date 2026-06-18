package internship

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"

	"lms-be/infras"
	"lms-be/shared/pagination"
	"lms-be/shared/socket"
)

type TaskService interface {
	Create(ctx context.Context, req RequestTaskFormat) (newTask Task, err error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []TaskDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []TaskDTO, err error)
	ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error)
	ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error)
	SubmitTaskFile(ctx context.Context, req RequestSubmitTaskFileFormat) (newTaskFile TaskFile, err error)
	GradeTask(ctx context.Context, id uuid.UUID, req RequestGradeTaskFormat) (newTask Task, err error)
	Update(ctx context.Context, id uuid.UUID, req RequestTaskFormat) (newTask Task, err error)
	GetFilesByTaskID(ctx context.Context, taskID uuid.UUID) (data []TaskFile, err error)
	CreateZipFromTaskFiles(ctx context.Context, taskID uuid.UUID, baseDir string) (data []byte, filename string, err error)
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

	// Notify Student about new task
	socket.GetInstance().SendToUser(req.StudentID.String(), "new_notification", map[string]interface{}{
		"title":   "Tugas Baru",
		"message": "Anda mendapatkan tugas baru: " + req.Title + ". Deadline: " + req.Deadline.Format("02-01-2006 15:04") + ".",
		"type":    "task",
	})

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

		if err := fileRepo.DeleteByTaskIDTx(ctx, tx, req.TaskID); err != nil {
			return err
		}

		if err := fileRepo.CreateTx(ctx, tx, &newTaskFile); err != nil {
			return err
		}
		task.Status = "submitted"
		if req.SubmissionURL != "" {
			task.SubmissionURL = &req.SubmissionURL
		}
		if err := taskRepo.UpdateSubmittedTx(ctx, tx, task); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return TaskFile{}, err
	}

	// Notify Mentor about submission
	studentName := "Mahasiswa"
	if assignment, err := s.MentorAssignmentRepository.ResolveMentorByStudentID(ctx, task.StudentID); err == nil && assignment.StudentName != "" {
		studentName = assignment.StudentName
	}

	socket.GetInstance().SendToUser(task.MentorID.String(), "new_notification", map[string]interface{}{
		"title":   "Tugas Dikumpulkan",
		"message": "Mahasiswa " + studentName + " telah mengumpulkan tugas: " + task.Title,
		"type":    "task_submission",
	})

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

	// 1. Injeksi Timestamp untuk Feedback
	if req.Feedback != nil && *req.Feedback != "" {
		waktuSekarang := time.Now().Format("2006-01-02 15:04")
		prefix := ""
		if req.Status == "revision_needed" {
			prefix = " - Revisi"
		} else if req.Status == "graded" {
			prefix = " - Nilai"
		}
		formattedFeedback := "\n[" + waktuSekarang + prefix + "] " + *req.Feedback
		req.Feedback = &formattedFeedback
	}

	newTask, _ = task.UpdateGradeFormat(req)
	newTask.ID = task.ID

	// 2. Implementasi json.Marshal untuk Kolom Nilai (grade)
	if req.Grade != nil {
		gradeBytes, err := json.Marshal(req.Grade)
		if err == nil {
			gradeStr := string(gradeBytes)
			newTask.Grade = &gradeStr
		}
	}

	err = s.TaskRepository.UpdateGrade(ctx, newTask)
	if err != nil {
		return Task{}, err
	}

	// Notify Student about grade/feedback
	if newTask.Status == "graded" {
		socket.GetInstance().SendToUser(task.StudentID.String(), "new_notification", map[string]interface{}{
			"message": "Tugas '" + task.Title + "' Anda telah dinilai. Lihat feedback.",
			"type":    "task_graded",
		})
	} else if newTask.Status == "revision_needed" {
		socket.GetInstance().SendToUser(task.StudentID.String(), "new_notification", map[string]interface{}{
			"title":   "Tugas Perlu Direvisi",
			"message": "Tugas '" + task.Title + "' Anda perlu direvisi. Baca catatan dari Mentor.",
			"type":    "task_revision",
		})
	}

	return newTask, nil
}

func (s *TaskServiceImpl) ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error) {
	if studentID == uuid.Nil {
		return pagination.Response{}, errors.New("student id is required")
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	return s.TaskRepository.ResolveAllByStudentID(ctx, studentID, req)
}

func (s *TaskServiceImpl) ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestTaskListFormat) (data pagination.Response, err error) {
	if mentorID == uuid.Nil {
		return pagination.Response{}, errors.New("mentor id is required")
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	return s.TaskRepository.ResolveAllByMentorID(ctx, mentorID, req)
}

func (s *TaskServiceImpl) Update(ctx context.Context, id uuid.UUID, req RequestTaskFormat) (newTask Task, err error) {
	if id == uuid.Nil {
		return Task{}, errors.New("task id is required")
	}
	if req.MentorID == uuid.Nil {
		return Task{}, errors.New("mentor id is required")
	}

	existing, err := s.TaskRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == uuid.Nil {
		return Task{}, errors.New("task not found")
	}

	if existing.MentorID != req.MentorID {
		return Task{}, errors.New("access denied")
	}

	newTask, _ = existing.UpdateFormat(req)
	err = s.TaskRepository.Update(ctx, newTask)
	if err != nil {
		return Task{}, err
	}
	return newTask, nil
}

func (s *TaskServiceImpl) GetFilesByTaskID(ctx context.Context, taskID uuid.UUID) (data []TaskFile, err error) {
	if taskID == uuid.Nil {
		return nil, errors.New("task id is required")
	}
	return s.TaskFileRepository.ResolveByTaskID(ctx, taskID)
}

func (s *TaskServiceImpl) CreateZipFromTaskFiles(ctx context.Context, taskID uuid.UUID, baseDir string) (data []byte, filename string, err error) {
	files, err := s.TaskFileRepository.ResolveByTaskID(ctx, taskID)
	if err != nil {
		return nil, "", err
	}
	if len(files) == 0 {
		return nil, "", errors.New("no files found for this task")
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	for _, file := range files {
		f, err := os.Open(filepath.Join(baseDir, file.FileURL))
		if err != nil {
			continue // skip if file not found
		}

		zf, err := w.Create(filepath.Base(file.FileURL))
		if err != nil {
			f.Close()
			return nil, "", err
		}
		_, err = io.Copy(zf, f)
		f.Close()
		if err != nil {
			return nil, "", err
		}
	}

	err = w.Close()
	if err != nil {
		return nil, "", err
	}

	return buf.Bytes(), fmt.Sprintf("task-%s-files.zip", taskID.String()), nil
}
