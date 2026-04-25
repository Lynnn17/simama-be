package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid/v5"

	"lms-be/configs"
	"lms-be/internal/domain/internship"
	"lms-be/shared"
	"lms-be/shared/failure"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
)

type TaskHandler struct {
	TaskService internship.TaskService
	Config      *configs.Config
}

func ProvideTaskHandler(service internship.TaskService, config *configs.Config) TaskHandler {
	return TaskHandler{TaskService: service, Config: config}
}

func (h *TaskHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/internship/task", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.Create)
			r.Get("/student/{studentId}", h.GetByStudentID)
			r.Get("/mentor/{mentorId}", h.GetByMentorID)
			r.Put("/{id}", h.Update)
			r.Post("/submit", h.SubmitTaskFile)
			r.Put("/grade/{id}", h.GradeTask)
			r.Get("/{id}/files", h.GetFiles)
			r.Get("/{id}/download-all", h.DownloadAll)
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

// @Summary Get tasks by student ID
// @Description Endpoint ini digunakan untuk mengambil daftar tugas berdasarkan ID mahasiswa.
// @Tags Task
// @Produce json
// @Security BearerAuth
// @Param studentId path string true "Student ID"
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search title/description"
// @Param status query string false "Status filter"
// @Param date query string false "Date filter (deadline)"
// @Success 200 {object} response.Base{data=pagination.Response{items=[]internship.TaskDTO}}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/task/student/{studentId} [get]
func (h *TaskHandler) GetByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID, err := uuid.FromString(chi.URLParam(r, "studentId"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid student id")))
		return
	}

	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")

	req := internship.RequestTaskListFormat{
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Search:     search,
		Status:     status,
		Date:       date,
	}

	data, err := h.TaskService.ResolveAllByStudentID(r.Context(), studentID, req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Get tasks by mentor ID
// @Description Endpoint ini digunakan untuk mengambil daftar tugas berdasarkan ID mentor.
// @Tags Task
// @Produce json
// @Security BearerAuth
// @Param mentorId path string true "Mentor ID"
// @Param pageNumber query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search title/description"
// @Param status query string false "Status filter"
// @Param date query string false "Date filter (deadline)"
// @Param studentSearch query string false "Search student name"
// @Success 200 {object} response.Base{data=pagination.Response{items=[]internship.TaskDTO}}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/task/mentor/{mentorId} [get]
func (h *TaskHandler) GetByMentorID(w http.ResponseWriter, r *http.Request) {
	mentorID, err := uuid.FromString(chi.URLParam(r, "mentorId"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid mentor id")))
		return
	}

	pageNumber, _ := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")
	date := r.URL.Query().Get("date")
	studentSearch := r.URL.Query().Get("studentSearch")

	req := internship.RequestTaskListFormat{
		PageSize:      pageSize,
		PageNumber:    pageNumber,
		Search:        search,
		Status:        status,
		Date:          date,
		StudentSearch: studentSearch,
	}

	data, err := h.TaskService.ResolveAllByMentorID(r.Context(), mentorID, req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// @Summary Update task
// @Description Endpoint ini digunakan oleh mentor untuk mengubah detail tugas (judul, deskripsi, deadline).
// @Tags Task
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Param payload body internship.RequestTaskFormat true "Update payload"
// @Success 200 {object} response.Base{data=internship.Task}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/internship/task/{id} [put]
func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid task id")))
		return
	}

	var req internship.RequestTaskFormat
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	req.MentorID, err = uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newTask, err := h.TaskService.Update(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, newTask)
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

func (h *TaskHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid task id")))
		return
	}

	data, err := h.TaskService.GetFilesByTaskID(r.Context(), id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

func (h *TaskHandler) DownloadAll(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid task id")))
		return
	}

	zipData, filename, err := h.TaskService.CreateZipFromTaskFiles(r.Context(), id, h.Config.App.File.Dir)
	if err != nil {
		response.WithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(zipData)
}
