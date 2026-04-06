package internship

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
)

type MentorAssignmentService interface {
	Create(ctx context.Context, req RequestMentorAssignmentFormat) (newMentorAssignment MentorAssignment, err error)
	GetStudentsByMentorID(ctx context.Context, mentorID uuid.UUID) (data []MentorAssignmentDTO, err error)
	ResolveMentorByStudentID(ctx context.Context, studentID uuid.UUID) (data MentorAssignmentDTO, err error)
}

type MentorAssignmentServiceImpl struct {
	MentorAssignmentRepository MentorAssignmentRepository
}

func ProvideMentorAssignmentServiceImpl(repository MentorAssignmentRepository) *MentorAssignmentServiceImpl {
	s := new(MentorAssignmentServiceImpl)
	s.MentorAssignmentRepository = repository
	return s
}

func (s *MentorAssignmentServiceImpl) Create(ctx context.Context, req RequestMentorAssignmentFormat) (newMentorAssignment MentorAssignment, err error) {
	if req.MentorID == uuid.Nil || req.StudentID == uuid.Nil {
		return MentorAssignment{}, errors.New("mentor id and student id are required")
	}

	exist, err := s.MentorAssignmentRepository.ExistByStudentID(ctx, req.StudentID)
	if err != nil {
		return MentorAssignment{}, err
	}
	if exist {
		return MentorAssignment{}, errors.New("student already has a mentor assignment")
	}

	newMentorAssignment, _ = newMentorAssignment.NewMentorAssignmentFormat(req)
	err = s.MentorAssignmentRepository.Create(ctx, &newMentorAssignment)
	if err != nil {
		return MentorAssignment{}, err
	}
	return newMentorAssignment, nil
}

func (s *MentorAssignmentServiceImpl) GetStudentsByMentorID(ctx context.Context, mentorID uuid.UUID) (data []MentorAssignmentDTO, err error) {
	if mentorID == uuid.Nil {
		return nil, errors.New("mentor id is required")
	}
	return s.MentorAssignmentRepository.GetStudentsByMentorID(ctx, mentorID)
}

func (s *MentorAssignmentServiceImpl) ResolveMentorByStudentID(ctx context.Context, studentID uuid.UUID) (data MentorAssignmentDTO, err error) {
	if studentID == uuid.Nil {
		return MentorAssignmentDTO{}, errors.New("student id is required")
	}
	data, err = s.MentorAssignmentRepository.ResolveMentorByStudentID(ctx, studentID)
	if err != nil {
		return MentorAssignmentDTO{}, err
	}
	return data, nil
}
