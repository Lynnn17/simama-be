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

type PersonnelHandler struct {
	PersonnelService master.PersonnelService
	Config          *configs.Config
}

func ProvidePersonnelHandler(service master.PersonnelService, config *configs.Config) PersonnelHandler {
	return PersonnelHandler{
		PersonnelService: service,
		Config:          config,
	}
}

func (h *PersonnelHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/personnel", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Post("/", h.Create)
			r.Get("/all", h.GetAllData)
			r.Get("/{id}", h.ResolveByID)
			r.Put("/{id}", h.Update)
			r.Delete("/{id}", h.DeleteSoft)
			r.Post("/import/preview", h.PreviewExcel)
			r.Post("/import", h.ImportFromPreview)
		})
	})
}

// ResolveAll mengambil semua data Personnel dengan filter, sort, dan pagination
// @Summary Get Personnel with Filter, Sort, and Pagination
// @Description Mengambil semua data Personnel dengan filter, sort, dan pagination
// @Tags Personnel
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param pageSize query int true "Page size"
// @Param pageNumber query int true "Page number"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order (asc or desc)"
// @Param ignorePaging query bool false "Ignore paging"
// @Success 200 {object} response.Base{data=[]master.PersonnelDTO}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel [get]
func (h *PersonnelHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.PersonnelService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetAllData mengambil semua data Personnel tanpa filter, sort, dan pagination
// @Summary Get All Personnel Data
// @Description Mengambil semua data Personnel tanpa filter, sort, dan pagination
// @Tags Personnel
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]master.PersonnelDTO}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel/all [get]
func (h *PersonnelHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	req := model.StandardRequest{}
	status, err := h.PersonnelService.GetAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// Create menambahkan data Personnel baru
// @Summary Create New Personnel
// @Description Menambahkan data Personnel baru
// @Tags Personnel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.RequestPersonnelFormat true "Request Body"
// @Success 201 {object} response.Base{data=master.Personnel}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel [post]
func (h *PersonnelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestPersonnelFormat
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
	newPersonnel, err := h.PersonnelService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newPersonnel)
}

// Update memperbarui data Personnel berdasarkan ID
// @Summary Update Personnel by ID
// @Description Memperbarui data Personnel berdasarkan ID
// @Tags Personnel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Personnel ID"
// @Param request body master.RequestPersonnelFormat true "Request Body"
// @Success 200 {object} response.Base{data=master.Personnel}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel/{id} [put]
func (h *PersonnelHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req master.RequestPersonnelFormat
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
	updatedPersonnel, err := h.PersonnelService.Update(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, updatedPersonnel)
}

// ResolveByID mengambil data Personnel berdasarkan ID
// @Summary Get Personnel by ID
// @Description Mengambil data Personnel berdasarkan ID
// @Tags Personnel
// @Produce json
// @Security BearerAuth
// @Param id path int true "Personnel ID"
// @Success 200 {object} response.Base{data=master.Personnel}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel/{id} [get]
func (h *PersonnelHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	data, err := h.PersonnelService.ResolveByID(r.Context(), id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// DeleteSoft menghapus data Personnel secara soft delete berdasarkan ID
// @Summary Soft Delete Personnel by ID
// @Description Menghapus data Personnel secara soft delete berdasarkan ID
// @Tags Personnel
// @Produce json
// @Security BearerAuth
// @Param id path int true "Personnel ID"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel/{id} [delete]
func (h *PersonnelHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
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

	err = h.PersonnelService.DeleteByID(r.Context(), id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK,
		"personnel deleted successfully",
	)
}

// PreviewExcel menampilkan preview data Personnel dari file Excel sebelum diimpor
// @Summary Preview Personnel Data from Excel
// @Description Menampilkan preview data Personnel dari file Excel sebelum diimpor
// @Tags Personnel
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Excel File"
// @Success 200 {object} response.Base{data=[]master.PreviewPersonnelRow}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel/import/preview [post]
func (h *PersonnelHandler) PreviewExcel(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.PersonnelService.PreviewFromExcel(r.Context(), fileBytes)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// ImportFromPreview imports Personnel data from Excel based on preview
// @Summary Import Personnel Data from Excel Preview
// @Description Imports Personnel data from Excel based on preview
// @Tags Personnel
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.ImportFromPreviewPersonnelRequest true "Request Body"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/personnel/import [post]
func (h *PersonnelHandler) ImportFromPreview(w http.ResponseWriter, r *http.Request) {
	var req master.ImportFromPreviewPersonnelRequest
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

	result, err := h.PersonnelService.ImportFromPreview(r.Context(), req, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}
