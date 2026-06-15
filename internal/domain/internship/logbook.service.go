package internship

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"lms-be/shared/pagination"
	"lms-be/shared/socket"

	"github.com/gofrs/uuid/v5"
)

type LogbookService interface {
	Create(ctx context.Context, req RequestLogbookFormat) (newLogbook Logbook, err error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []LogbookDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []LogbookDTO, err error)
	Update(ctx context.Context, id uuid.UUID, req RequestLogbookFormat) (newLogbook Logbook, err error)
	ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error)
	ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error)
}

type LogbookServiceImpl struct {
	LogbookRepository          LogbookRepository
	MentorAssignmentRepository MentorAssignmentRepository
}

func ProvideLogbookServiceImpl(logbookRepository LogbookRepository, mentorAssignmentRepository MentorAssignmentRepository) *LogbookServiceImpl {
	s := new(LogbookServiceImpl)
	s.LogbookRepository = logbookRepository
	s.MentorAssignmentRepository = mentorAssignmentRepository
	return s
}

func (s *LogbookServiceImpl) Create(ctx context.Context, req RequestLogbookFormat) (newLogbook Logbook, err error) {
	if req.StudentID == uuid.Nil {
		return Logbook{}, errors.New("student id is required")
	}
	if strings.TrimSpace(req.Blockers) == "" {
		return Logbook{}, errors.New("blockers is required")
	}
	if strings.TrimSpace(req.PlanTomorrow) == "" {
		return Logbook{}, errors.New("plan tomorrow is required")
	}
	if strings.TrimSpace(req.Activities) == "" {
		return Logbook{}, errors.New("activities is required")
	}
	if req.LogDate.Time().IsZero() {
		return Logbook{}, errors.New("log date is required")
	}

	// Validate progress status
	status := strings.ToLower(strings.TrimSpace(req.ProgressStatus))
	if status == "" {
		status = "in_progress"
	}
	if status != "in_progress" && status != "done" && status != "blocked" {
		return Logbook{}, errors.New("progress status must be one of: in_progress, done, blocked")
	}
	req.ProgressStatus = status

	newLogbook, _ = newLogbook.NewLogbookFormat(req)
	err = s.LogbookRepository.Create(ctx, &newLogbook)
	if err != nil {
		return Logbook{}, err
	}

	// Get student name for SSE payload
	studentName, err := s.LogbookRepository.GetStudentName(ctx, req.StudentID)
	if err != nil {
		studentName = "Mahasiswa" // Fallback
	}

	message := fmt.Sprintf("%s telah mengumpulkan tugas", studentName)

	// Trigger real-time refresh for HRD Dashboard
	socket.GetInstance().BroadcastToRole("HA03", "monitoring_update", message)

	// Trigger real-time refresh for Mentor Dashboard
	socket.GetInstance().BroadcastToRole("HA04", "refresh_logbooks", message)

	return newLogbook, nil
}

func (s *LogbookServiceImpl) GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []LogbookDTO, err error) {
	if studentID == uuid.Nil {
		return nil, errors.New("student id is required")
	}
	return s.LogbookRepository.GetByStudentID(ctx, studentID)
}

func (s *LogbookServiceImpl) GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []LogbookDTO, err error) {
	if mentorID == uuid.Nil {
		return nil, errors.New("mentor id is required")
	}
	return s.LogbookRepository.GetByMentorID(ctx, mentorID)
}

func (s *LogbookServiceImpl) Update(ctx context.Context, id uuid.UUID, req RequestLogbookFormat) (newLogbook Logbook, err error) {
	if id == uuid.Nil {
		return Logbook{}, errors.New("logbook id is required")
	}
	if req.StudentID == uuid.Nil {
		return Logbook{}, errors.New("student id is required")
	}

	existing, err := s.LogbookRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == uuid.Nil {
		return Logbook{}, errors.New("logbook not found")
	}

	if existing.StudentID != req.StudentID {
		return Logbook{}, errors.New("access denied")
	}

	// Validate progress status if provided
	if req.ProgressStatus != "" {
		status := strings.ToLower(strings.TrimSpace(req.ProgressStatus))
		if status != "in_progress" && status != "done" && status != "blocked" {
			return Logbook{}, errors.New("progress status must be one of: in_progress, done, blocked")
		}
		req.ProgressStatus = status
	}

	newLogbook, _ = existing.UpdateFormat(req)
	err = s.LogbookRepository.Update(ctx, newLogbook)
	if err != nil {
		return Logbook{}, err
	}
	return newLogbook, nil
}

func (s *LogbookServiceImpl) ResolveAllByStudentID(ctx context.Context, studentID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error) {
	if studentID == uuid.Nil {
		return pagination.Response{}, errors.New("student id is required")
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	return s.LogbookRepository.ResolveAllByStudentID(ctx, studentID, req)
}

func (s *LogbookServiceImpl) ResolveAllByMentorID(ctx context.Context, mentorID uuid.UUID, req RequestLogbookListFormat) (data pagination.Response, err error) {
	if mentorID == uuid.Nil {
		return pagination.Response{}, errors.New("mentor id is required")
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}

	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return pagination.Response{}, errors.New("invalid date format")
		}

		now := time.Now()
		todayStr := now.Format("2006-01-02")

		// If requested date is in the future
		if req.Date > todayStr {
			return pagination.Response{
				Items: make([]interface{}, 0),
				Meta:  pagination.CreateMeta(0, req.PageSize, req.PageNumber),
			}, nil
		}

		// If requested date is Saturday or Sunday
		if parsedDate.Weekday() == time.Saturday || parsedDate.Weekday() == time.Sunday {
			return pagination.Response{
				Items: make([]interface{}, 0),
				Meta:  pagination.CreateMeta(0, req.PageSize, req.PageNumber),
			}, nil
		}
	}

	data, err = s.LogbookRepository.ResolveAllByMentorID(ctx, mentorID, req)
	if err != nil {
		return data, err
	}

	// Post-process to handle "late" status
	if req.Date != "" {
		now := time.Now()
		isToday := req.Date == now.Format("2006-01-02")
		isPast := req.Date < now.Format("2006-01-02")
		isPast5PM := now.Hour() >= 17

		for i, item := range data.Items {
			logbook, ok := item.(LogbookDTO)
			if !ok {
				continue
			}
			if logbook.ProgressStatus == "pending" {
				if isPast || (isToday && isPast5PM) {
					logbook.ProgressStatus = "late"
					data.Items[i] = logbook
				}
			}
		}
	}

	return data, nil
}
