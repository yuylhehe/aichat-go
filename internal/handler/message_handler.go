package handler

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
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
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	message, err := h.messageService.Create(userID, &req)
	if err != nil {
		SendInternalError(c, "创建消息失败")
		return
	}

	SendCreated(c, message)
}

// GetList 获取消息列表
func (h *MessageHandler) GetList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	result, err := h.messageService.FindAll(userID)
	if err != nil {
		SendInternalError(c, "获取消息列表失败")
		return
	}

	SendSuccess(c, gin.H{
		"items": result,
	})
}

// GetByConversationID 根据对话ID获取消息
func (h *MessageHandler) GetByConversationID(c *gin.Context) {
	conversationID := c.Param("conversation_id")
	convID, err := strconv.ParseUint(conversationID, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的对话ID")
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.messageService.FindByConversationID(userID, uint(convID))
	if err != nil {
		SendInternalError(c, "获取消息列表失败")
		return
	}

	SendSuccess(c, gin.H{
		"items": result,
	})
}

// GetByID 获取单个消息
func (h *MessageHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的消息ID")
		return
	}

	userID := middleware.GetUserID(c)
	message, err := h.messageService.FindByID(userID, uint(messageID))
	if err != nil {
		SendNotFound(c, "消息不存在")
		return
	}

	SendSuccess(c, message)
}

// Update 更新消息
func (h *MessageHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的消息ID")
		return
	}

	var req dto.UpdateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	message, err := h.messageService.Update(userID, uint(messageID), &req)
	if err != nil {
		SendInternalError(c, "更新消息失败")
		return
	}

	SendSuccess(c, message)
}

// Delete 删除消息
func (h *MessageHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	messageID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的消息ID")
		return
	}

	userID := middleware.GetUserID(c)
	err = h.messageService.Delete(userID, uint(messageID))
	if err != nil {
		SendInternalError(c, "删除消息失败")
		return
	}

	SendSuccess(c, "删除成功")
}
