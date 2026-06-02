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

	isInCart, err := cih.ciserv.IsInCart(ctxnew, firebaseid.(string), cartItemdt.PhysicalBookId)
	if err != nil {
		cih.l.Error("Error getting result about physical book in the cart item", "id", cartItemdt.PhysicalBookId, "err", err)
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

	isInCart, err := cih.ciserv.AreAllInCart(ctxnew, firebaseid.(string), cartItemsdt.PhysicalBookIds)
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

	cartItem := dto.CartItemToDomain(time.Now(), cartItemdt.PhysicalBookId)

	_, err := cih.ciserv.Create(ctxnew, firebaseid.(string), cartItem)
	if err != nil {
		cih.l.Error("Error creating cart item", "id", cartItem.PhysicalBookId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully created cart item", "id", cartItem.PhysicalBookId)
	c.JSON(201, gin.H{"id": cartItem.PhysicalBookId})
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

	var cartItems = make([]domain.CartItem, 0, len(cartItemsdt.PhysicalBookIds))
	for _, physicalBookId := range cartItemsdt.PhysicalBookIds {
		cartItem := dto.CartItemToDomain(time.Now(), physicalBookId)
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

	physicalBookIdsStr := c.QueryArray("id")

	physicalBookIds := make([]uuid.UUID, 0, len(physicalBookIdsStr))
	for _, idStr := range physicalBookIdsStr {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid UUID: " + idStr})
			return
		}
		physicalBookIds = append(physicalBookIds, id)
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		cih.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := cih.ciserv.Delete(ctxnew, physicalBookIds, firebaseid.(string)); err != nil {
		cih.l.Error("Error deleting cart items", "ids", "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}
	cih.l.Info("Successfully deleted cart items")
}
