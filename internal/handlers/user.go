package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"lms-be/internal/domain/auth"
	"lms-be/shared"
	"lms-be/shared/failure"
	"lms-be/shared/model"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"

	"github.com/dchest/captcha"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

// UserHandler the HTTP handler for User domain.
type UserHandler struct {
	UserService auth.UserService
}

// ProvideUserHandler is the provider for this handler.
func ProvideUserHandler(userService auth.UserService) UserHandler {
	return UserHandler{
		UserService: userService,
	}
}

// Router sets up the router for this domain.
func (u *UserHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/login", u.Login)
		r.Post("/login-web", u.LoginWeb)
		r.Post("/validasi-login", u.ValidasiLogin)
		r.Get("/captcha", u.captchaHandler)
		r.Get("/captcha/image/{captchaId}", u.captchaImageHandler)
		r.Get("/captcha/refresh/{captchaId}", u.captchaRefreshHandler)
		r.Route("/", func(r chi.Router) {
			r.Use(middleware.VerifyToken)
			r.Post("/", u.CreateUser)
			r.Put("/{id}", u.UpdateUser)
			r.Put("/fcm-token/{id}", u.UpdateUserFcmToken)
			r.Delete("/{id}", u.DeleteUser)
			r.Put("/active-status/{id}", u.UpdateActiveStatus)
			r.Get("/all", u.GetAllData)
			r.Get("/", u.ResolveAll)
			r.Get("/{id}", u.ResolveUserById)
			r.Put("/password/{id}", u.ChangePassword)
			r.Put("/password/pw/{id}", u.ChangePassword)
			r.Put("/password/reset/{id}", u.ResetPassword)
			r.Post("/update-foto", u.UpdateFoto)
		})
	})
}

// ValidasiLogin sign in a user
// @Summary sign in a user
// @Description This endpoint sign in a user
// @Tags Users
// @Param users body auth.InputLogin true "The User to be sign in."
// @Produce json
// @Success 201 {object} response.Base{auth.ResponseLogin}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/validasi-login [post]
func (u *UserHandler) ValidasiLogin(w http.ResponseWriter, r *http.Request) {
	var input auth.InputLogin
	fmt.Println("INPUT:", input)
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	resp, exist, err := u.UserService.ValidasiLogin(input)
	if !exist {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, resp)
}

// Login sign in a user
// @Summary sign in a user
// @Description This endpoint sign in a user
// @Tags Users
// @Param users body auth.InputLogin true "The User to be sign in."
// @Produce json
// @Success 201 {object} response.Base{auth.ResponseLogin}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/login [post]
func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input auth.InputLogin
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	fmt.Println("ip address : ", r.Header.Get("x-forwarded-for"))

	resp, _, _, err := u.UserService.Login(input, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if err != nil {
		resJson := map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		response.WithJSON(w, http.StatusOK, resJson)
		return
	}

	response.WithJSON(w, http.StatusOK, resp)
}

// LoginWeb sign in a user for web with captcha
// @Summary sign in a user for web
// @Description This endpoint sign in a user for web with captcha validation
// @Tags Users
// @Param users body auth.InputLoginWeb true "The User to be sign in."
// @Produce json
// @Success 201 {object} response.Base{auth.ResponseLogin}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/login-web [post]
func (u *UserHandler) LoginWeb(w http.ResponseWriter, r *http.Request) {
	var input auth.InputLoginWeb
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if !captcha.VerifyString(input.CaptchaID, input.CaptchaAnswer) {
		newCaptchaId := captcha.New()
		resJson := map[string]interface{}{
			"success":      false,
			"message":      "Captcha salah",
			"newCaptchaId": newCaptchaId,
		}
		response.WithJSON(w, http.StatusBadRequest, resJson)
		return
	}

	fmt.Println("ip address : ", r.Header.Get("x-forwarded-for"))

	loginInput := auth.InputLogin{
		Username: input.Username,
		Password: input.Password,
	}

	resp, _, _, err := u.UserService.Login(loginInput, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if err != nil {
		resJson := map[string]interface{}{
			"success": false,
			"message": err.Error(),
		}
		response.WithJSON(w, http.StatusBadRequest, resJson)
		return
	}

	response.WithJSON(w, http.StatusOK, resp)
}

// CreateUser creates a new user
// @Summary Create a new User.
// @Description This endpoint creates a new User.
// @Tags Users
// @Security BearerAuth
// @Param users body auth.InputUser true "The User to be created."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user [post]
func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var input auth.InputUser
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = input.Validate()
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	exist, err := u.UserService.CreateUser(input, userID, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if exist {
		response.WithError(w, failure.Conflict("register", "user", err.Error()))
		return
	}

	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithMessage(w, http.StatusOK, "success")
}

// ChangePassword update user password
// @Summary update user password
// @Description This endpoint to update user password
// @Tags Users
// @Security BearerAuth
// @Param id path string true "The User identifier."
// @Param users body auth.InputChangePassword true "The User update a new password."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/password/{id} [put]
func (u *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var input auth.InputChangePassword

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.ChangePassword(id, input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Password berhasil diperbarui",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// ResetPassword reset user password
// @Summary reset user password
// @Description This endpoint to reset user password
// @Tags Users
// @Security BearerAuth
// @Param id path string true "The User identifier."
// @Param users body auth.InputChangePassword true "The User reset a new password."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/password/reset/{id} [put]
func (u *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input auth.InputChangePassword

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.ResetPassword(id, input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	payload := map[string]interface{}{
		"success": true,
		"message": "Password berhasil diperbarui",
	}
	response.WithJSON(w, http.StatusOK, payload)
}

// ResolveAll list all user.
// @Summary Get list all user.
// @Description endpoint ini digunakan untuk mendapatkan seluruh data user sesuai dengan filter yang dikirimkan.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param q query string false "Keyword search"
// @Param pageSize query int true "Set pageSize data"
// @Param pageNumber query int true "Set page number"
// @Param sortBy query string false "Set sortBy parameter is one of [id | kode | nama ]"
// @Param sortType query string false "Set sortType with asc or desc"
// @Param roleId query string false "id role"
// @Param active query bool false "active"
// @Success 200 {object} auth.UserDTO
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user [get]
func (h *UserHandler) ResolveAll(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("q")
	pageSizeStr := r.URL.Query().Get("pageSize")
	pageNumberStr := r.URL.Query().Get("pageNumber")
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
	activeStr := r.URL.Query().Get("active")

	roleId := r.URL.Query().Get("roleId")
	active, err := strconv.ParseBool(activeStr)
	if err != nil {
		active = true
	}

	req := model.StandardRequestUser{
		Keyword:    keyword,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		SortBy:     sortBy,
		SortType:   sortType,
		RoleID:     roleId,
		Active:     active,
	}

	err = shared.GetValidator().Struct(req)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	status, err := h.UserService.ResolveAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}

	response.WithJSON(w, http.StatusOK, status)
}

// GetAllData mengambil seluruh data User tanpa pagination.
// @Summary Ambil semua data User
// @Description Endpoint ini digunakan untuk mengambil semua data user tanpa pagination.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Base{data=[]auth.User}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/all [get]
func (h *UserHandler) GetAllData(w http.ResponseWriter, r *http.Request) {
	req := model.StandardRequestUser{}
	status, err := h.UserService.GetAll(r.Context(), req)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, status)
}

// UpdateUser update user data
// @Summary update user data
// @Description This endpoint to update user entity
// @Tags Users
// @Security BearerAuth
// @Param id path string true "The User identifier."
// @Param users body auth.UserUpdateFormat true "The User update data"
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/{id} [put]
func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input auth.UserUpdateFormat

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	err = u.UserService.UpdateUser(id, input, userID, r.Header.Get("x-forwarded-for"), r.UserAgent())
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithMessage(w, http.StatusOK, "success")
}

// UpdateUserFcmToken update data fcm token user
// @Summary update data fcm token user
// @Description This endpoint to update user entity
// @Tags Users
// @Security BearerAuth
// @Param id path string true "The User identifier."
// @Param users body auth.UserUpdateFcmTokenFormat true "The User update Fcm Token data"
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/fcm-token/{id} [put]
func (u *UserHandler) UpdateUserFcmToken(w http.ResponseWriter, r *http.Request) {
	var input auth.UserUpdateFcmTokenFormat

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		err = errors.New("Format input tidak sesuai!")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(input)
	if err != nil {
		err = errors.New("Format input tidak sesuai!")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.UpdateUserFcmToken(id, input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithMessage(w, http.StatusOK, "success")
}

// UpdateUser delete user data
// @Summary delete user data
// @Description This endpoint to delete user entity
// @Tags Users
// @Security BearerAuth
// @Param id path string true "The User identifier."
// @Produce json
// @Success 201 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/{id} [delete]
func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.SoftDelete(id, userID)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// UpdateActiveStatus updates the active status of a user
// @Summary Update user active status
// @Description This endpoint updates the active status (true/false) of a user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Param id path string true "The User identifier."
// @Param body body object{active=bool} true "Active status"
// @Produce json
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/active-status/{id} [put]
func (u *UserHandler) UpdateActiveStatus(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Active bool `json:"active"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	userID, err := uuid.FromString(middleware.GetClaimsValue(r.Context(), "userId").(string))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	id, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = u.UserService.UpdateActiveStatus(id, userID, input.Active)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusOK, "success")
}

// ResolveUserByID adalah untuk mendapatkan satu data User berdasarkan ID.
// @Summary Mendapatkan satu data User berdasarkan ID.
// @Description Endpoint ini adalah untuk mendapatkan User ID.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID"
// @Success 200 {object} response.Base{data=auth.User}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/{id} [get]
func (h *UserHandler) ResolveUserById(w http.ResponseWriter, r *http.Request) {
	ID, err := uuid.FromString(chi.URLParam(r, "id"))
	if err != nil {
		response.WithError(w, err)
		return
	}
	user, err := h.UserService.ResolveUserById(ID)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, user)
}

// UpdateFotoProfile adalah untuk mengupdate foto pegawai.
// @Summary mengupdate data foto pegawai.
// @Description Endpoint ini adalah untuk mengupdate data foto pegawai.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id formData string false "id pegawai"
// @Param file formData file true "Foto Baru"
// @Success 200 {object} response.Base{data=auth.User}
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/user/update-foto [post]
func (u *UserHandler) UpdateFoto(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	userID, err := uuid.FromString(id)
	if err != nil {
		response.WithError(w, err)
		return
	}

	user, err := u.UserService.ResolveUserById(userID)
	if err != nil {
		response.WithError(w, err)
		return
	}

	uploadedFile, _, _ := r.FormFile("file")
	var path string
	if uploadedFile != nil {
		filepath, err := u.UserService.UploadFile(w, r, "")
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		path = filepath
	} else {
		path = ""
	}

	reqFormat := auth.UpdateFotoRequest{
		Id:   userID,
		Foto: path,
	}

	if user.Foto != nil {
		reqFormat.FotoLama = *user.Foto
	}

	data, err := u.UserService.UpdateFoto(reqFormat)
	if err != nil {
		fmt.Print("error response")
		response.WithError(w, failure.BadRequest(err))
		return
	}

	response.WithJSON(w, http.StatusCreated, data)
}

// captchaHandler generates a new captcha ID and returns it along with the image URL.
// @Summary Generate a new captcha
// @Description This endpoint generates a new captcha ID and returns it along with the image URL.
// @Tags Users
// @Produce json
// @Success 200 {object} response.Base{data=map[string]string}
// @Failure 500 {object} response.Base
// @Router /v1/user/captcha [get]
func (u *UserHandler) captchaHandler(w http.ResponseWriter, r *http.Request) {
	captchaId := captcha.New()
	response.WithJSON(w, http.StatusOK, map[string]string{
		"captchaId":  captchaId,
		"captchaUrl": "/user/captcha/image/" + captchaId,
	})
}

// captchaImageHandler serves the captcha image for a given captcha ID.
// @Summary Get captcha image
// @Description This endpoint serves the captcha image for a given captcha ID.
// @Tags Users
// @Param captchaId path string true "Captcha ID"
// @Produce image/png
// @Success 200 {file} binary
// @Failure 404 {object} response.Base
// @Router /v1/user/captcha/image/{captchaId} [get]
func (u *UserHandler) captchaImageHandler(w http.ResponseWriter, r *http.Request) {
	captchaId := chi.URLParam(r, "captchaId")
	w.Header().Set("Content-Type", "image/png")
	if err := captcha.WriteImage(w, captchaId, 240, 80); err != nil {
		http.Error(w, "Captcha not found", http.StatusNotFound)
	}
}

// captchaRefreshHandler reloads the captcha for a given captcha ID.
// @Summary Refresh captcha
// @Description This endpoint reloads the captcha for a given captcha ID.
// @Tags Users
// @Param captchaId path string true "Captcha ID"
// @Produce json
// @Success 200 {object} response.Base{data=map[string]string}
// @Failure 400 {object} response.Base
// @Router /v1/user/captcha/refresh/{captchaId} [get]
func (u *UserHandler) captchaRefreshHandler(w http.ResponseWriter, r *http.Request) {
	captchaId := chi.URLParam(r, "captchaId")
	if captchaId == "" {
		response.WithError(w, failure.BadRequest(errors.New("Captcha ID is required")))
		return
	}
	captcha.Reload(captchaId)
	response.WithJSON(w, http.StatusOK, map[string]string{
		"captchaId":  captchaId,
		"captchaUrl": "/user/captcha/image/" + captchaId,
	})
}
