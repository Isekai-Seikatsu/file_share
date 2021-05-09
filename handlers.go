package main

import (
	"file_share/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddBucketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		bucket := models.NewBucket(identity.UUID, 60*time.Second)
		if err := db.Create(bucket).Error; err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.String(http.StatusOK, "Create Bucket %s", bucket)
	}
}

func ListBucketHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		var buckets []models.Bucket
		if err := db.Model(identity).Association("Buckets").Find(&buckets); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		for i, buck := range buckets {
			c.String(http.StatusOK, "Num %s Bucket %s\n", i, buck.ID)
		}
	}
}
