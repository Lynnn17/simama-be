package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
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
		r.Post("/", h.Create)
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Put("/status/{id}", h.UpdateStatus)

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

	// Get userID from claims if available (optional)
	var userIDPtr *uuid.UUID
	val := r.Context().Value(middleware.ValueKeyContext)
	if val != nil {
		if claims, ok := val.(jwt.MapClaims); ok {
			if userIdStr, ok := claims["userId"].(string); ok {
				id, err := uuid.FromString(userIdStr)
				if err == nil {
					userIDPtr = &id
				}
			}
		}
	}

	req.UserID = userIDPtr
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

// @Summary Ambil data pendaftaran magang dengan pagination
// @Description Endpoint ini digunakan untuk mengambil data pendaftaran magang dengan pagination.
// @Tags Internship Registration
// @Produce json
// @Security BearerAuth
// @Param pageSize query int false "Page size"
// @Param pageNumber query int false "Page number"
// @Success 200 {object} response.Base{data=pagination.Response}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/registration/ [get]
func (h *RegistrationHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {

	pageSize := 10
	pageNumber := 1

	pageSizeStr := r.URL.Query().Get("pageSize")
	if pageSizeStr != "" {
		parsed, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		pageSize = parsed
	}

	pageNumberStr := r.URL.Query().Get("pageNumber")
	if pageNumberStr != "" {
		parsed, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		pageNumber = parsed
	}

	req := internship.RequestRegistrationListFormat{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Status:     r.URL.Query().Get("status"),
	}

	err := shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.RegistrationService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
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
	req := internship.RequestRegistrationListFormat{
		Status: r.URL.Query().Get("status"),
	}

	data, err := h.RegistrationService.GetAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Update status registration magang
// @Description Endpoint ini digunakan untuk accept/reject registration magang.
// @Tags Internship Registration
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Registration ID"
// @Param request body internship.RequestUpdateRegistrationStatusFormat true "Status registration"
// @Success 200 {object} response.Base{data=internship.Registration}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/registration/status/{id} [put]
func (h *RegistrationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid registration id")))
		return
	}

	var req internship.RequestUpdateRegistrationStatusFormat
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	req.UserID = &userID

	updated, err := h.RegistrationService.UpdateStatus(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, updated)
}
