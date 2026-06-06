package ai_chat

import (
	"context"

	"github.com/bookshop/internal/transport/web/middleware"
	"github.com/gin-gonic/gin"
)

func (acr *aiChatRouter) AIChatRegRouters(ctx context.Context, gr *gin.RouterGroup) {
	aiChats := gr.Group("/ai-chats")
	{
		aiChats.POST("", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			acr.ach.CreateChat(c.Request.Context(), c)
		})
		aiChats.POST("/:chatId/messages", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			acr.ach.CreateMessage(c.Request.Context(), c)
		})
		aiChats.GET("/:chatId/messages", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			acr.ach.Messages(c.Request.Context(), c)
		})
		aiChats.DELETE("/:chatId/messages", middleware.VerifyTokenMiddleware(), func(c *gin.Context) {
			acr.ach.DeleteMessages(c.Request.Context(), c)
		})
	}
}
