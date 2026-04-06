package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"lms-be/configs"
	"lms-be/internal/files"
	"lms-be/shared/failure"
	"lms-be/transport/http/middleware"
	"lms-be/transport/http/response"

	"github.com/go-chi/chi"
	"github.com/nfnt/resize"
)

type FileHandler struct {
	Config      *configs.Config
	FileService files.FileService
}

func ProvideFileHandler(conf *configs.Config, fs files.FileService) FileHandler {
	return FileHandler{
		Config:      conf,
		FileService: fs,
	}
}

func (h *FileHandler) Router(r chi.Router, middleware *middleware.JWT) {
	r.Route("/files", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", h.ReadFile)
			r.Post("/upload", h.UploadFile)
			r.Get("/image", h.ReadImage)
		})
	})
}

// ReadFile untuk membaca file sesuai dengan path
// @Summary Get file yang sudah pernah di upload
// @Description End point ini digunakan untuk ownload file yang sudah pernah di upload
// @Tags Files
// @Produce json
// @Param path query string true "Nama file atau path yang tersimpan di database"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/files [get]
func (h *FileHandler) ReadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("path")
	dir := h.Config.App.File.Dir

	// Validate path to prevent path traversal attacks
	fileLocation, err := isPathSafe(dir, filename)
	if err != nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	img, err := os.Open(fileLocation)

	if err != nil {
		http.Error(w, "File Not Found", http.StatusInternalServerError)
	}
	defer img.Close()
	// w.Header().Set("Content-Type", "image/png") // <-- set the content-type header
	io.Copy(w, img)
}

// untuk ngambil 1 folder saja
func (h *FileHandler) ReadFileParams(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "path")
	dir := h.Config.App.File.Dir
	fileDir := h.Config.App.File.FotoProfile

	path := filepath.Join(fileDir, filename)

	// Validate path to prevent path traversal attacks
	fileLocation, err := isPathSafe(dir, path)
	if err != nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	img, err := os.Open(fileLocation)

	if err != nil {
		http.Error(w, "File Not Found", http.StatusInternalServerError)
		return
	}
	defer img.Close()
	io.Copy(w, img)
}

// UploadFile untuk upload file attachment
// @Summary Upload file attachment
// @Description End point ini digunakan untuk mengupload file attachment
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File"
// @Param filePath formData string true "Path File"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/files/upload [post]
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	filePath := r.FormValue("filePath")
	path, err := h.FileService.UploadFile(filePath, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response.WithJSON(w, http.StatusOK, path)
}

// ReadImage untuk membaca file sesuai dengan path dan ukuran
// @Summary Get file yang sudah pernah di upload
// @Description End point ini digunakan untuk ownload file yang sudah pernah di upload
// @Tags Files
// @Produce json
// @Param path query string false "Nama file atau path yang tersimpan di database"
// @Param size query string true "Ukuran gambar, misal: size=100x100"
// @Success 200 {object} response.Base
// @Failure 400 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/files/image [get]
func (h *FileHandler) ReadImage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filename := r.URL.Query().Get("path")
	size := query.Get("size")

	dir := h.Config.App.File.Dir

	// Validate path to prevent path traversal attacks
	fileLocation, err := isPathSafe(dir, filename)
	if err != nil {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	img, err := os.Open(fileLocation)

	if err != nil {
		http.Error(w, "File Not Found", http.StatusInternalServerError)
		return
	}
	defer img.Close()
	w.Header().Set("Content-Type", "image/png") // <-- set the content-type header

	if len(size) > 0 {
		width, height, err := parseSize(size)
		if err != nil {
			failure.BadRequest(errors.New("size is not valid"))
			return
		}

		// // Decode the image (from PNG to image.Image):
		if img == nil {
			http.Error(w, "File Not Found", http.StatusInternalServerError)
			return
		}
		src, err := png.Decode(img)
		if err != nil {
			failure.BadRequest(errors.New("size is not valid"))
			return
		}
		m := resize.Resize(uint(width), uint(height), src, resize.Lanczos2)
		buffer := new(bytes.Buffer)
		if err := jpeg.Encode(buffer, m, nil); err != nil {
			log.Println("unable to encode image.")
		}

		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
		}
	} else {
		io.Copy(w, img)
	}
}

func parseSize(size string) (width, height int, err error) {
	fmt.Println("Size ---> ", size)
	dimensions := strings.Split(size, "x")
	switch len(dimensions) {
	case 2:
		width, err = strconv.Atoi(dimensions[0])
		if err != nil {
			return
		}

		height, err = strconv.Atoi(dimensions[1])
		if err != nil {
			return
		}
	case 1:
		width, err = strconv.Atoi(dimensions[0])
		if err != nil {
			return
		}
		height = 0
	default:
		err = errors.New("Invalid Size")
	}
	return
}

// Validate path to prevent path traversal attacks
func isPathSafe(baseDir, requestedPath string) (string, error) {
	if requestedPath == "" {
		return "", errors.New("path is required")
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", errors.New("invalid base directory")
	}

	cleanPath := filepath.Clean(requestedPath)
	if cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) {
		return "", errors.New("invalid path: path traversal detected")
	}

	fullPath := filepath.Join(absBaseDir, cleanPath)
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", errors.New("invalid file path")
	}

	if absFullPath == absBaseDir {
		return "", errors.New("path must specify a file")
	}
	if !strings.HasPrefix(absFullPath, absBaseDir+string(filepath.Separator)) {
		return "", errors.New("access denied: path outside allowed directory")
	}

	return absFullPath, nil
}
