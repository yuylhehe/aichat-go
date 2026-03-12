package handler

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ConversationHandler 对话处理器
type ConversationHandler struct {
	conversationService *service.ConversationService
	messageService      *service.MessageService
}

// NewConversationHandler 创建对话处理器
func NewConversationHandler(
	conversationService *service.ConversationService,
	messageService *service.MessageService,
) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
		messageService:      messageService,
	}
}

// Create 创建对话
func (h *ConversationHandler) Create(c *gin.Context) {
	var req dto.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	conversation, err := h.conversationService.Create(userID, &req)
	if err != nil {
		SendInternalError(c, "创建对话失败")
		return
	}

	SendCreated(c, conversation)
}

// GetList 获取对话列表
func (h *ConversationHandler) GetList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	search := c.Query("search")

	result, err := h.conversationService.FindByUserID(userID, search)
	if err != nil {
		SendInternalError(c, "获取对话列表失败")
		return
	}

	SendSuccess(c, gin.H{
		"items": result,
	})
}

// GetByID 获取单个对话
func (h *ConversationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	conversationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的对话ID")
		return
	}

	userID := middleware.GetUserID(c)

	conversation, err := h.conversationService.FindByID(userID, uint(conversationID))
	if err != nil {
		SendNotFound(c, "对话不存在")
		return
	}

	// 获取消息列表
	messages, err := h.messageService.FindByConversationID(userID, uint(conversationID))

	if err != nil {
		SendInternalError(c, "获取消息列表失败")
		return
	}

	SendSuccess(c, gin.H{
		"conversation": conversation,
		"messages":     messages,
	})
}

// Update 更新对话
func (h *ConversationHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	conversationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的对话ID")
		return
	}

	var req dto.UpdateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	conversation, err := h.conversationService.Update(userID, uint(conversationID), &req)
	if err != nil {
		SendInternalError(c, "更新对话失败")
		return
	}

	SendSuccess(c, conversation)
}

// Delete 删除对话
func (h *ConversationHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	conversationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		SendBadRequest(c, "无效的对话ID")
		return
	}

	userID := middleware.GetUserID(c)

	// 删除对话
	err = h.conversationService.Delete(userID, uint(conversationID))
	if err != nil {
		SendInternalError(c, "删除对话失败")
		return
	}

	SendSuccess(c, "删除成功")
}
