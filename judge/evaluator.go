package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
)

func CallJudgeLLM(prompt string) (string, error) {
	ctx := context.Background()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://openrouter.ai/api/v1"
	client := openai.NewClientWithConfig(config)

	// 考官通常需要用最聪明的模型，比如 GPT-4o 或 Claude-3.5-Sonnet
	// 这里演示用 gpt5
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "openai/gpt-5",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "你是一个公正的 AI 评测员。请根据提供的标准进行打分，只输出 0 到 10 之间的数字。"},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: 0, // 评测需要严谨，不需要创意
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	os.Setenv("HTTP_PROXY", "127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "127.0.0.1:7890")

	// 假设这是我们 RAG 运行的一次现场记录
	userQuestion := "Temporal 的核心优势是什么？"

	// 检索到的 Context (从 Qdrant 拿到的)
	retrievedContext := `
	Temporal 是一个持久化执行平台。
	它最大的优势是 Fault Tolerance (容错性) 和 State Management (状态管理)。
	即使进程崩溃，Workflow 也能从断点恢复。
	`
	// Agent 生成的 Answer
	agentAnswer := "Temporal 的核心优势是它的容错性和状态管理能力，可以保证代码在崩溃后恢复运行。"

	// ----------------------------------------------------
	// 指标 1: 忠实度 (Faithfulness)
	// 检查 Answer 是否包含 Context 中不存在的信息（幻觉）
	// ----------------------------------------------------
	faithfulnessPrompt := fmt.Sprintf(`
	请评估【生成的回答】是否完全基于【参考资料】？
	如果是，打 10 分；如果有编造信息，打 0 分。
	
	【参考资料】: %s
	【生成的回答】: %s
	
	请只输出分数：
	`, retrievedContext, agentAnswer)

	score1, _ := CallJudgeLLM(faithfulnessPrompt)
	fmt.Printf("📊 忠实度评分 (Faithfulness): %s/10\n", strings.TrimSpace(score1))

	// ----------------------------------------------------
	// 指标 2: 相关度 (Relevance)
	// 检查 Answer 是否回答了 User Question
	// ----------------------------------------------------
	relevancePrompt := fmt.Sprintf(`
	请评估【生成的回答】是否直接且准确地回答了【用户问题】？
	如果是，打 10 分；答非所问，打 0 分。
	
	【用户问题】: %s
	【生成的回答】: %s
	
	请只输出分数：
	`, userQuestion, agentAnswer)
	score2, _ := CallJudgeLLM(relevancePrompt)
	fmt.Printf("🎯 相关度评分 (Relevance):   %s/10\n", strings.TrimSpace(score2))
}
