package HumanintheLoopAI

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

// GetEmbedding 调用 OpenRouter 获取文本向量
func GetEmbedding(text string) ([]float32, error) {
	// 1. 从环境变量获取 Key (安全起见，不要硬编码)
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("未设置 OPENROUTER_API_KEY 环境变量")
	}

	// 2. 配置 Client 指向 OpenRouter
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1" // 关键：修改 BaseURL

	client := openai.NewClientWithConfig(config)

	// 3. 发起请求
	// OpenRouter 上可以用 "qwen/qwen3-embedding-8b"
	res, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input:          text,
		Model:          "qwen/qwen3-embedding-8b",
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
		Dimensions:     1536,
	})
	if err != nil {
		log.Printf("Embedding API 调用失败: %v", err)
		return nil, err
	}
	if len(res.Data) == 0 {
		return nil, fmt.Errorf("返回数据为空")
	}

	// 4. 返回向量 (OpenAI Small 模型通常是 1536 维)
	vec := res.Data[0].Embedding
	fmt.Printf("向量维度检查: %d\n", len(vec)) // <--- 必须确认是 1536

	if len(vec) != 1536 {
		// 如果这里打印出 4096，你需要修改 Qdrant 创建 Collection 时的 Size 为 4096
		log.Fatalf("维度不匹配！期望 1536，实际返回 %d", len(vec))
	}
	return vec, nil
}

// ChatWithLLM 发送提示词给 LLM 并获取回复
func ChatWithLLM(systemPrompt, userPrompt string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1" // 关键：修改 BaseURL

	client := openai.NewClientWithConfig(config)

	res, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "moonshotai/kimi-k2-0905",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: userPrompt},
		},
	})
	if err != nil {
		return "", err
	}

	return res.Choices[0].Message.Content, nil
}
