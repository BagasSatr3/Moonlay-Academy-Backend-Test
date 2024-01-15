package handlers

import (
	"backend/api/models"
	"backend/config"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetSublists(c echo.Context) error {
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

	listID := c.Param("listID")

	db := config.GetDB()
	var subLists []models.Sublist

	query := db.Model(&subLists).Where("list_id = ?", listID).Offset(offset).Limit(limit).Preload("Files")

	if titleFilter != "" {
		query = query.Where("title LIKE ?", "%"+titleFilter+"%")
	}

	if descriptionFilter != "" {
		query = query.Where("description LIKE ?", "%"+descriptionFilter+"%")
	}

	query.Find(&subLists)

	return c.JSON(http.StatusOK, subLists)
}

func GetSublistByID(c echo.Context) error {
	db := config.GetDB()
	id := c.Param("id")
	var subList models.Sublist

	if err := db.Preload("Files").First(&subList, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "SubList not found!")
	}

	return c.JSON(http.StatusOK, subList)
}

func CreateSublist(c echo.Context) error {
	db := config.GetDB()
	var list models.List
	var subList models.Sublist

	id := c.Param("listID")

	if err := c.Bind(&subList); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input!")
	}

	if err := subList.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := db.First(&list, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "List not found!")
	}

	subList.ListID = list.ID

	fileNames := c.Get("fileNames").([]string)
	subList.Files = make([]models.File, len(fileNames))
	for i, fileName := range fileNames {
		subList.Files[i] = models.File{
			FileName:  fileName,
			SublistID: &subList.ID,
		}
	}

	db.Create(&subList)

	return c.JSON(http.StatusOK, subList)
}

func UpdateSublist(c echo.Context) error {
	db := config.GetDB()
	id := c.Param("id")
	var subList models.Sublist

	if err := db.First(&subList, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "SubList not found!")
	}

	if err := c.Bind(&subList); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input!")
	}

	if err := subList.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	db.Save(&subList)

	return c.JSON(http.StatusOK, subList)
}

func DeleteSublist(c echo.Context) error {
	db := config.GetDB()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid list ID")
	}

	var subList models.Sublist
	if err := db.Preload("Files").First(&subList, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "SubList not found!")
	}

	for _, file := range subList.Files {
		if err := db.Delete(&file).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to delete associated files")
		}
	}

	if err := db.Delete(&subList).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete subList")
	}

	return c.JSON(http.StatusOK, "List deleted successfully")
}
