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

type MentorAssignmentHandler struct {
	MentorAssignmentService internship.MentorAssignmentService
}

func ProvideMentorAssignmentHandler(service internship.MentorAssignmentService) MentorAssignmentHandler {
	return MentorAssignmentHandler{
		MentorAssignmentService: service,
	}
}

func (h *MentorAssignmentHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/internship/mentor-assignment", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Get("/", h.ResolveAll)
			r.Get("/mentor/{mentorId}/students", h.GetStudentsByMentorID)
			r.Get("/student/{studentId}", h.ResolveMentorByStudentID)
		})
	})
}

// @Summary Assign mentor to student
// @Description Endpoint ini digunakan untuk menugaskan mentor ke mahasiswa.
// @Tags Mentor Assignment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body internship.RequestMentorAssignmentFormat true "Data assignment mentor"
// @Success 201 {object} response.Base{data=internship.MentorAssignment}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/mentor-assignment [post]
func (h *MentorAssignmentHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req internship.RequestMentorAssignmentFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	assignedBy, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	req.AssignedBy = assignedBy

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newMentorAssignment, err := h.MentorAssignmentService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newMentorAssignment)
}

// @Summary Get students by mentor ID
// @Description Endpoint ini digunakan untuk mengambil daftar mahasiswa berdasarkan mentor_id.
// @Tags Mentor Assignment
// @Produce json
// @Security BearerAuth
// @Param mentorId path string true "Mentor ID"
// @Success 200 {object} response.Base{data=[]internship.MentorAssignmentDTO}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/mentor-assignment/mentor/{mentorId}/students [get]
func (h *MentorAssignmentHandler) GetStudentsByMentorID(w http.ResponseWriter, r *http.Request) {
	mentorIDStr := chi.URLParam(r, "mentorId")
	mentorID, err := uuid.FromString(mentorIDStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid mentor id")))
		return
	}

	data, err := h.MentorAssignmentService.GetStudentsByMentorID(r.Context(), mentorID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Get mentor by student ID
// @Description Endpoint ini digunakan untuk mengambil detail mentor berdasarkan student_id.
// @Tags Mentor Assignment
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Success 200 {object} response.Base{data=internship.MentorAssignmentDTO}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/mentor-assignment/student/{studentId} [get]
func (h *MentorAssignmentHandler) ResolveMentorByStudentID(w http.ResponseWriter, r *http.Request) {
	studentIDStr := chi.URLParam(r, "studentId")
	studentID, err := uuid.FromString(studentIDStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid student id")))
		return
	}

	data, err := h.MentorAssignmentService.ResolveMentorByStudentID(r.Context(), studentID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Get all mentor assignments with pagination
// @Description Endpoint ini digunakan untuk mengambil data penugasan mentor dengan pagination.
// @Tags Mentor Assignment
// @Produce json
// @Security BearerAuth
// @Param pageSize query int false "Page size"
// @Param pageNumber query int false "Page number"
// @Success 200 {object} response.Base{data=pagination.Response}
// @Failure 400 {object} response.Base
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/mentor-assignment/ [get]
func (h *MentorAssignmentHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {

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

	req := internship.RequestMentorAssignmentListFormat{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	data, err := h.MentorAssignmentService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}
