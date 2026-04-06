package files

import (
	"fmt"
	"io"
	"log"
	"lms-be/configs"
	"lms-be/infras"
	"lms-be/shared/failure"
	"lms-be/transport/http/response"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

type FileServiceImpl struct {
	DB     *infras.PostgresqlConn
	Config *configs.Config
}

func ProvideFileServiceImpl(db *infras.PostgresqlConn, c *configs.Config) *FileServiceImpl {
	return &FileServiceImpl{
		DB:     db,
		Config: c,
	}
}

type FileService interface {
	UploadFile(filePath string, w http.ResponseWriter, r *http.Request) (path string, err error)
}

func (s *FileServiceImpl) UploadFile(filePath string, w http.ResponseWriter, r *http.Request) (path string, err error) {
	if err = r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	newID, _ := uuid.NewV4()
	filename := fmt.Sprintf("%s%s", newID.String(), filepath.Ext(handler.Filename))
	dir := s.Config.App.File.Dir

	path = filepath.Join(filePath, filename)
	fileLocation := filepath.Join(dir, path)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("File " + path + " was Created")

	return
}
