package fav_items

import (
	"context"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (fih *favItemsHandler) IsInFavs(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var favItemdt dto.FavItemDTO

	if err := c.ShouldBindJSON(&favItemdt); err != nil {
		fih.l.Error("Error getting result about book in the fav item: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		fih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isInFavs, err := fih.fiserv.IsInFavs(ctxnew, firebaseid.(string), favItemdt.BookId)
	if err != nil {
		fih.l.Error("Error getting result about book in the fav item", "id", favItemdt.BookId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fih.l.Info("Successfully getting result about book in the fav item", isInFavs)

	c.JSON(201, gin.H{"isInFavs": isInFavs})
}

func (fih *favItemsHandler) Count(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		fih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := fih.fiserv.Count(ctxnew, firebaseid.(string))
	if err != nil {
		fih.l.Error("error about count in the fav", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	fih.l.Info("Successfully got fav items count")
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (fih *favItemsHandler) AllFavsItems(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		fih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	favItemsPreview, err := fih.fiserv.AllFavsItems(ctxnew, firebaseid.(string))
	if err != nil {
		fih.l.Error("error in GetAllFavItems", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	fih.l.Info("Successfully got fav items")
	c.JSON(http.StatusOK, gin.H{"favItems": favItemsPreview})
}

func (fih *favItemsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var favItemdt dto.FavItemDTO

	if err := c.ShouldBindJSON(&favItemdt); err != nil {
		fih.l.Error("Error creating fav item: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		fih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	favItem := dto.FavItemToDomain(time.Now(), favItemdt.BookId)

	_, err := fih.fiserv.Create(ctxnew, firebaseid.(string), favItem)
	if err != nil {
		fih.l.Error("Error creating fav item", "id", favItem.BookId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fih.l.Info("Successfully created fav item", "id", favItem.BookId)
	c.JSON(201, gin.H{"id": favItem.BookId})
}

func (fih *favItemsHandler) Delete(ctx context.Context, c *gin.Context) {

	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	bookIdsStr := c.QueryArray("id")

	bookIds := make([]uuid.UUID, 0, len(bookIdsStr))
	for _, idStr := range bookIdsStr {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid UUID: " + idStr})
			return
		}
		bookIds = append(bookIds, id)
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		fih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := fih.fiserv.Delete(ctxnew, bookIds, firebaseid.(string)); err != nil {
		fih.l.Error("Error deleting fav items", "ids", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}
	fih.l.Info("Successfully deleted fav items")
}
