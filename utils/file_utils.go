package utils

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func IsValidFileType(file []byte) bool {
	fileType := http.DetectContentType(file)
	return strings.HasPrefix(fileType, "image/")
}

func CreateFile(filename string) (*os.File, error) {
	//create uploads directory if not peresent
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	// build file path and create it
	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return nil, err
	}

	return dst, nil
}
