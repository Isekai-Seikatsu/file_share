package main

import (
	"context"
	"file_share/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func BucketHandler(db *gorm.DB, minioClient *minio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)

		var bucket models.Bucket
		tx := db.Where(
			models.Bucket{IdentityUUID: identity.UUID},
		).Attrs(
			models.NewBucket(identity.UUID, 1*time.Hour),
		).FirstOrCreate(&bucket)
		if err := tx.Error; err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		existed, err := minioClient.BucketExists(context.Background(), bucket.ID.String())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !existed {
			if err := minioClient.MakeBucket(context.Background(), bucket.ID.String(), minio.MakeBucketOptions{}); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		c.String(http.StatusOK, "Bucket %s", bucket)
	}
}

func ListBucketHandler(db *gorm.DB, minioClient *minio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		var buckID uuid.UUID
		if err := (&buckID).UnmarshalText([]byte(c.Param("buckid"))); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var bucket models.Bucket
		if err := db.First(&bucket, buckID, identity.UUID).Error; err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		for info := range minioClient.ListObjects(context.Background(), bucket.ID.String(), minio.ListObjectsOptions{}) {
			c.String(http.StatusOK, "%s", info)
		}
	}
}

func RedirectwUploadLink(db *gorm.DB, minioClient *minio.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		filename := c.Param("filename")
		var buckID uuid.UUID
		if err := (&buckID).UnmarshalText([]byte(c.Param("buckid"))); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		existed, err := minioClient.BucketExists(context.Background(), buckID.String())
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if !existed {
			c.AbortWithError(http.StatusNotFound, err)
		} else {
			u, err := GetUploadLink(minioClient, buckID.String(), filename)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			log.Println(u)
			c.Redirect(http.StatusTemporaryRedirect, u.String())
		}
	}
}
