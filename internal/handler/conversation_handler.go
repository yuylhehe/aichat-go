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
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	userID := middleware.GetUserID(c)

	conversationReq := &service.CreateConversationRequest{
		Name:   req.Name,
		UserID: userID,
	}

	conversation, err := h.conversationService.Create(conversationReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "创建对话失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": dto.ConversationResponse{
			ID:           conversation.ID,
			Name:         conversation.Name,
			UserID:       conversation.UserID,
			IsActive:     conversation.IsActive,
			SystemPrompt: conversation.SystemPrompt,
			Model:        conversation.Model,
			Temperature:  conversation.Temperature,
			CreatedAt:    conversation.CreatedAt.Format(common.TimeLayout),
			UpdatedAt:    conversation.UpdatedAt.Format(common.TimeLayout),
			Messages:     conversation.Messages,
		},
	})
}

// GetList 获取对话列表
func (h *ConversationHandler) GetList(c *gin.Context) {
	userID := middleware.GetUserID(c)
	search := c.Query("search")

	result, err := h.conversationService.FindByUserID(userID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取对话列表失败",
			"details": err.Error(),
		})
		return
	}

	// 转换服务层的响应为处理器响应
	items := make([]*dto.ConversationResponse, len(result))
	for i, item := range result {
		items[i] = &dto.ConversationResponse{
			ID:           item.ID,
			Name:         item.Name,
			UserID:       item.UserID,
			IsActive:     item.IsActive,
			SystemPrompt: item.SystemPrompt,
			Model:        item.Model,
			Temperature:  item.Temperature,
			CreatedAt:    item.CreatedAt.Format(common.TimeLayout),
			UpdatedAt:    item.UpdatedAt.Format(common.TimeLayout),
			Messages:     item.Messages,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"items": items,
		},
	})
}

// GetByID 获取单个对话
func (h *ConversationHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	conversationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的对话ID",
		})
		return
	}

	userID := middleware.GetUserID(c)

	conversation, err := h.conversationService.FindByID(userID, uint(conversationID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "对话不存在",
		})
		return
	}

	// 获取消息列表
	messages, err := h.messageService.FindByConversationID(userID, uint(conversationID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "获取消息列表失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"conversation": dto.ConversationResponse{
				ID:           conversation.ID,
				Name:         conversation.Name,
				UserID:       conversation.UserID,
				IsActive:     conversation.IsActive,
				SystemPrompt: conversation.SystemPrompt,
				Model:        conversation.Model,
				Temperature:  conversation.Temperature,
				CreatedAt:    conversation.CreatedAt.Format(common.TimeLayout),
				UpdatedAt:    conversation.UpdatedAt.Format(common.TimeLayout),
				Messages:     conversation.Messages,
			},
			"messages": messages,
		},
	})
}

// Update 更新对话
func (h *ConversationHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	conversationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的对话ID",
		})
		return
	}

	var req dto.UpdateConversationRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	userID := middleware.GetUserID(c)

	// 更新对话
	updateReq := &service.UpdateConversationRequest{
		Name: req.Name,
	}
	conversation, err := h.conversationService.Update(userID, uint(conversationID), updateReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "更新对话失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": dto.ConversationResponse{
			ID:           conversation.ID,
			Name:         conversation.Name,
			UserID:       conversation.UserID,
			IsActive:     conversation.IsActive,
			SystemPrompt: conversation.SystemPrompt,
			Model:        conversation.Model,
			Temperature:  conversation.Temperature,
			CreatedAt:    conversation.CreatedAt.Format(common.TimeLayout),
			UpdatedAt:    conversation.UpdatedAt.Format(common.TimeLayout),
			Messages:     conversation.Messages,
		},
	})
}

// Delete 删除对话
func (h *ConversationHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	conversationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的对话ID",
		})
		return
	}

	userID := middleware.GetUserID(c)

	// 删除对话
	err = h.conversationService.Delete(userID, uint(conversationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "删除对话失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}
