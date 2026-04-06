package auth

import (
	"context"
	"lms-be/configs"
	"lms-be/internal/files"
	"lms-be/shared/logger"
	"net/http"

	"github.com/gofrs/uuid"
)

type AppConfigService interface {
	ResolveByID(ctx context.Context, id int, isConfig bool) (data AppConfigDTO, err error)
	ResolveDTOByID(ctx context.Context, id int) (data AppConfigDTOByID, err error)
	Update(ctx context.Context, id int, req RequestAppConfigFormat, userId uuid.UUID) (data AppConfig, err error)
	UploadFile(w http.ResponseWriter, r *http.Request) (path string, err error)
}

type AppConfigServiceImpl struct {
	AppConfigRepository AppConfigRepository
	FileService         files.FileService
	Config              *configs.Config
}

func ProvideAppConfigServiceImpl(AppConfigRepository AppConfigRepository, file files.FileService, config *configs.Config) *AppConfigServiceImpl {
	s := new(AppConfigServiceImpl)
	s.AppConfigRepository = AppConfigRepository
	s.FileService = file
	s.Config = config
	return s
}

func (s *AppConfigServiceImpl) ResolveByID(ctx context.Context, id int, isConfig bool) (data AppConfigDTO, err error) {
	return s.AppConfigRepository.ResolveByID(ctx, id, isConfig)
}

func (s *AppConfigServiceImpl) Update(ctx context.Context, id int, req RequestAppConfigFormat, userId uuid.UUID) (data AppConfig, err error) {
	data, err = s.AppConfigRepository.GetByID(ctx, id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	data.AppConfigFormat(req, userId)

	err = s.AppConfigRepository.Update(ctx, data)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (s *AppConfigServiceImpl) UploadFile(w http.ResponseWriter, r *http.Request) (path string, err error) {
	path, err = s.FileService.UploadFile(s.Config.App.File.Config, w, r)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (s *AppConfigServiceImpl) ResolveDTOByID(ctx context.Context, id int) (data AppConfigDTOByID, err error) {
	return s.AppConfigRepository.ResolveDTOByID(ctx, id)
}
