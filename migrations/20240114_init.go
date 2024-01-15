package migrations

import (
	"backend/api/models"
	"backend/config"
)

func RunMigrations() {
	db := config.GetDB()

	db.AutoMigrate(&models.List{}, &models.Sublist{}, &models.File{})
}
