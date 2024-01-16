// config/database.go
package config

import (
	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type DBConfig struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

var PgConfig = DBConfig{
	Driver:   "postgres",
	User:     "postgres",
	Password: "1234",
	Host:     "localhost",
	Port:     "5432",
	DBName:   "moonlay-todo-list",
	SSLMode:  "disable",
}

func InitDB(config DBConfig) *gorm.DB {
	var (
		database *gorm.DB
		err      error
	)

	switch config.Driver {
	case "postgres":
		dsn := "user=" + config.User + " dbname=" + config.DBName + " password=" + config.Password +
			" host=" + config.Host + " port=" + config.Port + " sslmode=" + config.SSLMode
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql":
		// dsn := config.User + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.DBName +
		// 	"?charset=utf8mb4&parseTime=True&loc=Local"
		// database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		panic("Unsupported database driver")
	}

	if err != nil {
		panic("Failed to connect to database!")
	}

	db = database
	return db
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get underlying SQL database")
	}
	sqlDB.Close()
}
