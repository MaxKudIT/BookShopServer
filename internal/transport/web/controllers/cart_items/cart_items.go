package cart_items

import (
	"context"
	"github.com/bookshop/internal/domain"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (cih *cartItemsHandler) IsInCart(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var cartItemdt dto.CartItemDTO

	if err := c.ShouldBindJSON(&cartItemdt); err != nil {
		cih.l.Error("Error getting result about book in the cart item: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isInCart, err := cih.ciserv.IsInCart(ctxnew, firebaseid.(string), cartItemdt.BookId)
	if err != nil {
		cih.l.Error("Error getting result about book in the cart item", "id", cartItemdt.BookId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully getting result about book in the cart item", isInCart)

	c.JSON(201, gin.H{"isInCart": isInCart})
}

func (cih *cartItemsHandler) AreAllInCart(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var cartItemsdt dto.CartItemsDTO

	if err := c.ShouldBindJSON(&cartItemsdt); err != nil {
		cih.l.Error("Error getting result about books in the cart item: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	isInCart, err := cih.ciserv.AreAllInCart(ctxnew, firebaseid.(string), cartItemsdt.BookIds)
	if err != nil {
		cih.l.Error("Error getting result about books in the cart item", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully getting result about books in the cart item", isInCart)

	c.JSON(201, gin.H{"areAllInCart": isInCart})
}

func (cih *cartItemsHandler) AllCartItems(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cartItemsPreview, err := cih.ciserv.AllCartItems(ctxnew, firebaseid.(string))
	if err != nil {
		cih.l.Error("error in GetAllCartItems", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	cih.l.Info("Successfully got cart items")
	c.JSON(http.StatusOK, gin.H{"cartItems": cartItemsPreview})
}

func (cih *cartItemsHandler) Count(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := cih.ciserv.Count(ctxnew, firebaseid.(string))
	if err != nil {
		cih.l.Error("error about count in the cart", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	cih.l.Info("Successfully got cart items count")
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func (cih *cartItemsHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var cartItemdt dto.CartItemDTO

	if err := c.ShouldBindJSON(&cartItemdt); err != nil {
		cih.l.Error("Error creating cart item: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cartItem := dto.CartItemToDomain(time.Now(), cartItemdt.BookId)

	_, err := cih.ciserv.Create(ctxnew, firebaseid.(string), cartItem)
	if err != nil {
		cih.l.Error("Error creating cart item", "id", cartItem.BookId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully created cart item", "id", cartItem.BookId)
	c.JSON(201, gin.H{"id": cartItem.BookId})
}

func (cih *cartItemsHandler) CreateItems(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var cartItemsdt dto.CartItemsDTO

	if err := c.ShouldBindJSON(&cartItemsdt); err != nil {
		cih.l.Error("Error creating cart items: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var cartItems = make([]domain.CartItem, 0, len(cartItemsdt.BookIds))
	for _, bookId := range cartItemsdt.BookIds {
		cartItem := dto.CartItemToDomain(time.Now(), bookId)
		cartItems = append(cartItems, cartItem)
	}

	cartId, err := cih.ciserv.CreateItems(ctxnew, firebaseid.(string), cartItems)
	if err != nil {
		cih.l.Error("Error creating cart items", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully created cart items")
	c.JSON(201, gin.H{"id": cartId})
}

func (cih *cartItemsHandler) Delete(ctx context.Context, c *gin.Context) {

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
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := cih.ciserv.Delete(ctxnew, bookIds, firebaseid.(string)); err != nil {
		cih.l.Error("Error deleting cart items", "ids", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully deleted cart items")
}
