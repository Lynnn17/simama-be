package internship

import (
	"context"
	"errors"
	"strings"

	"github.com/gofrs/uuid/v5"
	"lms-be/shared/pagination"
)

type LogbookService interface {
	Create(ctx context.Context, req RequestLogbookFormat) (newLogbook Logbook, err error)
	GetByStudentID(ctx context.Context, studentID uuid.UUID) (data []LogbookDTO, err error)
	GetByMentorID(ctx context.Context, mentorID uuid.UUID) (data []LogbookDTO, err error)
	UpdateStatus(ctx context.Context, id uuid.UUID, req RequestUpdateLogbookStatusFormat) (newLogbook Logbook, err error)
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
	if req.LogDate.IsZero() {
		return Logbook{}, errors.New("log date is required")
	}

	newLogbook, _ = newLogbook.NewLogbookFormat(req)
	err = s.LogbookRepository.Create(ctx, &newLogbook)
	if err != nil {
		return Logbook{}, err
	}
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

func (s *LogbookServiceImpl) UpdateStatus(ctx context.Context, id uuid.UUID, req RequestUpdateLogbookStatusFormat) (newLogbook Logbook, err error) {
	if id == uuid.Nil {
		return Logbook{}, errors.New("logbook id is required")
	}
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status != "approved" && status != "rejected" {
		return Logbook{}, errors.New("invalid logbook status")
	}

	existing, err := s.LogbookRepository.ResolveByID(ctx, id)
	if err != nil || existing.ID == uuid.Nil {
		return Logbook{}, errors.New("logbook not found")
	}

	assignment, err := s.MentorAssignmentRepository.ResolveMentorByStudentID(ctx, existing.StudentID)
	if err != nil || assignment.ID == uuid.Nil {
		return Logbook{}, errors.New("mentor assignment not found")
	}

	if assignment.MentorID != req.UserID {
		return Logbook{}, errors.New("access denied")
	}

	req.Status = status
	newLogbook, _ = existing.UpdateStatusFormat(req)
	newLogbook.ID = existing.ID
	err = s.LogbookRepository.UpdateStatus(ctx, newLogbook)
	if err != nil {
		return Logbook{}, err
	}
	return newLogbook, nil
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

	if existing.Status != "pending" {
		return Logbook{}, errors.New("only pending logbook can be updated")
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
	return s.LogbookRepository.ResolveAllByMentorID(ctx, mentorID, req)
}
