package auth

import (
	"lms-be/shared/model"
	"lms-be/shared/pagination"

	"github.com/gofrs/uuid"
)

type LogSystemService interface {
	CreateLogSystem(reqFormat RequestLogSystemFormat, userId uuid.UUID, ipAddress string, userAgent string) (logSystem LogSystem, error error)
	ResolveAll(request model.StandardRequest) (orders pagination.Response, err error)
}

type LogSystemServiceImpl struct {
	LogSystemRepository LogSystemRepository
}

func ProvideLogSystemServiceImpl(LogSystemRepository LogSystemRepository) *LogSystemServiceImpl {
	s := new(LogSystemServiceImpl)
	s.LogSystemRepository = LogSystemRepository
	return s
}

func (s *LogSystemServiceImpl) CreateLogSystem(reqFormat RequestLogSystemFormat, userId uuid.UUID, ipAddress string, userAgent string) (newLogSystem LogSystem, err error) {
	if err != nil {
		return LogSystem{}, err
	}
	newLogSystem, _ = newLogSystem.NewLogSystemFormat(reqFormat, userId, ipAddress, userAgent)
	err = s.LogSystemRepository.CreateLogSystem(newLogSystem)
	if err != nil {
		return LogSystem{}, err
	}
	return newLogSystem, nil
}

func (s *LogSystemServiceImpl) ResolveAll(request model.StandardRequest) (orders pagination.Response, err error) {
	return s.LogSystemRepository.ResolveAll(request)
}
