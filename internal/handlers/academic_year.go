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

type AcademicYearHandler struct {
	AcademicYearService master.AcademicYearService
	Config              *configs.Config
}

func ProvideAcademicYearHandler(service master.AcademicYearService, config *configs.Config) AcademicYearHandler {
	return AcademicYearHandler{
		AcademicYearService: service,
		Config:              config,
	}
}

func (h *AcademicYearHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/master/academic-year", func(r chi.Router) {
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

// ResolveAll mengambil daftar data Academic Year dengan filter, sort, dan pagination.
// @Summary Daftar Academic Year dengan filter, sort, dan pagination
// @Description Endpoint ini digunakan untuk mengambil daftar data Academic Year berdasarkan pencarian (q), urutan (sort), dan pagination.
// @Tags Academic Year
// @Produce json
// @Security BearerAuth
// @Param q query string false "Kata kunci pencarian"
// @Param pageSize query int true "Jumlah data per halaman"
// @Param pageNumber query int true "Nomor halaman yang diambil"
// @Param sortBy query string false "Parameter pengurutan [name, code]"
// @Param sortType query string false "Tipe pengurutan [asc | desc]"
// @Param ignorePaging query bool false "Jika true, pagination diabaikan dan semua data dikembalikan (default: false)"
// @Success 200 {object} response.Base{data=[]master.AcademicYearDTO}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year [get]
func (h *AcademicYearHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
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

	status, err := h.AcademicYearService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, status)
}

// GetAllData mengambil seluruh data Academic Year tanpa pagination.
// @Summary Ambil semua data Academic Year
// @Description Endpoint ini digunakan untuk mengambil semua data academic year tanpa pagination.
// @Tags Academic Year
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]master.AcademicYear}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year/all [get]
func (h *AcademicYearHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	req := model.StandardRequest{}
	status, err := h.AcademicYearService.GetAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, status)
}

// Create menambahkan data academic year baru.
// @Summary Tambah data Academic Year baru
// @Description Endpoint ini digunakan untuk menambahkan academic year baru.
// @Tags Academic Year
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param academic_year body master.RequestAcademicYearFormat true "Data academic year"
// @Success 200 {object} response.Base{data=master.AcademicYear}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year [post]
func (h *AcademicYearHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req master.RequestAcademicYearFormat
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
	newAcademicYear, err := h.AcademicYearService.Create(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, newAcademicYear)
}

// Update memperbarui data academic year berdasarkan ID.
// @Summary Perbarui data Academic Year
// @Description Endpoint ini digunakan untuk memperbarui data academic year berdasarkan ID yang dikirimkan di path.
// @Tags Academic Year
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Academic Year"
// @Param academic_year body master.RequestAcademicYearFormat true "Data academic year yang diperbarui"
// @Success 200 {object} response.Base{data=master.AcademicYear}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year/{id} [put]
func (h *AcademicYearHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req master.RequestAcademicYearFormat
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
	newAcademicYear, err := h.AcademicYearService.Update(r.Context(), id, req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, newAcademicYear)
}

// ResolveByID mengambil data academic year berdasarkan ID.
// @Summary Ambil detail Academic Year by ID
// @Description Endpoint ini digunakan untuk mengambil satu data academic year berdasarkan ID.
// @Tags Academic Year
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Academic Year"
// @Success 200 {object} response.Base{data=master.AcademicYear}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year/{id} [get]
func (h *AcademicYearHandler) ResolveByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid id format")))
		return
	}

	data, err := h.AcademicYearService.ResolveByID(r.Context(), id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, data)
}

// DeleteSoft menghapus data academic year secara soft delete berdasarkan ID.
// @Summary Hapus Academic Year by ID (Soft Delete)
// @Description Endpoint ini digunakan untuk menghapus data academic year berdasarkan ID (soft delete).
// @Tags Academic Year
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Academic Year"
// @Success 200 {object} response.Base{data=string}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year/{id} [delete]
func (h *AcademicYearHandler) DeleteSoft(w http.ResponseWriter, r *http.Request) {
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

	err = h.AcademicYearService.Delete(r.Context(), id, userID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, "Academic Year successfully deleted")
}

// PreviewExcel previews Academic Year data from an Excel file without saving.
// @Summary Preview Import Academic Year from Excel
// @Description Preview data dari file Excel sebelum import. Response menunjukkan data yang valid, error, dan apakah data sudah ada di database.
// @Tags Academic Year
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param institutionId formData int true "Institution ID"
// @Param file formData file true "Excel file to preview"
// @Success 200 {object} response.Base{data=master.PreviewAcademicYearResult}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year/import/preview [post]
func (h *AcademicYearHandler) PreviewExcel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("failed to parse form: "+err.Error())))
		return
	}

	institutionIdStr := r.FormValue("institutionId")
	if institutionIdStr == "" {
		response.WithError(w, failure.BadRequest(errors.New("institutionId is required")))
		return
	}
	institutionID, err := strconv.Atoi(institutionIdStr)
	if err != nil {
		response.WithError(w, failure.BadRequest(errors.New("invalid institutionId")))
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

	result, err := h.AcademicYearService.PreviewFromExcel(r.Context(), fileBytes, institutionID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}

// ImportFromPreview imports Academic Year data from preview JSON.
// @Summary Import Academic Year from Preview Data
// @Description Import data dari preview yang sudah di-review. Kirim JSON, bukan file.
// @Tags Academic Year
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body master.ImportFromPreviewRequest true "Import request with preview data"
// @Success 200 {object} response.Base{data=model.ImportResult}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/master/academic-year/import [post]
func (h *AcademicYearHandler) ImportExcel(w http.ResponseWriter, r *http.Request) {
	var req master.ImportFromPreviewRequest
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

	result, err := h.AcademicYearService.ImportFromPreview(r.Context(), req, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, result)
}
