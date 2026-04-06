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

type CompanyHandler struct {
	CompanyService master.CompanyService
	Config         *configs.Config
}

func ProvideCompanyHandler(service master.CompanyService, config *configs.Config) CompanyHandler {
	return CompanyHandler{
		CompanyService: service,
		Config:         config,
	}
}

func (h *CompanyHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/company", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)	
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllData)
			r.Post("/", h.Create)
			r.Post("/import/preview", h.PreviewExcel)
			r.Post("/import", h.ImportExcel)
			r.Put("/{id}", h.Update)
			r.Get("/{id}", h.ResolveByID)
			r.Delete("/{id}", h.DeleteSoft)
		})
	})
}

// ResolveAll mengambil semua data Company dengan filter, sort, dan pagination
// @Summary Get Company with Filter, Sort, and Pagination
// @Description Mengambil semua data Company dengan filter, sort, dan pagination
// @Tags Company
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param pageSize query int true "Page size"
// @Param pageNumber query int true "Page number"
// @Param sortBy query string false "Sort by field"
// @Param sortOrder query string false "Sort order (asc or desc)"
// @Param ignorePaging query bool false "Ignore paging"
// @Success 200 {object} response.Base{data=[]master.CompanyDTO}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company [get]
func (h *CompanyHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.CompanyService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetAllData mengambil semua data Company tanpa filter, sort, dan pagination
// @Summary Get All Company Data
// @Description Mengambil semua data Company tanpa filter, sort, dan pagination
// @Tags Company
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]master.CompanyDTO}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company/all [get]
func (h *CompanyHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	req := model.StandardRequest{}
	status, err := h.CompanyService.GetAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// Create menambahkan data Company baru
// @Summary Create New Company
// @Description Menambahkan data Company baru
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.RequestCompanyFormat true "Request Body"
// @Success 201 {object} response.Base{data=master.Company}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company [post]
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestCompanyFormat
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
	newCompany, err := h.CompanyService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusCreated, newCompany)
}

// Update memperbarui data Company berdasarkan ID
// @Summary Update Company by ID
// @Description Memperbarui data Company berdasarkan ID
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param request body master.RequestCompanyFormat true "Request Body"
// @Success 200 {object} response.Base{data=master.Company}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company/{id} [put]
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WithError(w, failure.BadRequest(errors.New("missing id in path")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid id format")))
		return
	}

	var req master.RequestCompanyFormat
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid JSON body: "+err.Error())))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	req.UserID = userID
	updatedCompany, err := h.CompanyService.Update(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, updatedCompany)
}

// ResolveByID mengambil data Company berdasarkan ID
// @Summary Get Company by ID
// @Description Mengambil data Company berdasarkan ID
// @Tags Company
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} response.Base{data=master.Company}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company/{id} [get]
func (h *CompanyHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid id format")))
		return
	}

	data, err := h.CompanyService.ResolveByID(r.Context(), id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}

// DeleteSoft menghapus data Company secara soft delete berdasarkan ID
// @Summary Soft Delete Company by ID
// @Description Menghapus data Company secara soft delete berdasarkan ID
// @Tags Company
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company/{id} [delete]
func (h *CompanyHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.WithError(w, failure.BadRequest(errors.New("missing id in path")))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid id format")))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = h.CompanyService.DeleteByID(r.Context(), id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, "company deleted successfully")
}

// PreviewExcel menampilkan preview data Company dari file Excel sebelum diimpor
// @Summary Preview Company Data from Excel
// @Description Menampilkan preview data Company dari file Excel sebelum diimpor
// @Tags Company
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Excel File"
// @Success 200 {object} response.Base{data=[]master.PreviewCompanyRow}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company/import/preview [post]
func (h *CompanyHandler) PreviewExcel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("failed to parse form: "+err.Error())))
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("file is required")))
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.WithError(w, failure.InternalError(err))
		return
	}

	result, err := h.CompanyService.PreviewFromExcel(r.Context(), fileBytes)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// ImportFromPreview imports Company data from Excel based on preview
// @Summary Import Company Data from Excel Preview
// @Description Imports Company data from Excel based on preview
// @Tags Company
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.ImportFromPreviewCompanyRequest true "Request Body"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base{data=string}
// @Failure 500 {object} response.Base{data=string}
// @Router /v1/master/company/import [post]
func (h *CompanyHandler) ImportExcel(w http.ResponseWriter, r *http.Request) {
	var req master.ImportFromPreviewCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid JSON: "+err.Error())))
		return
	}

	if len(req.Data) == 0 {
		response.WithError(w, failure.BadRequest(errors.New("data is required")))
		return
	}

	if req.Mode != "" && req.Mode != "insert" && req.Mode != "upsert" {
		response.WithError(w, failure.BadRequest(errors.New("mode must be 'insert' or 'upsert'")))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	result, err := h.CompanyService.ImportFromPreview(r.Context(), req, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}
