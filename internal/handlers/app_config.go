package handlers

import (
	"encoding/json"
	"errors"
	"lms-be/configs"
	"lms-be/internal/domain/auth"
	"lms-be/shared/failure"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type AppConfigHandler struct {
	AppConfigService auth.AppConfigService
	Config           *configs.Config
}

func ProvideAppConfigHandler(service auth.AppConfigService, config *configs.Config) AppConfigHandler {
	return AppConfigHandler{
		AppConfigService: service,
		Config:           config,
	}
}

func (h *AppConfigHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/app-config", func(r chi.Router) {
		r.Get("/public/{id}", h.ResolveDTOByID)
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Put("/{id}", h.Update)
			r.Get("/{id}", h.ResolveByID)
			r.Post("/upload", h.UploadFile)
		})
	})
}

// Update memperbarui data app config berdasarkan ID.
// @Summary Perbarui data AppConfig
// @Description Endpoint ini digunakan untuk memperbarui data app config berdasarkan ID yang dikirimkan di path.
// @Tags AppConfig
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID AppConfig"
// @Param config body auth.RequestAppConfigFormat true "Data app config yang diperbarui"
// @Success 200 {object} response.Base{data=auth.AppConfig}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/app-config/{id} [put]
func (h *AppConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var req auth.RequestAppConfigFormat
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid JSON body: "+err.Error())))
		return
	}

	userId, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.AppConfigService.Update(r.Context(), id, req, userId)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByID mengambil data app config berdasarkan ID.
// @Summary Ambil detail AppConfig by ID
// @Description Endpoint ini digunakan untuk mengambil satu data app config berdasarkan ID.
// @Tags AppConfig
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID AppConfig"
// @Success 200 {object} response.Base{data=auth.AppConfig}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/app-config/{id} [get]
func (h *AppConfigHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.AppConfigService.ResolveByID(r.Context(), id, false)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// ResolveByDTO mengambil data app config berdasarkan ID.
// @Summary Ambil detail AppConfig DTO by ID
// @Description Endpoint ini digunakan untuk mengambil satu data app config berdasarkan ID.
// @Tags AppConfig
// @Produce json
// @Param id path string true "ID AppConfig"
// @Success 200 {object} response.Base{data=auth.AppConfig}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/app-config/public/{id} [get]
func (h *AppConfigHandler) ResolveDTOByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.AppConfigService.ResolveDTOByID(r.Context(), id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// UploadFile untuk upload file attachment
// @Summary Upload file attachment
// @Description End point ini digunakan untuk mengupload file attachment
// @Tags AppConfig
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/app-config/upload [post]
func (h *AppConfigHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	path, err := h.AppConfigService.UploadFile(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.WithJSON(w, http.StatusOK, path)
}
