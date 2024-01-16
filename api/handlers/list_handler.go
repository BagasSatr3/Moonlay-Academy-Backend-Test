package handlers

import (
	"backend/api/models"
	"backend/config"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetLists(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	titleFilter := c.QueryParam("title")
	descriptionFilter := c.QueryParam("description")

	db := config.GetDB()
	var lists []models.List

	query := db.Model(&lists).Offset(offset).Limit(limit).Preload("Files")

	if titleFilter != "" {
		query = query.Where("title LIKE ?", "%"+titleFilter+"%")
	}
	if descriptionFilter != "" {
		query = query.Where("description LIKE ?", "%"+descriptionFilter+"%")
	}

	if c.QueryParam("preloadSublists") == "true" {
		query = query.Preload("Sublists").Preload("Sublists.Files")
	}

	query.Find(&lists)

	return c.JSON(http.StatusOK, lists)
}

func GetListByID(c echo.Context) error {
	db := config.GetDB()
	id := c.Param("id")
	var list models.List

	if err := db.Preload("Files").Preload("Sublists").Preload("Sublists.Files").First(&list, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "List not found!")
	}

	return c.JSON(http.StatusOK, list)
}

func CreateList(c echo.Context) error {
	db := config.GetDB()
	var list models.List

	if err := c.Bind(&list); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input!")
	}

	if err := list.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	fileNames := c.Get("fileNames").([]string)
	list.Files = make([]models.File, len(fileNames))
	for i, fileName := range fileNames {
		list.Files[i] = models.File{
			FileName: fileName,
			ListID:   &list.ID,
		}
	}

	db.Create(&list)

	return c.JSON(http.StatusOK, list)
}

func UpdateList(c echo.Context) error {
	db := config.GetDB()
	id := c.Param("id")
	var list models.List

	if err := db.First(&list, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "List not found!")
	}

	if err := c.Bind(&list); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input!")
	}

	if err := list.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db.Save(&list)

	return c.JSON(http.StatusOK, list)
}

func DeleteList(c echo.Context) error {
	db := config.GetDB()
	db = db.Debug()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid list ID")
	}

	var list models.List
	if err := db.Preload("Files").Preload("Sublists").First(&list, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "List not found!")
	}

	if len(list.Files) > 0 {
		if err := db.Unscoped().Model(&list).Association("Files").Delete(&list.Files); err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to delete associated files")
		}
	}

	if len(list.Sublists) > 0 {
		if err := db.Unscoped().Model(&list).Association("Sublists").Delete(&list.Sublists); err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to delete associated sublists")
		}
	}

	if err := db.Unscoped().Delete(&list).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete list")
	}

	return c.JSON(http.StatusOK, "List deleted successfully")
}
