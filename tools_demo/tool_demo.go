package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"log"
	"os"
)

// --- 1. 定义真实的业务函数 ---
// 实际上这里应该调用 OpenWeatherMap 之类的 API，这里我们 Mock 一下
func GetCurrentWeather(location string, unit string) string {
	fmt.Printf(">>> 正在调用外部 API 查询 %s 的天气 (单位: %s)...\n", location, unit)
	if location == "Beijing" {
		return `{"temperature": "26", "description": "Sunny", "unit": "celsius"}`
	} else if location == "New York" {
		return `{"temperature": "15", "description": "Rainy", "unit": "celsius"}`
	}
	return `{"temperature": "Unknown", "description": "Unknown"}`
}

func main() {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1" // 关键：修改 BaseURL
	client := openai.NewClientWithConfig(config)
	ctx := context.Background()

	// --- 2. 定义工具描述 (Schema) ---
	// 这是告诉 LLM：“我会什么”
	tools := []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_current_weather",
				Description: "获取指定城市的当前天气情况",
				Parameters: jsonschema.Definition{
					Type: jsonschema.Object,
					Properties: map[string]jsonschema.Definition{
						"location": {
							Type:        jsonschema.String,
							Description: "城市名称，例如: Beijing, San Francisco",
						},
						"unit": {
							Type:        jsonschema.String,
							Description: "温度单位，摄氏度或华氏度",
						},
					},
					Required: []string{"location"},
				},
			},
		},
	}

	// --- 3. 第一轮对话：用户提问 ---
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "北京今天出门需要带伞吗？",
		},
	}

	fmt.Println("🤖 第一轮：发送用户问题给 LLM...")
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    "openai/gpt-5", // Qwen 对 Tool Use 支持很好
		Messages: messages,
		Tools:    tools,
	})
	if err != nil {
		log.Fatalln(err)
	}

	msg := resp.Choices[0].Message

	// --- 4. 检查 LLM 是否想调用工具 ---
	if len(msg.ToolCalls) > 0 {
		fmt.Println("⚡ LLM 决定调用工具！")

		// 把 LLM 的回复（包含了它想调用的函数名和参数）加入历史记录
		messages = append(messages, msg)

		// 遍历所有工具调用 (LLM 可能一次想调好几个)
		for _, toolCall := range msg.ToolCalls {
			if toolCall.Function.Name == "get_current_weather" {
				// A. 解析参数
				var args map[string]string
				_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

				// B. 执行本地代码 (Hands)
				location := args["location"]
				unit := args["unit"]
				if unit == "" {
					unit = "celsius"
				}

				result := GetCurrentWeather(location, unit)
				// C. 将结果打包回 Message (Role: Tool)
				// 这一步至关重要：必须带上 ToolCallID，否则 LLM 不知道这是哪个调用的结果
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					Content:    result,
					ToolCallID: toolCall.ID,
				})
			}
		}

		// --- 5. 第二轮对话：把结果发回给 LLM ---
		fmt.Println("🤖 第二轮：把工具运行结果发回给 LLM，等待最终回答...")
		finalResp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    "openai/gpt-5", // Qwen 对 Tool Use 支持很好
			Messages: messages,
			Tools:    tools,
		})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("\n💬 最终回答:")
		fmt.Println(finalResp.Choices[0].Message.Content)
	} else {
		fmt.Println("LLM 没有调用工具，直接回复了：", msg.Content)
	}
}
