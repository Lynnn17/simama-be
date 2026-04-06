package internship

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
)

type RegistrationService interface {
	Create(ctx context.Context, req RequestRegistrationFormat) (newRegistration Registration, err error)
	GetAll(ctx context.Context) (data []RegistrationDTO, err error)
}

type RegistrationServiceImpl struct {
	RegistrationRepository RegistrationRepository
}

func ProvideRegistrationServiceImpl(repository RegistrationRepository) *RegistrationServiceImpl {
	s := new(RegistrationServiceImpl)
	s.RegistrationRepository = repository
	return s
}

func (s *RegistrationServiceImpl) Create(ctx context.Context, req RequestRegistrationFormat) (newRegistration Registration, err error) {
	if req.UserID == uuid.Nil {
		return Registration{}, errors.New("user id is required")
	}

	exist, err := s.RegistrationRepository.ExistByUserID(ctx, req.UserID)
	if err != nil {
		return Registration{}, err
	}
	if exist {
		return Registration{}, errors.New("registration already exists")
	}

	newRegistration, _ = newRegistration.NewRegistrationFormat(req)
	err = s.RegistrationRepository.Create(ctx, &newRegistration)
	if err != nil {
		return Registration{}, err
	}
	return newRegistration, nil
}

func (s *RegistrationServiceImpl) GetAll(ctx context.Context) (data []RegistrationDTO, err error) {
	return s.RegistrationRepository.GetAll(ctx)
}
