package HumanintheLoopAI

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

type ReviewSignal struct {
	Action   string // "APPROVE" or "REJECT"
	Feedback string // 如果拒绝，这里放修改建议
}

// 定义用户提问的信号结构
type QuestionSignal struct {
	Question      string
	ReplyChanName string // 这是一个演示概念，实际通过 Polling 或其他方式获取结果
}

func NewsAgentWorkflow(ctx workflow.Context, topic string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// ==========================================
	// 步骤 1: 抓取新闻
	// ==========================================
	var newsList []NewsItem
	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow: 抓取新闻中...", "topic", topic)
	// ExecuteActivity 是阻塞的，直到 Activity 完成（或重试耗尽）
	// Get(ctx, &newsList) 会把结果解析到 newsList 变量中
	if err := workflow.ExecuteActivity(ctx, FetchNews, topic).Get(ctx, &newsList); err != nil {
		return "", err
	}

	// ==========================================
	// 🆕 新增步骤：存入记忆 (RAG)
	// ==========================================
	var storeResult string
	if err := workflow.ExecuteActivity(ctx, StoreNewsToKnowledgeBase, newsList).Get(ctx, &storeResult); err != nil {
		// 即使存储失败，也许不应该阻断整个流程？这里演示我们选择“报错返回”
		logger.Error("知识库存储失败", "error", err)
		return "", err
	}
	logger.Info("知识库更新完毕", "result", storeResult)

	// ==========================================
	// 步骤 2: AI 总结
	// ==========================================
	var summary string
	if err := workflow.ExecuteActivity(ctx, AISummarize, newsList).Get(ctx, &summary); err != nil {
		return "", err
	}

	// 创建 Selector 来处理多个信号 (Review 和 Question)
	selector := workflow.NewSelector(ctx)
	var approvalResult string

	// ==========================================
	// 重点：进入审核阻塞
	// ==========================================
	// 获取信号通道，名字叫 "ReviewSignal"
	reviewChan := workflow.GetSignalChannel(ctx, "ReviewSignal")
	selector.AddReceive(reviewChan, func(c workflow.ReceiveChannel, more bool) {
		var signal ReviewSignal
		c.Receive(ctx, &signal)
		approvalResult = signal.Action
		if signal.Action == "APPROVE" {
			logger.Info("内容已通过审核")
		} else if signal.Action == "REJECT" {
			logger.Info("审核拒绝，要求重写", "feedback", signal.Feedback)
			if err := workflow.ExecuteActivity(ctx, AISummarize, newsList).Get(ctx, &summary); err != nil {
				logger.Error("重写总结失败", "error", err)
			}
		}
	})

	for approvalResult != "APPROVE" {
		// 阻塞等待另一个信号的到来
		selector.Select(ctx)
	}

	// ==========================================
	// 步骤 3: 发送邮件 (暂时直接发送)
	// ==========================================
	var emailStatus string
	if err := workflow.ExecuteActivity(ctx, SendEmail, summary).Get(ctx, &emailStatus); err != nil {
		return "", err
	}

	logger.Info("Workflow 完成", "emailStatus", emailStatus)
	return summary, nil
}
