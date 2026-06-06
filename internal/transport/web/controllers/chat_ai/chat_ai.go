package chat_ai

import (
	"context"
	"github.com/bookshop/internal/transport/web/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (h *Handler) Ask(ctx context.Context, c *gin.Context) {
	ctxnew, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var askDTO dto.AskDTO
	if err := c.ShouldBindJSON(&askDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.service.Ask(ctxnew, askDTO.Question)
	if err != nil {
		h.l.Error("failed to ask ai", "error", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AskResponseDTO{Answer: answer})
}
