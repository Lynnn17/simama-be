package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
			r.Put("/{id}", h.Update)
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
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search activity/blocker/plan"
// @Param progressStatus query string false "Progress status filter (in_progress, done, blocked)"
// @Success 200 {object} response.Base{data=pagination.Response{items=[]internship.LogbookDTO}}
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

	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	search := r.URL.Query().Get("search")
	progressStatus := r.URL.Query().Get("progressStatus")
	date := r.URL.Query().Get("date")

	req := internship.RequestLogbookListFormat{
		PageSize:       pageSize,
		PageNumber:     pageNumber,
		Search:         search,
		ProgressStatus: progressStatus,
		Date:           date,
	}

	data, err := h.LogbookService.ResolveAllByStudentID(r.Context(), studentID, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Get logbooks by mentor ID
// @Description Endpoint ini digunakan untuk mengambil logbook mahasiswa bimbingan mentor (read-only).
// @Tags Logbook
// @Produce json
// @Security BearerAuth
// @Param mentorId path string true "Mentor ID"
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search activity/blocker/plan/student name"
// @Param progressStatus query string false "Progress status filter (in_progress, done, blocked)"
// @Param date query string false "Date filter (YYYY-MM-DD)"
// @Success 200 {object} response.Base{data=pagination.Response{items=[]internship.LogbookDTO}}
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

	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	search := r.URL.Query().Get("search")
	progressStatus := r.URL.Query().Get("progressStatus")
	date := r.URL.Query().Get("date")

	req := internship.RequestLogbookListFormat{
		PageSize:       pageSize,
		PageNumber:     pageNumber,
		Search:         search,
		ProgressStatus: progressStatus,
		Date:           date,
	}

	data, err := h.LogbookService.ResolveAllByMentorID(r.Context(), mentorID, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Update logbook
// @Description Endpoint ini digunakan untuk mengupdate logbook harian.
// @Tags Logbook
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Logbook ID"
// @Param request body internship.RequestLogbookFormat true "Data logbook"
// @Success 200 {object} response.Base{data=internship.Logbook}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/logbook/{id} [put]
func (h *LogbookHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var req internship.RequestLogbookFormat
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
		response.WithError(w, failure.Unauthorized("invalid user id in token"))
		return
	}
	req.StudentID = userID
	req.ID = id

	res, err := h.LogbookService.Update(r.Context(), id, req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, res)
}
