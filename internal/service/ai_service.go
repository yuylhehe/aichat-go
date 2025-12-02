package service

import (
	"ai-chat/config"
	"ai-chat/internal/repository"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Message 消息结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages    []Message `json:"messages"`
	Model       *string   `json:"model,omitempty"`
	Temperature *float64  `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Thinking    *struct {
		Type string `json:"type"`
	} `json:"thinking,omitempty"`
}

// StreamResponse 流式响应
type StreamResponse struct {
	Content string
	Type    string // "content" or "reasoning"
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
		Finish  string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// AIService AI服务
type AIService struct {
	db     *gorm.DB
	cfg    *config.Config
	client *http.Client
}

// NewAIService 创建AI服务
func NewAIService(db *gorm.DB, cfg *config.Config) *AIService {
	return &AIService{
		db:  db,
		cfg: cfg,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5分钟超时
			Transport: &http.Transport{
				Proxy:               http.ProxyFromEnvironment, // 加上这行，支持系统代理
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// ChatCompletion 单次聊天完成
func (s *AIService) ChatCompletion(req *ChatRequest) (*ChatResponse, error) {
	if req.Model == nil {
		model := s.cfg.Model
		req.Model = &model
	}
	if req.Temperature == nil {
		temperature := 0.7
		req.Temperature = &temperature
	}

	// 构建请求体
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建请求
	url := s.cfg.BaseURL + "/chat/completions"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.cfg.OpenAIKey)

	// 发送请求
	log.Printf("发送AI请求: URL=%s", url)
	log.Printf("发送AI请求: Body=%s", string(requestBody))

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	log.Printf("收到AI响应: Status=%s", resp.Status)
	log.Printf("收到AI响应: Body=%s", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI API错误: %s", string(body))
	}

	// 解析响应
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	log.Printf("AI响应成功: choices=%d, usage=%d", len(chatResp.Choices), chatResp.Usage.TotalTokens)
	return &chatResp, nil
}

// StreamChat 流式聊天
func (s *AIService) StreamChat(req *ChatRequest) (<-chan StreamResponse, <-chan error) {
	if req.Model == nil {
		model := s.cfg.Model
		req.Model = &model
	}
	if req.Temperature == nil {
		temperature := 0.7
		req.Temperature = &temperature
	}

	// 默认开启思考模式，如果模型支持
	if req.Thinking == nil {
		req.Thinking = &struct {
			Type string `json:"type"`
		}{
			Type: "enabled",
		}
	}

	responses := make(chan StreamResponse, 100)
	errors := make(chan error, 1)

	go func() {
		defer close(responses)
		defer close(errors)

		// 设置流式请求
		streamReq := *req
		streamReq.Stream = true

		requestBody, err := json.Marshal(streamReq)
		if err != nil {
			errors <- fmt.Errorf("序列化流式请求失败: %w", err)
			return
		}

		// 构建请求
		url := s.cfg.BaseURL + "/chat/completions"
		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			errors <- fmt.Errorf("创建流式请求失败: %w", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+s.cfg.OpenAIKey)

		// 发送请求
		log.Printf("发送流式AI请求: model=%s, messages=%d", *streamReq.Model, len(streamReq.Messages))
		log.Printf("发送流式AI请求: 请求体=%s", string(requestBody))

		startTime := time.Now()
		resp, err := s.client.Do(httpReq)
		log.Printf("流式请求响应耗时: %v", time.Since(startTime))

		if err != nil {
			errors <- fmt.Errorf("发送流式请求失败: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errors <- fmt.Errorf("AI流式API错误: %s", string(body))
			return
		}

		// 打印响应头用于调试
		log.Printf("流式响应头: %+v", resp.Header)

		// 处理流式响应
		reader := bufio.NewReader(resp.Body)

		for {
			// 读取一行
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					log.Println("流式响应读取到 EOF")
					break
				}
				log.Printf("读取流式响应错误: %v", err)
				errors <- fmt.Errorf("读取流式响应错误: %w", err)
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if !strings.HasPrefix(line, "data:") {
				continue
			}

			data := strings.TrimPrefix(line, "data:")
			data = strings.TrimSpace(data)
			if data == "[DONE]" {
				log.Println("收到流式响应结束标记 [DONE]")
				return
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content          string `json:"content"`
						ReasoningContent string `json:"reasoning_content"`
					} `json:"delta"`
					Finish interface{} `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				log.Printf("解析流式响应数据错误: %v, data: %s", err, data)
				continue
			}

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta

				// 处理思考内容
				if delta.ReasoningContent != "" {
					responses <- StreamResponse{
						Content: delta.ReasoningContent,
						Type:    "reasoning",
					}
				}

				// 处理普通内容
				if delta.Content != "" {
					responses <- StreamResponse{
						Content: delta.Content,
						Type:    "content",
					}
				}
			}
		}
	}()

	return responses, errors
}

// GetConversationHistory 获取对话历史
func (s *AIService) GetConversationHistory(conversationID uint) ([]Message, error) {
	var messages []repository.Message
	err := s.db.Where("conversation_id = ?", conversationID).
		Order("sort asc").
		Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("获取对话历史失败: %w", err)
	}

	var chatMessages []Message
	for _, msg := range messages {
		role := "assistant"
		if msg.Type == "user" {
			role = "user"
		} else if msg.Type == "system" {
			role = "system"
		}

		chatMessages = append(chatMessages, Message{
			Role:    role,
			Content: msg.Content,
		})
	}

	return chatMessages, nil
}
