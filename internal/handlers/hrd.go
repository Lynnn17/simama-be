package handlers

import (
	"net/http"

	"github.com/go-chi/chi"

	"lms-be/internal/domain/internship"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"
)

type HRDHandler struct {
	HRDMonitoringService internship.HRDMonitoringService
}

func ProvideHRDHandler(service internship.HRDMonitoringService) HRDHandler {
	return HRDHandler{
		HRDMonitoringService: service,
	}
}

func (h *HRDHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/hrd", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/monitoring", h.GetMonitoringData)
		})
	})
}

// @Summary Get HRD monitoring data
// @Description Endpoint ini digunakan oleh HRD untuk memantau kehadiran dan logbook mahasiswa hari ini.
// @Tags HRD
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari nama atau universitas"
// @Success 200 {object} response.Base{data=[]internship.HRDMonitoringDTO}
// @Failure 401 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/hrd/monitoring [get]
func (h *HRDHandler) GetMonitoringData(w http.ResponseWriter, r *http.Request) {
	req := internship.RequestHRDMonitoringFormat{
		Search: r.URL.Query().Get("search"),
	}

	data, err := h.HRDMonitoringService.GetMonitoringData(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, data)
}
