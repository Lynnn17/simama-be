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
	"github.com/rs/zerolog/log"
)

// MenuHandler adalah HTTP handler untuk domain Role
type MenuHandler struct {
	MenuService auth.MenuService
}

// ProvideMenuHandler adalah provider untuk handler ini
func ProvideMenuHandler(MenuService auth.MenuService) MenuHandler {
	return MenuHandler{
		MenuService: MenuService,
	}
}

// Router untuk setup dari router untuk domain ini
func (h *MenuHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/menu-role", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveMenuByRoleID)
			r.Get("/trx", h.ResolveMenuByRoleIDTrx)
			r.Post("/bulk", h.CreateBulkMenuRole)
			r.Put("/update-permission", h.UpdateMenuPermission)
		})
	})
	r.Route("/menu", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Get("/", h.ResolveAll)
			r.Get("/all", h.GetAllMenu)
			r.Get("/{id}", h.ResolveMenuByID)
			r.Post("/", h.CreateMenu)
			r.Put("/{id}", h.UpdateMenu)
			r.Delete("/{id}", h.DeleteMenu)
		})
	})
}

/* =============== FUNCTION MASTER MENU ============== */

// GetAllMenu list all menu.
// @Summary Get list all menu.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data menu sesuai dengan filter yang dikirimkan.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/all [get]
func (h *MenuHandler) GetAllMenu(w http.ResponseWriter, r *http.Request) {
	status, err := h.MenuService.GetAllMenu()
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// ResolveAll list all Menu.
// @Summary Get list all Menu.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data Menu sesuai dengan filter yang dikirimkan.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [id | namaMenu ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu [get]
func (h *MenuHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
	fmt.Println("pageSizeStr", pageSizeStr)
	fmt.Println("pageNumberStr", pageNumberStr)
	sortBy := r.URL.Query().Get("sortBy")
	if sortBy == "" {
		sortBy = "createdAt"
	}

	sortType := r.URL.Query().Get("sortType")
	if sortType == "" {
		sortType = "DESC"
	}
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

	req := model.StandardRequestMenu{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.MenuService.ResolveAll(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// CreateMenu adalah untuk menambah data Menu.
// @Summary menambahkan data Menu.
// @Description Endpoint ini adalah untuk menambahkan data Menu.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param Menu body auth.RequestMenuFormat true "Menu yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.Menu}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu [post]
func (h *MenuHandler) CreateMenu(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestMenuFormat
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	reqFormat.UserID = userID
	newMenu, err := h.MenuService.CreateMenu(reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenu)
}

// ResolveMenuByID adalah untuk mendapatkan satu data Menu berdasarkan ID.
// @Summary Mendapatkan satu data Menu berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan Menu ID.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=auth.Menu}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/{id} [get]
func (h *MenuHandler) ResolveMenuByID(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	menu, err := h.MenuService.ResolveMenuByID(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, menu)
}

// UpdateMenu adakan untuk update data Menu
// @Summary update data Menu
// @Description endpoint ini adalah untuk mengubah data Menu
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Param Menu body auth.RequestMenuFormat true "Menu yang akan ditambahkan"
// @Success 200 {object} response.Base{data=[]auth.Menu}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/{id} [put]
func (h *MenuHandler) UpdateMenu(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
	}

	var newMenu auth.RequestMenuFormat
	err = json.NewDecoder(r.Body).Decode(&newMenu)

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		fmt.Print("error user id")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newMenu.UserID = userID
	menu, err := h.MenuService.UpdateMenu(id, newMenu)
	if err != nil {
		log.Info().Msg("Error: " + err.Error())
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, menu)
}

// DeleteMenu adalah untuk menghapus data Menu.
// @Summary hapus data Menu.
// @Description Endpoint ini adalah untuk menghapus data Menu.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu/{id} [delete]
func (h *MenuHandler) DeleteMenu(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newID, err := uuid.FromString(id)

	err = h.MenuService.DeleteMenuByID(newID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Data Berhasil di Hapus",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

/* =============== FUNCTION MENU ROLE ACCESS =============== */

// resolveMenuByID adalah untuk mendapatkan semua menu berdasarkan RoleID.
// @Summary Mendapatkan semua data Menu.
// @Description Endpoint ini adalah untuk mendapatkan Role dan semua menu berdasarkan RoleID.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param roleId query string true "Set RoleID"
// @Success 200 {object} response.Base{data=auth.MenuResponse}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-role [get]
func (h *MenuHandler) ResolveMenuByRoleID(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	req := auth.MenuRequest{
		RoleId: roleID,
	}

	menu, err := h.MenuService.ResolveMenuByRoleID(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// ResolveMenuByRoleIDTransaksi adalah untuk mendapatkan semua menu berdasarkan RoleID.
// @Summary Mendapatkan semua data Menu.
// @Description Endpoint ini adalah untuk mendapatkan Role dan semua menu berdasarkan RoleID.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param roleId query string true "Set RoleID"
// @Success 200 {object} response.Base{data=auth.MenuResponse}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-role/trx [get]
func (h *MenuHandler) ResolveMenuByRoleIDTrx(w http.ResponseWriter, r *http.Request) {
	roleID := r.URL.Query().Get("roleId")
	req := auth.MenuRequest{
		RoleId: roleID,
	}

	menu, err := h.MenuService.ResolveMenuByRoleIDTrx(req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, menu)
}

// UpdateMenuPermission adakan untuk update Menu permission
// @Summary Sort data menu
// @Description endpoint ini adalah update Menu permission
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param Menu body auth.RequestMenuPermissionFormat true "Menu Permission"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-role/update-permission [put]
func (h *MenuHandler) UpdateMenuPermission(w http.ResponseWriter, r *http.Request) {
	var menu auth.RequestMenuPermissionFormat
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = h.MenuService.UpdateMenuPermission(menu)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, "success")
}

// CreateMenuRole adalah untuk menambah data bulk Menu Role.
// @Summary menambahkan data bulk MenuUser.
// @Description Endpoint ini adalah untuk menambahkan data Bulk Menu Role.
// @Tags Menus
// @Produce json
// @Security BearerAuth
// @Param MenuUser body auth.RequestBulkMenuRole true "Menu yang akan ditambahkan"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/menu-role/bulk [post]
func (h *MenuHandler) CreateBulkMenuRole(w http.ResponseWriter, r *http.Request) {
	var reqFormat auth.RequestBulkMenuRole
	err := json.NewDecoder(r.Body).Decode(&reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	newMenuUser, err := h.MenuService.CreateBulkMenuRole(reqFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, newMenuUser)
}
