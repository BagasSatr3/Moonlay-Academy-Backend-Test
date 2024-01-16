package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/api/handlers"
	"backend/api/middleware"
	"backend/api/models"
	"backend/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetSublists(t *testing.T) {
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

	sampleSublist := models.Sublist{
		Title:       "Sample Sublist",
		Description: "This is a sample sublist.",
		Files:       []models.File{{FileName: "sublist.txt"}},
		ListID:      sampleList.ID,
	}
	db.Create(&sampleSublist)

	req := httptest.NewRequest(http.MethodGet, "/lists/"+strconv.Itoa(int(sampleList.ID))+"/sublists", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("listID")
	c.SetParamValues(strconv.Itoa(int(sampleList.ID)))

	if assert.NoError(t, handlers.GetSublists(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response []models.Sublist
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestGetSublistByID(t *testing.T) {
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

	sampleSublist := models.Sublist{
		Title:       "Sample Sublist",
		Description: "This is a sample sublist.",
		Files:       []models.File{{FileName: "sublist.txt"}},
		ListID:      sampleList.ID,
	}
	db.Create(&sampleSublist)

	req := httptest.NewRequest(http.MethodGet, "/sublists/"+strconv.Itoa(int(sampleSublist.ID)), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(sampleSublist.ID)))

	if assert.NoError(t, handlers.GetSublistByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Sublist
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestCreateSublist(t *testing.T) {
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

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	err := writer.WriteField("title", "Test Sublist")
	assert.NoError(t, err)
	err = writer.WriteField("description", "This is a test sublist.")
	assert.NoError(t, err)

	fileContents := []byte("test file content")
	err = writer.WriteField("file", "sublist.txt")
	assert.NoError(t, err)
	part, err := writer.CreateFormFile("file", "sublist.txt")
	assert.NoError(t, err)
	_, err = part.Write(fileContents)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/lists/"+strconv.Itoa(int(sampleList.ID))+"/sublists", &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("listID")
	c.SetParamValues(strconv.Itoa(int(sampleList.ID)))

	middleware.FileUpload(handlers.CreateSublist)(c)

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var response models.Sublist
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestUpdateSublist(t *testing.T) {
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

	fmt.Println(sampleList.ID)

	sampleSublist := models.Sublist{
		Title:       "Sample Sublist",
		Description: "This is a sample sublist.",
		Files:       []models.File{{FileName: "sublist.txt"}},
		ListID:      sampleList.ID,
	}
	db.Create(&sampleSublist)

	fmt.Println(sampleSublist)

	updatedData := models.Sublist{
		ListID:      sampleSublist.ListID,
		Title:       "Updated Sublist",
		Description: "This sublist has been updated.",
	}

	body, err := json.Marshal(updatedData)
	assert.NoError(t, err)

	fmt.Println("Request Body:", string(body))

	req := httptest.NewRequest(http.MethodPut, "/sublists/"+strconv.Itoa(int(sampleSublist.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	fmt.Println(req)
	fmt.Println(rec)
	c := e.NewContext(req, rec)
	fmt.Println(c)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(sampleSublist.ID)))

	if assert.NoError(t, handlers.UpdateSublist(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response models.Sublist
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
	}
}

func TestDeleteSublist(t *testing.T) {
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

	sampleSublist := models.Sublist{
		Title:       "Sample Sublist",
		Description: "This is a sample sublist.",
		Files:       []models.File{{FileName: "sublist.txt"}},
		ListID:      sampleList.ID,
	}
	db.Create(&sampleSublist)

	req := httptest.NewRequest(http.MethodDelete, "/sublists/"+strconv.Itoa(int(sampleSublist.ID)), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(sampleSublist.ID)))

	if assert.NoError(t, handlers.DeleteSublist(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response string
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Sublist deleted successfully", response)
	}
}
