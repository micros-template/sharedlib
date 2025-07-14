package utils

import (
	"errors"
	"io"
	"mime/multipart"
)

func FileToByte(fileHeader *multipart.FileHeader) ([]byte, error) {
	f, err := fileHeader.Open()
	if err != nil {
		return nil, errors.New("error open file")
	}
	defer f.Close()
	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.New("error read file")
	}
	return fileBytes, nil
}
