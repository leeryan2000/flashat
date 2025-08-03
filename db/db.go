package db

import (
	
	"github.com/leeryan2000/flashat/config"
	"github.com/leeryan2000/flashat/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var configuration = config.GetConfig()

func InitDB() (*gorm.DB, error) {
	dsn := "host=" + configuration.DB_DBHOST + " "
	dsn += "user=" + configuration.DB_USER + " "
	dsn += "password=" + configuration.DB_PASS + " "
	dsn += "dbname=" + configuration.DB_DBNAME + " "
	dsn += "port=" + configuration.DB_PORT + " "
	dsn += "sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.User{})

	return db, nil
}
