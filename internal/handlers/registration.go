package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid/v5"

	"lms-be/internal/domain/internship"
	"lms-be/shared"
	"lms-be/shared/failure"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
)

type RegistrationHandler struct {
	RegistrationService internship.RegistrationService
}

func ProvideRegistrationHandler(service internship.RegistrationService) RegistrationHandler {
	return RegistrationHandler{
		RegistrationService: service,
	}
}

func (h *RegistrationHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/internship/registration", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.GetAllData)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
		})
	})
}

// @Summary Tambah data internship registration baru
// @Description Endpoint ini digunakan untuk menambahkan pendaftaran magang baru.
// @Tags Internship Registration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param registration body internship.RequestRegistrationFormat true "Data pendaftaran magang"
// @Success 201 {object} response.Base{data=internship.Registration}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/registration [post]
func (h *RegistrationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req internship.RequestRegistrationFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	req.UserID = userID
	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newRegistration, err := h.RegistrationService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newRegistration)
}

// @Summary Ambil semua data pendaftaran magang
// @Description Endpoint ini digunakan untuk mengambil seluruh data pendaftaran magang.
// @Tags Internship Registration
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]internship.RegistrationDTO}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/registration/all [get]
func (h *RegistrationHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	roleID := middleware.GetClaimsValue(r.Context(), "roleId").(string)
	if roleID != "HA01" {
		response.WithError(w, failure.Unauthorized("Access denied"))
		return
	}

	data, err := h.RegistrationService.GetAll(r.Context())
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}
