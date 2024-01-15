package routes

import (
	"backend/api/handlers"
	"backend/api/middleware"

	"github.com/labstack/echo/v4"
)

func ListRoutes(e *echo.Echo) {
	e.GET("/lists", handlers.GetLists)
	e.GET("/lists/:id", handlers.GetListByID)
	e.POST("/lists", handlers.CreateList, middleware.FileUpload)
	e.PUT("/lists/:id", handlers.UpdateList)
	e.DELETE("/lists/:id", handlers.DeleteList)
}
