package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"lms-be/internal/domain/auth"
	"lms-be/shared"
	"lms-be/shared/failure"
	"lms-be/shared/model"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type LogSystemHandler struct {
	LogSystemService auth.LogSystemService
}

func ProvideLogSystemHandler(LogSystemService auth.LogSystemService) LogSystemHandler {
	return LogSystemHandler{
		LogSystemService: LogSystemService,
	}
}

func (h *LogSystemHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/log-system", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", h.CreateLogSystem)
			r.Get("/", h.ResolveAll)
		})
	})
}

// CreateLogSystem adalah untuk menambah data Log System.
// @Summary menambahkan data Log System.
// @Description Endpoint ini adalah untuk menambahkan data Log System.
// @Tags Log System
// @Produce json
// @Security BearerAuth
// @Param logSystem body auth.RequestLogSystemFormat true "Log System yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.LogSystem}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/log-system [post]
func (h *LogSystemHandler) CreateLogSystem(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestLogSystemFormat
	fmt.Println("reqFormat", reqFormat)
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	newMenu, err := h.LogSystemService.CreateLogSystem(reqFormat, userID, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenu)
}

// ResolveAll list data Log System.
// @Summary Get list data Log System.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Log System sesuai dengan filter yang dikirimkan.
// @Tags log-system
// @Produce json
// @Security BearerAuth
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [ kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param startDate query string false "Start date filter (YYYY-MM-DD)"
// @Param endDate query string false "End date filter (YYYY-MM-DD)"
// @Success 200 {object} auth.LogSystemDto
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/log-system [get]
func (h *LogSystemHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	roleId := middleware.GetClaimsValue(r.Context(), "roleId").(string)
	if roleId != "HA01" {
		response.WithError(w, failure.Unauthorized("Access denied"))
		return
	}
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "jam"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "DESC"
	}
	startDate := r.URL.Query().Get("startDate")
	endDate := r.URL.Query().Get("endDate")
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
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.LogSystemService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}
