package files

import (
	"io"
	"os"

	"lms-be/shared/random"
)

const (
	// TemporaryFileFolder is a temporary folder file
	TemporaryFileFolder = "temp"
)

// CreateEmptyFile is create empty file (temporary)
func CreateEmptyFile() (fileName string, file *os.File, err error) {
	randStr, err := random.RandStringBytes(5, random.LettersOnly)
	if err != nil {
		return
	}
	fileName = TemporaryFileFolder + "/" + randStr + ".csv"
	file, err = os.Create(fileName)
	if err != nil {
		return
	}

	err = os.Chmod(fileName, 777)
	if err != nil {
		return
	}

	return
}

// CopyFile is copy file to empty file (temporary)
func CopyFile(dst io.Writer, src io.Reader) (err error) {

	_, err = io.Copy(dst, src)
	if err != nil {
		return
	}

	return
}
