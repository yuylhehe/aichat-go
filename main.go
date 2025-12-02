package main

import (
	"log"

	"ai-chat/config"
	"ai-chat/internal/handler"
	"ai-chat/internal/repository"
	"ai-chat/internal/router"
	"ai-chat/internal/service"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	db, err := repository.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	// 自动迁移数据库
	if err := repository.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化服务层
	userService := service.NewUserService(db)
	authService := service.NewAuthService(db, cfg)
	conversationService := service.NewConversationService(db)
	messageService := service.NewMessageService(db)
	fixedPromptService := service.NewFixedPromptService(db)
	aiService := service.NewAIService(db, cfg)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	conversationHandler := handler.NewConversationHandler(conversationService, messageService)
	messageHandler := handler.NewMessageHandler(messageService)
	fixedPromptHandler := handler.NewFixedPromptHandler(fixedPromptService)
	aiHandler := handler.NewAIHandler(aiService, conversationService, messageService, fixedPromptService)

	// 创建路由配置
	routerConfig := &router.RouterConfig{
		JWTSecret:           cfg.JWTSecret,
		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		ConversationHandler: conversationHandler,
		MessageHandler:      messageHandler,
		FixedPromptHandler:  fixedPromptHandler,
		AIHandler:           aiHandler,
	}

	// 初始化路由
	r := router.NewRouter(routerConfig)

	// 启动服务器
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Engine().Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
