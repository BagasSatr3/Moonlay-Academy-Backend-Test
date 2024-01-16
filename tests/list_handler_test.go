package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/api/handlers"
	"backend/api/middleware"
	"backend/api/models"
	"backend/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetLists(t *testing.T) {
	e := echo.New()
	config.InitDB(config.PgConfig)
	defer config.CloseDB()

	db := config.GetDB()
	log := config.GetDB().Logger
	sampleList := models.List{
		Title:       "Sample List",
		Description: "This is a sample list.",
		Files:       []models.File{{FileName: "sample.txt"}},
	}
	db.Create(&sampleList)

	fmt.Println(log)

	req := httptest.NewRequest(http.MethodGet, "/lists", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.GetLists(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.List
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestGetListByID(t *testing.T) {
	e := echo.New()
	config.InitDB(config.PgConfig)
	defer config.CloseDB()

	db := config.GetDB()
	sampleList := models.List{
		Title:       "Sample List",
		Description: "This is a sample list.",
		Files:       []models.File{{FileName: "sample.txt"}},
	}
	db.Create(&sampleList)

	var lastInsertedID int
	db.Table("lists").Select("id").Order("id desc").Limit(1).Row().Scan(&lastInsertedID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/lists/%d", lastInsertedID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	if lastInsertedID > 0 {
		c.SetParamValues(fmt.Sprintf("%d", lastInsertedID))
	} else {
		c.SetParamValues("1")
	}

	if assert.NoError(t, handlers.GetListByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.List
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestCreateList(t *testing.T) {
	e := echo.New()
	config.InitDB(config.PgConfig)
	defer config.CloseDB()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	err := writer.WriteField("title", "Test List")
	assert.NoError(t, err)
	err = writer.WriteField("description", "This is a test list.")
	assert.NoError(t, err)

	fileContents := []byte("test file content")
	err = writer.WriteField("file", "test.txt")
	assert.NoError(t, err)
	part, err := writer.CreateFormFile("file", "test.txt")
	assert.NoError(t, err)
	_, err = part.Write(fileContents)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/lists", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middleware.FileUpload(handlers.CreateList)(c)

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var response models.List
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestUpdateList(t *testing.T) {
	e := echo.New()
	config.InitDB(config.PgConfig)
	defer config.CloseDB()

	db := config.GetDB()
	sampleList := models.List{
		Title:       "Sample List",
		Description: "This is a sample list.",
		Files:       []models.File{{FileName: "sample.txt"}},
	}
	db.Create(&sampleList)

	var lastInsertedID int
	db.Table("lists").Select("id").Order("id desc").Limit(1).Row().Scan(&lastInsertedID)

	updatedData := models.List{
		Title:       "Updated List",
		Description: "This list has been updated.",
	}

	body, err := json.Marshal(updatedData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/lists/%d", lastInsertedID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	if lastInsertedID > 0 {
		c.SetParamValues(fmt.Sprintf("%d", lastInsertedID))
	} else {
		c.SetParamValues("1")
	}

	if assert.NoError(t, handlers.UpdateList(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.List
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestDeleteList(t *testing.T) {
	e := echo.New()
	config.InitDB(config.PgConfig)
	defer config.CloseDB()

	db := config.GetDB()
	sampleList := models.List{
		Title:       "Sample List",
		Description: "This is a sample list.",
		Files:       []models.File{{FileName: "sample.txt"}},
	}

	db.Create(&sampleList)

	var lastInsertedID int
	db.Table("lists").Select("id").Order("id desc").Limit(1).Row().Scan(&lastInsertedID)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/lists/%d", lastInsertedID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	if lastInsertedID > 0 {
		c.SetParamValues(fmt.Sprintf("%d", lastInsertedID))
	} else {
		c.SetParamValues("1")
	}

	if assert.NoError(t, handlers.DeleteList(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "List deleted successfully", response)
	}
}
