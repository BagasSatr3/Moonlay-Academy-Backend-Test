package middleware

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func FileUpload(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		files := form.File["file"]

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return err
		}

		var fileNames []string
		for _, file := range files {
			if file.Size > (10 * 1024 * 1024) { // 10 MB
				return c.JSON(http.StatusBadRequest, "File size exceeds the limit.")
			}

			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".txt" && ext != ".pdf" {
				return c.JSON(http.StatusBadRequest, "Invalid file type. Only txt and pdf are allowed.")
			}

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			randomFileName := uuid.New().String() + filepath.Ext(file.Filename)
			dst, err := os.Create(filepath.Join(uploadDir, randomFileName))
			if err != nil {
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			fileNames = append(fileNames, randomFileName)
		}

		c.Set("fileNames", fileNames)

		return next(c)
	}
}
