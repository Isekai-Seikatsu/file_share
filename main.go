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

	db.AutoMigrate(&models.Identity{}, &models.Bucket{})

	mc, err := GetMinioClient()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.Use(IdentityHandler(db))

	router.GET("/bucket", BucketHandler(db, mc))
	router.GET("/bucket/:buckid", ListBucketHandler(db, mc))

	router.PUT("/bucket/:buckid/:filename", RedirectwUploadLink(db, mc))

	router.Run()
}
