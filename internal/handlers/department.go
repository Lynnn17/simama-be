package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"lms-be/configs"
	"lms-be/internal/domain/master"
	"lms-be/shared"
	"lms-be/shared/failure"
	"lms-be/shared/model"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid/v5"
)

type DepartmentHandler struct {
	DepartmentService master.DepartmentService
	Config            *configs.Config
}

func ProvideDepartmentHandler(service master.DepartmentService, config *configs.Config) DepartmentHandler {
	return DepartmentHandler{
		DepartmentService: service,
		Config:            config,
	}
}

func (h *DepartmentHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/department", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
			r.Put("/{id}", h.Update)
			r.Get("/{id}", h.ResolveByID)
			r.Delete("/{id}", h.DeleteSoft)
			r.Post("/import/preview", h.PreviewExcel)
			r.Post("/import", h.ImportExcel)
		})
	})
}

// ResolveAll mengambil semua data Department dengan filter, sort, dan pagination
// @Summary Get Department with Filter, Sort, and Pagination
// @Description Mengambil semua data Department dengan filter, sort, dan pagination
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param pageSize query int true "Page size"
// @Param pageNumber query int true "Page number"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order (asc or desc)"
// @Param ignorePaging query bool false "Ignore paging"
// @Success 200 {object} response.Base{data=[]master.DepartmentDTO}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department [get]
func (h *DepartmentHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	ignorePagingStr := r.URL.Query().Get("ignorePaging")

	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "desc"
	}

	ignorePaging, _ := strconv.ParseBool(ignorePagingStr)

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	pageNumber, err := strconv.Atoi(pageNumberStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	req := model.StandardRequest{
		Keyword:      keyword,
		PageSize:     pageSize,
		PageNumber:   pageNumber,
		SortBy:       sortBy,
		SortType:     sortType,
		IgnorePaging: ignorePaging,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.DepartmentService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetAllData mengambil semua data Department tanpa filter, sort, dan pagination
// @Summary Get All Department Data
// @Description Mengambil semua data Department tanpa filter, sort, dan pagination
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]master.DepartmentDTO}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department/all [get]
func (h *DepartmentHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	req := model.StandardRequest{}
	status, err := h.DepartmentService.GetAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// Create menambahkan data Department baru
// @Summary Create New Department
// @Description Menambahkan data Department baru
// @Tags Department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.RequestDepartmentFormat true "Request Body"
// @Success 201 {object} response.Base{data=master.Department}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department [post]
func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestDepartmentFormat
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
	newDepartment, err := h.DepartmentService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newDepartment)
}

// Update memperbarui data Department berdasarkan ID
// @Summary Update Department by ID
// @Description Memperbarui data Department berdasarkan ID
// @Tags Department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Param request body master.RequestDepartmentFormat true "Request Body"
// @Success 200 {object} response.Base{data=master.Department}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department/{id} [put]
func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WithError(w, failure.BadRequest(errors.New("invalid department id")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	var req master.RequestDepartmentFormat
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
	updatedDepartment, err := h.DepartmentService.Update(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, updatedDepartment)
}

// ResolveByID mengambil data Department berdasarkan ID
// @Summary Get Department by ID
// @Description Mengambil data Department berdasarkan ID
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} response.Base{data=master.Department}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department/{id} [get]
func (h *DepartmentHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.DepartmentService.ResolveByID(r.Context(), id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// DeleteSoft menghapus data Department secara soft delete berdasarkan ID
// @Summary Soft Delete Department by ID
// @Description Menghapus data Department secara soft delete berdasarkan ID
// @Tags Department
// @Produce json
// @Security BearerAuth
// @Param id path int true "Department ID"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department/{id} [delete]
func (h *DepartmentHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WithError(w, failure.BadRequest(errors.New("invalid department id")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = h.DepartmentService.DeleteByID(r.Context(), id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "department deleted successfully")
}

// PreviewExcel menampilkan preview data Department dari file Excel sebelum diimpor
// @Summary Preview Department Data from Excel
// @Description Menampilkan preview data Department dari file Excel sebelum diimpor
// @Tags Department
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Excel File"
// @Success 200 {object} response.Base{data=[]master.PreviewDepartmentRow}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department/import/preview [post]
func (h *DepartmentHandler) PreviewExcel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	result, err := h.DepartmentService.PreviewFromExcel(r.Context(), fileBytes)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// ImportFromPreview imports Department data from Excel based on preview
// @Summary Import Department Data from Excel Preview
// @Description Imports Department data from Excel based on preview
// @Tags Department
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.ImportFromPreviewDepartmentRequest true "Request Body"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/department/import [post]
func (h *DepartmentHandler) ImportExcel(w http.ResponseWriter, r *http.Request) {
	var req master.ImportFromPreviewDepartmentRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if len(req.Data) == 0 {
		response.WithError(w, failure.BadRequest(errors.New("data is empty")))
		return
	}

	if req.Mode != "" && req.Mode != "insert" && req.Mode != "upsert" {
		response.WithError(w, failure.BadRequest(errors.New("invalid mode")))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	result, err := h.DepartmentService.ImportFromPreview(r.Context(), req, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}
