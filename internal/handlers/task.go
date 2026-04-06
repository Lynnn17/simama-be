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

type TaskHandler struct {
	TaskService internship.TaskService
}

func ProvideTaskHandler(service internship.TaskService) TaskHandler {
	return TaskHandler{TaskService: service}
}

func (h *TaskHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/internship/task", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Get("/student/{studentId}", h.GetByStudentID)
			r.Get("/mentor/{mentorId}", h.GetByMentorID)
			r.Post("/submit", h.SubmitTaskFile)
			r.Put("/grade/{id}", h.GradeTask)
		})
	})
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req internship.RequestTaskFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newTask, err := h.TaskService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newTask)
}

func (h *TaskHandler) GetByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID, err := uuid.FromString(chi.URLParam(r, "studentId"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid student id")))
		return
	}
	data, err := h.TaskService.GetByStudentID(r.Context(), studentID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

func (h *TaskHandler) GetByMentorID(w http.ResponseWriter, r *http.Request) {
	mentorID, err := uuid.FromString(chi.URLParam(r, "mentorId"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid mentor id")))
		return
	}
	data, err := h.TaskService.GetByMentorID(r.Context(), mentorID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

func (h *TaskHandler) SubmitTaskFile(w http.ResponseWriter, r *http.Request) {
	var req internship.RequestSubmitTaskFileFormat
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	uploadedBy, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	req.UploadedBy = uploadedBy

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newTaskFile, err := h.TaskService.SubmitTaskFile(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newTaskFile)
}

func (h *TaskHandler) GradeTask(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid task id")))
		return
	}

	var req internship.RequestGradeTaskFormat
	err = json.NewDecoder(r.Body).Decode(&req)
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

	newTask, err := h.TaskService.GradeTask(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newTask)
}
