package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	HumanintheLoopAI "humanintheloop"
	"log"
)

func main() {
	// 1. 创建 Temporal Client
	// 默认连接到本地 localhost:7233
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("无法创建 Temporal Client", err)
	}
	defer c.Close()

	// 2. 启动 Worker
	// Worker 负责监听 Task Queue，并执行具体的 Workflow 和 Activity 代码
	we := worker.New(c, HumanintheLoopAI.TaskQueueName, worker.Options{})
	we.RegisterWorkflow(HumanintheLoopAI.NewsAgentWorkflow)
	we.RegisterWorkflow(HumanintheLoopAI.ManagerAgentWorkflow)
	we.RegisterWorkflow(HumanintheLoopAI.ResearcherAgentWorkflow)
	we.RegisterWorkflow(HumanintheLoopAI.WriterAgentWorkflow)
	we.RegisterActivity(HumanintheLoopAI.FetchNews)
	we.RegisterActivity(HumanintheLoopAI.AISummarize)
	we.RegisterActivity(HumanintheLoopAI.SendEmail)
	we.RegisterActivity(HumanintheLoopAI.StoreNewsToKnowledgeBase)

	// 3. 运行worker
	if err := we.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("无法启动 Worker", err)
	}
}
