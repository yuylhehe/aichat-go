package handler

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	messageService *service.MessageService
}

// NewMessageHandler 创建消息处理器
func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// Create 创建消息
func (h *MessageHandler) Create(c *gin.Context) {
	var req dto.CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	messageReq := &service.CreateMessageRequest{
		ConversationID: req.ConversationID,
		Content:        req.Content,
		Type:           req.Type,
	}

	userID := middleware.GetUserID(c)
	message, err := h.messageService.Create(userID, messageReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "创建消息失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": h.toResponse(message),
	})
}

// GetList 获取消息列表
func (h *MessageHandler) GetList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	result, err := h.messageService.FindAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取消息列表失败",
			"details": err.Error(),
		})
		return
	}

	items := make([]*dto.MessageResponse, len(result))
	for i, item := range result {
		items[i] = h.toResponse(item)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items": items,
		},
	})
}

// GetByConversationID 根据对话ID获取消息
func (h *MessageHandler) GetByConversationID(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	convID, err := strconv.ParseUint(conversationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的对话ID",
		})
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.messageService.FindByConversationID(userID, uint(convID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取消息列表失败",
			"details": err.Error(),
		})
		return
	}

	items := make([]*dto.MessageResponse, len(result))
	for i, item := range result {
		items[i] = h.toResponse(item)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items": items,
		},
	})
}

// GetByID 获取单个消息
func (h *MessageHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的消息ID",
		})
		return
	}

	userID := middleware.GetUserID(c)
	message, err := h.messageService.FindByID(userID, uint(messageID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "消息不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": h.toResponse(message),
	})
}

// Update 更新消息
func (h *MessageHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的消息ID",
		})
		return
	}

	var req dto.UpdateMessageRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	updateReq := &service.UpdateMessageRequest{
		Content: req.Content,
	}
	userID := middleware.GetUserID(c)
	message, err := h.messageService.Update(userID, uint(messageID), updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "更新消息失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": h.toResponse(message),
	})
}

// Delete 删除消息
func (h *MessageHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的消息ID",
		})
		return
	}

	userID := middleware.GetUserID(c)
	err = h.messageService.Delete(userID, uint(messageID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "删除消息失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

func (h *MessageHandler) toResponse(item *service.MessageResponse) *dto.MessageResponse {
	return &dto.MessageResponse{
		ID:               item.ID,
		ConversationID:   item.ConversationID,
		Content:          item.Content,
		ReasoningContent: item.ReasoningContent,
		Sort:             item.Sort,
		Type:             item.Type,
		Tokens:           item.Tokens,
		Model:            item.Model,
		ParentID:         item.ParentID,
		Metadata:         item.Metadata,
		CreatedAt:        item.CreatedAt.Format(common.TimeLayout),
	}
}
