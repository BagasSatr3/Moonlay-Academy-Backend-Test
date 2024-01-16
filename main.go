package main

import (
	"backend/api/routes"
	"backend/config"
	"backend/migrations"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize Database
	config.InitDB(config.PgConfig)
	defer config.CloseDB()

	// Run Migrations
	migrations.RunMigrations()

	// Initialize Routes
	routes.ListRoutes(e)
	routes.SublistRoutes(e)

	// Start the server
	e.Start(":8000")
}
