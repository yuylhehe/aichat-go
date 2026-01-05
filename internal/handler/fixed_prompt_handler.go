package handler

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FixedPromptHandler 固定提示词处理器
type FixedPromptHandler struct {
	fixedPromptService *service.FixedPromptService
}

// NewFixedPromptHandler 创建固定提示词处理器
func NewFixedPromptHandler(fixedPromptService *service.FixedPromptService) *FixedPromptHandler {
	return &FixedPromptHandler{
		fixedPromptService: fixedPromptService,
	}
}

// Create 创建固定提示词
func (h *FixedPromptHandler) Create(c *gin.Context) {
	var req dto.CreateFixedPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	response, err := h.fixedPromptService.Create(userID, &req)
	if err != nil {
		SendInternalError(c, "创建固定提示词失败: "+err.Error())
		return
	}

	SendCreated(c, response)
}

// GetList 获取固定提示词列表
func (h *FixedPromptHandler) GetList(c *gin.Context) {
	var req dto.GetFixedPromptsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	result, err := h.fixedPromptService.FindAll(userID, req.Page, req.PageSize, req.Search)
	if err != nil {
		SendInternalError(c, "获取固定提示词列表失败: "+err.Error())
		return
	}

	SendSuccess(c, result)
}

// GetByID 根据ID获取固定提示词
func (h *FixedPromptHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		SendBadRequest(c, "无效的ID格式")
		return
	}

	userID := middleware.GetUserID(c)
	response, err := h.fixedPromptService.FindByID(userID, uint(id))
	if err != nil {
		SendNotFound(c, "固定提示词不存在或无权访问")
		return
	}

	SendSuccess(c, response)
}

// Update 更新固定提示词
func (h *FixedPromptHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		SendBadRequest(c, "无效的ID格式")
		return
	}

	var req dto.UpdateFixedPromptRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	response, err := h.fixedPromptService.Update(userID, uint(id), &req)
	if err != nil {
		SendInternalError(c, "更新固定提示词失败或无权更新: "+err.Error())
		return
	}

	SendSuccess(c, response)
}

// Delete 删除固定提示词
func (h *FixedPromptHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		SendBadRequest(c, "无效的ID格式")
		return
	}

	userID := middleware.GetUserID(c)
	err = h.fixedPromptService.Delete(userID, uint(id))
	if err != nil {
		SendInternalError(c, "删除固定提示词失败或无权删除: "+err.Error())
		return
	}

	SendSuccess(c, "删除成功")
}
