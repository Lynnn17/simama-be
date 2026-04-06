package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid/v5"

	"lms-be/internal/domain/internship"
	"lms-be/shared"
	"lms-be/shared/failure"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
)

type LogbookHandler struct {
	LogbookService internship.LogbookService
}

func ProvideLogbookHandler(service internship.LogbookService) LogbookHandler {
	return LogbookHandler{LogbookService: service}
}

func (h *LogbookHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/internship/logbook", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Get("/student/{studentId}", h.GetByStudentID)
			r.Get("/mentor/{mentorId}", h.GetByMentorID)
			r.Put("/status/{id}", h.UpdateStatus)
		})
	})
}

// @Summary Create logbook
// @Description Endpoint ini digunakan untuk membuat logbook harian.
// @Tags Logbook
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body internship.RequestLogbookFormat true "Data logbook"
// @Success 201 {object} response.Base{data=internship.Logbook}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/logbook [post]
func (h *LogbookHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req internship.RequestLogbookFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	studentID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	req.StudentID = studentID

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newLogbook, err := h.LogbookService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newLogbook)
}

// @Summary Get logbooks by student ID
// @Description Endpoint ini digunakan untuk mengambil riwayat logbook mahasiswa.
// @Tags Logbook
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Success 200 {object} response.Base{data=[]internship.LogbookDTO}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/logbook/student/{studentId} [get]
func (h *LogbookHandler) GetByStudentID(w http.ResponseWriter, r *http.Request) {
	studentIDStr := chi.URLParam(r, "studentId")
	studentID, err := uuid.FromString(studentIDStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid student id")))
		return
	}

	data, err := h.LogbookService.GetByStudentID(r.Context(), studentID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Get logbooks by mentor ID
// @Description Endpoint ini digunakan untuk mengambil logbook mahasiswa berdasarkan mentor assignment.
// @Tags Logbook
// @Produce json
// @Security BearerAuth
// @Param mentorId path string true "Mentor ID"
// @Success 200 {object} response.Base{data=[]internship.LogbookDTO}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/logbook/mentor/{mentorId} [get]
func (h *LogbookHandler) GetByMentorID(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := chi.URLParam(r, "mentorId")
	mentorID, err := uuid.FromString(mentorIDStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid mentor id")))
		return
	}

	data, err := h.LogbookService.GetByMentorID(r.Context(), mentorID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Update logbook status
// @Description Endpoint ini digunakan untuk approve/reject logbook oleh mentor.
// @Tags Logbook
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Logbook ID"
// @Param request body internship.RequestUpdateLogbookStatusFormat true "Status logbook"
// @Success 200 {object} response.Base{data=internship.Logbook}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/logbook/status/{id} [put]
func (h *LogbookHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid logbook id")))
		return
	}

	var req internship.RequestUpdateLogbookStatusFormat
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	mentorID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	req.UserID = mentorID

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newLogbook, err := h.LogbookService.UpdateStatus(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newLogbook)
}
