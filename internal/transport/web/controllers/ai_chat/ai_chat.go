package ai_chat

import (
	"context"
	"net/http"
	"time"

	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (ach *aiChatHandler) CreateChat(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var aiChatDTO dto.AIChatDTO
	if err := c.ShouldBindJSON(&aiChatDTO); err != nil {
		ach.l.Error("Error creating ai chat: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ach.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	aiChat := dto.AIChatToDomain(uuid.New(), aiChatDTO, time.Now())
	aiChatId, err := ach.acs.CreateChat(ctxnew, aiChat, firebaseid.(string))
	if err != nil {
		ach.l.Error("Error creating ai chat", "id", aiChatId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ach.l.Info("Successfully created ai chat", "id", aiChatId)
	c.JSON(http.StatusCreated, gin.H{"id": aiChatId})
}

func (ach *aiChatHandler) CreateMessage(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	chatId, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var aiMessageDTO dto.AIMessageDTO
	if err := c.ShouldBindJSON(&aiMessageDTO); err != nil {
		ach.l.Error("Error creating ai message: data not valid")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ach.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	aiMessage := dto.AIMessageToDomain(uuid.New(), chatId, aiMessageDTO, time.Now())
	aiMessageId, err := ach.acs.CreateMessage(ctxnew, aiMessage, firebaseid.(string))
	if err != nil {
		ach.l.Error("Error creating ai message", "id", aiMessageId, "err", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ach.l.Info("Successfully created ai message", "id", aiMessageId)
	c.JSON(http.StatusCreated, gin.H{"id": aiMessageId})
}

func (ach *aiChatHandler) Messages(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	chatId, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ach.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	aiMessages, err := ach.acs.Messages(ctxnew, firebaseid.(string), chatId)
	if err != nil {
		ach.l.Error("Error getting ai messages", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ach.l.Info("Successfully got ai messages", "chatId", chatId)
	c.JSON(http.StatusOK, gin.H{"messages": aiMessages})
}

func (ach *aiChatHandler) DeleteMessages(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	chatId, err := uuid.Parse(c.Param("chatId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	firebaseid, exists := c.Get("firebase_id")
	if !exists {
		ach.l.Error("firebaseid not found in context")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := ach.acs.DeleteMessages(ctxnew, firebaseid.(string), chatId); err != nil {
		ach.l.Error("Error deleting ai messages", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ach.l.Info("Successfully deleted ai messages", "chatId", chatId)
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
