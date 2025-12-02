package handler

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AIHandler AI对话处理器
type AIHandler struct {
	aiService           *service.AIService
	conversationService *service.ConversationService
	messageService      *service.MessageService
	fixedPromptService  *service.FixedPromptService
}

// NewAIHandler 创建AI处理器
func NewAIHandler(
	aiService *service.AIService,
	conversationService *service.ConversationService,
	messageService *service.MessageService,
	fixedPromptService *service.FixedPromptService,
) *AIHandler {
	return &AIHandler{
		aiService:           aiService,
		conversationService: conversationService,
		messageService:      messageService,
		fixedPromptService:  fixedPromptService,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	ConversationID *uint    `json:"conversationId,omitempty"`
	Model          *string  `json:"model,omitempty"`
	Temperature    *float64 `json:"temperature,omitempty"`
	Message        string   `json:"message" binding:"required,min=1"`
	FixedPromptID  *uint    `json:"fixedPromptId,omitempty"`
	Thinking       *struct {
		Type string `json:"type"`
	} `json:"thinking,omitempty"`
}

// StreamChatRequest 流式聊天请求
type StreamChatRequest struct {
	ChatRequest
	UseFixedPrompt bool `json:"useFixedPrompt"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message service.Message `json:"message"`
		Finish  string          `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
	PromptFilterResults []PromptFilterResult `json:"prompt_filter_results,omitempty"`
}

// PromptFilterResult 提示过滤结果
type PromptFilterResult struct {
	PromptIndex  int    `json:"prompt_index"`
	FinishReason string `json:"finish_reason"`
}

// SendMessage 发送消息
func (h *AIHandler) SendMessage(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	userID := middleware.GetUserID(c)

	// 创建或获取会话
	var conversationID uint
	if req.ConversationID != nil {
		conversationID = *req.ConversationID
	} else {
		conversationReq := &service.CreateConversationRequest{
			Name:   req.Message[:min(len(req.Message), 50)],
			UserID: userID,
		}
		conversation, err := h.conversationService.Create(conversationReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "创建会话失败",
				"details": err.Error(),
			})
			return
		}
		conversationID = conversation.ID
	}

	// 获取固定提示词内容
	var systemPrompt string
	if req.FixedPromptID != nil {
		fixedPrompt, err := h.fixedPromptService.FindByID(userID, *req.FixedPromptID)
		if err == nil && fixedPrompt.IsActive {
			systemPrompt = fixedPrompt.Content
		}
	}

	// 构建消息列表
	chatMessages, err := h.buildChatMessages(userID, conversationID, systemPrompt, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取消息历史失败",
			"details": err.Error(),
		})
		return
	}

	// 调用AI服务
	chatReq := &service.ChatRequest{
		Messages:    chatMessages,
		Model:       req.Model,
		Temperature: req.Temperature,
	}

	result, err := h.aiService.ChatCompletion(chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "AI服务调用失败",
			"details": err.Error(),
		})
		return
	}

	// 保存用户消息
	userMessage := &service.CreateMessageRequest{
		ConversationID: conversationID,
		Content:        req.Message,
		Type:           "user",
	}
	_, err = h.messageService.Create(userID, userMessage)
	if err != nil {
		// 即使保存消息失败，也返回AI回复
		c.JSON(http.StatusOK, gin.H{
			"data": ChatResponse{
				ID:      "resp_" + time.Now().Format("20060102150405"),
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   *req.Model,
				Choices: result.Choices,
				Usage:   result.Usage,
			},
		})
		return
	}

	// 保存AI回复
	if len(result.Choices) > 0 {
		assistantMessage := &service.CreateMessageRequest{
			ConversationID: conversationID,
			Content:        result.Choices[0].Message.Content,
			Type:           "assistant",
			Model:          req.Model,
		}
		_, err = h.messageService.Create(userID, assistantMessage)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data": ChatResponse{
					ID:      "resp_" + time.Now().Format("20060102150405"),
					Object:  "chat.completion",
					Created: time.Now().Unix(),
					Model:   *req.Model,
					Choices: result.Choices,
					Usage:   result.Usage,
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": ChatResponse{
			ID:      "resp_" + time.Now().Format("20060102150405"),
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   *req.Model,
			Choices: result.Choices,
			Usage:   result.Usage,
		},
	})
}

// StreamChat 流式聊天
func (h *AIHandler) StreamChat(c *gin.Context) {
	var req StreamChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	userID := middleware.GetUserID(c)

	// 创建或获取会话
	var conversationID uint
	if req.ConversationID != nil {
		conversationID = *req.ConversationID
	} else {
		conversationReq := &service.CreateConversationRequest{
			Name:   req.Message[:min(len(req.Message), 50)],
			UserID: userID,
		}
		conversation, err := h.conversationService.Create(conversationReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "创建会话失败",
				"details": err.Error(),
			})
			return
		}
		conversationID = conversation.ID
	}

	// 保存用户消息
	userMessage := &service.CreateMessageRequest{
		ConversationID: conversationID,
		Content:        req.Message,
		Type:           "user",
	}
	_, err := h.messageService.Create(userID, userMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "保存用户消息失败",
			"details": err.Error(),
		})
		return
	}

	// 获取系统提示词
	var systemPrompt string
	if req.UseFixedPrompt && req.FixedPromptID != nil {
		var fixedPrompt *dto.FixedPromptResponse
		fixedPrompt, err = h.fixedPromptService.FindByID(userID, *req.FixedPromptID)
		if err == nil && fixedPrompt.IsActive {
			systemPrompt = fixedPrompt.Content
		}
	}

	// 构建消息列表
	chatMessages, err := h.buildChatMessages(userID, conversationID, systemPrompt, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取消息历史失败",
			"details": err.Error(),
		})
		return
	}

	// 流式处理AI响应
	chatReq := &service.ChatRequest{
		Messages:    chatMessages,
		Model:       req.Model,
		Temperature: req.Temperature,
		Stream:      true,
		Thinking:    req.Thinking,
	}

	h.processStreamResponse(c, userID, conversationID, chatReq, "message")
}

// GetModels 获取可用模型列表
func (h *AIHandler) GetModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": []string{
			"gpt-3.5-turbo",
			"gpt-3.5-turbo-16k",
			"gpt-4",
			"gpt-4-turbo",
			"gpt-4-turbo-preview",
		},
	})
}

// StreamChatByConversationID 根据会话ID进行流式聊天
func (h *AIHandler) StreamChatByConversationID(c *gin.Context) {
	// 获取会话ID
	conversationIDStr := c.Param("conversationId")
	if conversationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的会话ID",
		})
		return
	}

	// 获取提示词（可选）
	prompt := c.Query("prompt")

	userID := middleware.GetUserID(c)

	// 构建消息列表
	chatMessages, err := h.buildChatMessages(userID, uint(conversationID), "", prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取消息历史失败",
			"details": err.Error(),
		})
		return
	}

	// 流式处理AI响应
	chatReq := &service.ChatRequest{
		Messages:    chatMessages,
		Model:       nil,
		Temperature: nil,
		Stream:      true,
	}

	// 从Query中获取thinking
	thinkingType := c.Query("thinking")
	if thinkingType != "" {
		chatReq.Thinking = &struct {
			Type string `json:"type"`
		}{
			Type: thinkingType,
		}
	}

	h.processStreamResponse(c, userID, uint(conversationID), chatReq, "token")
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// buildChatMessages 构建聊天消息列表
func (h *AIHandler) buildChatMessages(userID, conversationID uint, systemPrompt string, currentMessage string) ([]service.Message, error) {
	// 构建消息历史
	messages, err := h.messageService.FindByConversationID(userID, conversationID)
	if err != nil {
		return nil, err
	}

	// 构造消息列表用于OpenAI
	var chatMessages []service.Message
	if systemPrompt != "" {
		chatMessages = append(chatMessages, service.Message{
			Role:    "system",
			Content: systemPrompt,
		})
	}

	// 添加历史消息
	for _, msg := range messages {
		chatMessages = append(chatMessages, service.Message{
			Role:    msg.Type, // 使用消息的实际类型
			Content: msg.Content,
		})
	}

	// 添加当前用户消息
	if currentMessage != "" {
		chatMessages = append(chatMessages, service.Message{
			Role:    "user",
			Content: currentMessage,
		})
	}

	return chatMessages, nil
}

// processStreamResponse 处理流式响应通用逻辑
func (h *AIHandler) processStreamResponse(c *gin.Context, userID, conversationID uint, chatReq *service.ChatRequest, messageType string) {
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 刷新器
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		return
	}

	// 发送初始连接成功消息
	c.Writer.Write([]byte("data: {\"conversationId\":" + strconv.Itoa(int(conversationID)) + "}\n\n"))
	flusher.Flush()

	var fullContent string
	var fullReasoningContent string
	chunkCount := 0

	// 调用流式接口
	tokens, errors := h.aiService.StreamChat(chatReq)

	// 定义保存消息的闭包
	saveMessage := func() {
		if fullContent == "" {
			return
		}

		msgReq := &service.CreateMessageRequest{
			ConversationID:   conversationID,
			Content:          fullContent,
			ReasoningContent: fullReasoningContent,
			Type:             "assistant",
			Model:            chatReq.Model,
		}

		_, err := h.messageService.Create(userID, msgReq)
		if err != nil {
			log.Printf("保存AI回答失败: %v", err)
		}
	}

Loop:
	for {
		select {
		case <-c.Request.Context().Done():
			log.Println("客户端断开连接，保存已生成内容")
			saveMessage()
			return

		case response, ok := <-tokens:
			if !ok {
				tokens = nil
				if errors == nil {
					break Loop
				}
				continue
			}

			if response.Type == "content" {
				fullContent += response.Content
			}
			chunkCount++

			// 构建SSE数据
			data := map[string]interface{}{
				"type":           messageType,
				"content":        response.Content,
				"done":           false,
				"conversationId": conversationID,
			}

			if response.Type == "reasoning" {
				fullReasoningContent += response.Content
				data["type"] = "reasoning"
			}

			jsonData, _ := json.Marshal(data)
			c.Writer.Write([]byte("data: "))
			c.Writer.Write(jsonData)
			c.Writer.Write([]byte("\n\n"))
			flusher.Flush()

		case err, ok := <-errors:
			if !ok {
				errors = nil
				if tokens == nil {
					break Loop
				}
				continue
			}

			log.Printf("Stream error: %v", err)
			errorData := map[string]interface{}{
				"type":    "error",
				"message": "AI服务调用失败",
				"details": err.Error(),
			}
			jsonErrorData, _ := json.Marshal(errorData)
			c.Writer.Write([]byte("data: "))
			c.Writer.Write(jsonErrorData)
			c.Writer.Write([]byte("\n\n"))
			flusher.Flush()

			saveMessage()
			return
		}
	}

	// 发送完成信号
	finishData := map[string]interface{}{
		"type":           "finish",
		"conversationId": conversationID,
		"content":        fullContent,
		"chunkCount":     chunkCount,
	}
	jsonFinishData, _ := json.Marshal(finishData)
	c.Writer.Write([]byte("data: "))
	c.Writer.Write(jsonFinishData)
	c.Writer.Write([]byte("\n\n"))
	flusher.Flush()

	saveMessage()
}
