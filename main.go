package main

import (
	"file_share/models"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func main() {
	dsn := fmt.Sprintf(
		"host=localhost user=postgres password=%s dbname=postgres port=5432",
		os.Getenv("PG_PASS"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Identity{})

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		HomeHandler(c, db)
	})

	router.Run()
}
