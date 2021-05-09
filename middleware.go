package main

import (
	"file_share/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func IdentityHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("UUID")
		if err != nil {
			SetNewIdentity(c, db)
			return
		}
	
		var uid uuid.UUID
		if err := (&uid).UnmarshalText([]byte(cookie)); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var i models.Identity
		if err := db.First(&i, uid).Error; err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				SetNewIdentity(c, db)
			default:
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
		c.Set("identity", &i)
		c.String(http.StatusOK, "Old identity value: %s \n", &i)
	}
}

func SetNewIdentity(c *gin.Context, db *gorm.DB) {
	i := models.Identity{IP: c.ClientIP(), UserAgent: c.GetHeader("User-Agent")}
	if err := db.Create(&i).Error; err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.SetCookie("UUID", i.UUID.String(), 3600, "/", "localhost", false, true)
	c.Set("identity", &i)
	c.String(http.StatusOK, "New Identity value: %s \n", &i)
}
