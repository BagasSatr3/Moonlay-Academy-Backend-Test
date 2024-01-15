package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/api/handlers"
	"backend/api/models"
	"backend/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetLists(t *testing.T) {
	// Setup
	e := echo.New()
	config.InitDB()
	defer config.CloseDB()

	// Insert sample data for testing
	db := config.GetDB()
	sampleList := models.List{
		Title:       "Sample List",
		Description: "This is a sample list.",
		Files:       []models.File{{FileName: "sample.txt"}},
	}
	db.Create(&sampleList)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/lists", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler function
	if assert.NoError(t, handlers.GetLists(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate the response, you might need to adjust this based on your actual response structure
		var response []models.List
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Assert other conditions based on your business logic
	}
}

func TestGetListByID(t *testing.T) {
	// Setup
	e := echo.New()
	config.InitDB()
	defer config.CloseDB()

	// Insert sample data for testing
	db := config.GetDB()
	sampleList := models.List{
		Title:       "Sample List",
		Description: "This is a sample list.",
		Files:       []models.File{{FileName: "sample.txt"}},
	}
	db.Create(&sampleList)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/lists/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Call the handler function
	if assert.NoError(t, handlers.GetListByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		// Validate the response, you might need to adjust this based on your actual response structure
		var response models.List
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Assert other conditions based on your business logic
	}
}
