package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	HumanintheLoopAI "humanintheloop"
	"log"
)

func main() {
	// 1. 创建一个Temporal client
	// 默认连接到本地的localhost:7233
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("无法创建 Temporal Client", err)
	}
	defer c.Close()

	workflowID := "human-in-the-loop-ai-workflow" + uuid.New().String()
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: HumanintheLoopAI.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, HumanintheLoopAI.NewsAgentWorkflow, "人工智能")
	if err != nil {
		log.Fatalln("无法执行 Workflow", err)
	}

	var result string
	if err = we.Get(context.Background(), &result); err != nil {
		log.Fatalln("无法获取 Workflow 结果", err)
	}
	fmt.Println("Workflow 结果:", result)
}
