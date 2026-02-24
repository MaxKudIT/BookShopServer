package fav

import (
	"context"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (fh *favHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var favItemdt struct {
		Title string
	}

	if err := c.ShouldBindJSON(&favItemdt); err != nil {
		fh.l.Error("Error creating fav: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		fh.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	fav := dto.FavToDomain(uuid.New(), time.Now())

	favId, err := fh.fserv.Create(ctxnew, fav, firebaseid.(string))
	if err != nil {
		fh.l.Error("Error creating fav", "id", favId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fh.l.Info("Successfully created fav", "id", favId)
	fh.l.Info(favItemdt.Title)
	c.JSON(201, gin.H{"id": favId})
}
