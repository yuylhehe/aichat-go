package router

import (
	"ai-chat/assets"
	"ai-chat/internal/handler"
	"ai-chat/internal/middleware"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Router HTTP路由配置
type Router struct {
	engine    *gin.Engine
	jwtSecret string

	// 处理器
	authHandler         *handler.AuthHandler
	aiHandler           *handler.AIHandler
	conversationHandler *handler.ConversationHandler
	messageHandler      *handler.MessageHandler
	fixedPromptHandler  *handler.FixedPromptHandler
	userHandler         *handler.UserHandler
}

// RouterConfig 路由配置
type RouterConfig struct {
	JWTSecret           string
	AuthHandler         *handler.AuthHandler
	AIHandler           *handler.AIHandler
	ConversationHandler *handler.ConversationHandler
	MessageHandler      *handler.MessageHandler
	FixedPromptHandler  *handler.FixedPromptHandler
	UserHandler         *handler.UserHandler
}

// NewRouter 创建路由
func NewRouter(config *RouterConfig) *Router {
	r := &Router{
		engine:    gin.Default(),
		jwtSecret: config.JWTSecret,

		authHandler:         config.AuthHandler,
		aiHandler:           config.AIHandler,
		conversationHandler: config.ConversationHandler,
		messageHandler:      config.MessageHandler,
		fixedPromptHandler:  config.FixedPromptHandler,
		userHandler:         config.UserHandler,
	}

	r.setupRoutes()
	return r
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {

	publicFS, err := fs.Sub(assets.PublicFS, "public")
	if err != nil {
		panic(err)
	}

	httpFS := http.FS(publicFS)

	// 静态文件服务
	r.engine.StaticFS("/static", httpFS)

	// 根路径 - 返回前端页面
	r.engine.GET("/", func(c *gin.Context) {
		content, err := assets.PublicFS.ReadFile("public/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})

	// API版本组
	v1 := r.engine.Group("/api/v1")
	{
		// 认证路由
		auth := v1.Group("/auth")
		{
			// 公网环境不注册
			//auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.GET("/me",
				middleware.Auth(r.jwtSecret), r.authHandler.GetProfile)
		}

		// AI对话路由
		ai := v1.Group("/ai")
		ai.Use(middleware.Auth(r.jwtSecret))
		{
			ai.POST("/chat", r.aiHandler.SendMessage)
			ai.POST("/stream", r.aiHandler.StreamChat)
			ai.GET("/stream/:conversationId", r.aiHandler.StreamChatByConversationID)
			// TODO: 模型可选
			ai.GET("/models", r.aiHandler.GetModels)
		}

		// 对话路由
		conversations := v1.Group("/conversations")
		conversations.Use(middleware.Auth(r.jwtSecret))
		{
			conversations.POST("", r.conversationHandler.Create)
			conversations.GET("", r.conversationHandler.GetList)
			conversations.GET("/:id", r.conversationHandler.GetByID)
			conversations.PUT("/:id", r.conversationHandler.Update)
			conversations.DELETE("/:id", r.conversationHandler.Delete)
		}

		// 消息路由
		messages := v1.Group("/messages")
		messages.Use(middleware.Auth(r.jwtSecret))
		{
			messages.POST("", r.messageHandler.Create)
			messages.GET("", r.messageHandler.GetList)
			messages.GET("/conversation/:conversation_id", r.messageHandler.GetByConversationID)
			messages.GET("/:id", r.messageHandler.GetByID)
			messages.PUT("/:id", r.messageHandler.Update)
			messages.DELETE("/:id", r.messageHandler.Delete)
		}

		// 固定提示词路由
		fixedPrompts := v1.Group("/fixed-prompts")
		fixedPrompts.Use(middleware.Auth(r.jwtSecret))
		{
			fixedPrompts.POST("", r.fixedPromptHandler.Create)
			fixedPrompts.GET("", r.fixedPromptHandler.GetList)
			fixedPrompts.GET("/:id", r.fixedPromptHandler.GetByID)
			fixedPrompts.PUT("/:id", r.fixedPromptHandler.Update)
			fixedPrompts.DELETE("/:id", r.fixedPromptHandler.Delete)
		}

		// 用户路由
		users := v1.Group("/users")
		users.Use(middleware.Auth(r.jwtSecret))
		{
			users.GET("/profile", r.userHandler.GetProfile)
			users.PUT("/profile", r.userHandler.UpdateProfile)
			// users.PUT("/password", r.userHandler.UpdatePassword)
			users.DELETE("/account", r.userHandler.DeleteAccount)

			// 管理员路由，暂时不提供
			// admin := users.Group("/admin")
			// {
			// admin.GET("", r.userHandler.GetUserList)
			// admin.GET("/:id", r.userHandler.GetUserByID)
			// }
		}
	}

	// 健康检查路由
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	
	// 统一 404 处理
	r.engine.NoRoute(func(c *gin.Context) {
		// API 路径返回 JSON 格式错误
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Not Found",
			})
			return
		}

		// 其他路径尝试返回前端入口文件（支持 SPA History 模式）
		// 如果是静态资源路径，则直接返回 404
		if strings.HasPrefix(c.Request.URL.Path, "/static") {
			c.String(http.StatusNotFound, "404 page not found")
			return
		}

		// 返回 index.html
		content, err := assets.PublicFS.ReadFile("public/index.html")
		if err != nil {
			c.String(http.StatusNotFound, "404 page not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})
}

// Engine 获取Gin引擎实例
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
