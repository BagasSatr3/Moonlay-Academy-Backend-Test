package routes

import (
	"backend/api/handlers"
	"backend/api/middleware"

	"github.com/labstack/echo/v4"
)

func SublistRoutes(e *echo.Echo) {
	e.GET("/lists/:listID/sublists", handlers.GetSublists)
	e.GET("/sublists/:id", handlers.GetSublistByID)
	e.POST("/lists/:listID/sublists", handlers.CreateSublist, middleware.FileUpload)
	e.PUT("/sublists/:id", handlers.UpdateSublist)
	e.DELETE("/sublists/:id", handlers.DeleteSublist)
}
