package main

import (
	"file_share/models"
	"fmt"
	"net/http"
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

	db.AutoMigrate(&models.Identity{}, &models.Bucket{})

	router := gin.Default()
	router.Use(IdentityHandler(db))
	
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Home Page")
	})

	router.Run()
}
