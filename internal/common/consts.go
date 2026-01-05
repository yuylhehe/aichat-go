package common

import "time"

const TimeLayout = "2006-01-02 15:04:05"

// 认证配置
const (
	AccessTokenExpire  = 24 * time.Hour     // 24小时过期
	RefreshTokenExpire = 7 * 24 * time.Hour // 7天过期
)

// 消息角色
const (
	RoleSystem    = "system"
	RoleUser      = "user"
	RoleAssistant = "assistant"
)

// SSE 消息类型
const (
	SSETypeMessage   = "message"
	SSETypeToken     = "token"
	SSETypeReasoning = "reasoning"
	SSETypeContent   = "content"
	SSETypeFinish    = "finish"
	SSETypeError     = "error"
)

// 模型列表
var Models = []string{
	"gpt-3.5-turbo",
	"gpt-3.5-turbo-16k",
	"gpt-4",
	"gpt-4-turbo",
	"gpt-4-turbo-preview",
}

// AI 默认配置
const (
	DefaultTemperature = 0.7
	ThinkingEnabled    = "enabled"
)
