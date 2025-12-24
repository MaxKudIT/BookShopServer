package user

import (
	"context"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/bookshop/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (uh *userHandler) Create(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var userdt dto.UserDTO

	if err := c.ShouldBindJSON(&userdt); err != nil {
		uh.l.Error("Error creating user: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := dto.Validate(userdt); err != nil {
		uh.l.Error("Error creating user: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	passwordHash, err := utils.GenerateHash(userdt.Password)
	if err != nil {
		uh.l.Error("Error creating user: password hash not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userp := dto.UserToDomain(uuid.New(), passwordHash, userdt)
	user, err := uh.us.Create(ctxnew, userp)
	if err != nil {
		uh.l.Error("Error creating user", "id", user.Id, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	uh.l.Info("Successfully created user", "id", user.Id)
	c.JSON(201, gin.H{"id:": user.Id})
}

func (uh *userHandler) Delete(ctx context.Context, c *gin.Context) {

	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	if err := uh.us.Delete(ctxnew, uuid); err != nil {
		uh.l.Error("Error deleting user", "id", id, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
		return
	}
	uh.l.Info("Successfully deleted user", "id", id)
}
