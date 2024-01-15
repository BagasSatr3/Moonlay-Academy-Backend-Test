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

		// Create a directory to store the uploaded files
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return err
		}

		// Process each uploaded file
		var fileNames []string
		for _, file := range files {
			// Validate file size
			if file.Size > (10 * 1024 * 1024) { // 10 MB
				return c.JSON(http.StatusBadRequest, "File size exceeds the limit.")
			}

			// Validate file type (allow only txt and pdf)
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".txt" && ext != ".pdf" {
				return c.JSON(http.StatusBadRequest, "Invalid file type. Only txt and pdf are allowed.")
			}

			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			// Generate a unique file name using uuid
			randomFileName := uuid.New().String() + filepath.Ext(file.Filename)
			dst, err := os.Create(filepath.Join(uploadDir, randomFileName))
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy the file to the destination
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

			fileNames = append(fileNames, randomFileName)
		}

		// Attach file information to the context
		c.Set("fileNames", fileNames)

		// Continue to the next handler
		return next(c)
	}
}
