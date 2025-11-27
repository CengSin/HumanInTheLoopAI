package HumanintheLoopAI

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/activity"
	"time"
)

type NewsItem struct {
	Title   string
	Content string
}

// Activity 1: 抓取新闻
// 模拟网络请求，去获取特定主题的新闻
func FetchNews(ctx context.Context, topic string) ([]NewsItem, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("正在抓取新闻", "topic", topic)

	// 模拟网络延迟
	time.Sleep(2 * time.Second)

	// 返回一些假数据
	return []NewsItem{
		{Title: "Temporal Go SDK 发布", Content: "Go 语言是构建高并发后端服务的首选..."},
		{Title: "AI Agent 架构演进", Content: "从简单的 Prompt 工程走向复杂的编排..."},
	}, nil
}

// Activity 2: AI 总结
// 模拟调用 OpenAI 或其他 LLM 接口
func AISummarize(ctx context.Context, news []NewsItem) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("正在调用 LLM 进行总结...", "count", len(news))

	// 模拟思考时间
	time.Sleep(2 * time.Second)

	return fmt.Sprintf("【今日简报】\n共获取 %d 条新闻。\n重点：Temporal 结合 Go 语言能构建极其稳定的 Agent 后端。", len(news)), nil
}

// Activity 3: 发送邮件
// 模拟发送操作
func SendEmail(ctx context.Context, content string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("正在发送邮件...", "preview", content[:10])

	time.Sleep(1 * time.Second)
	return "邮件发送成功", nil
}
