package main

import (
	"context"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	HumanintheLoopAI "humanintheloop"
	"log"
	"strings"
)

func main() {
	userQuestion := "Temporal 是什么？它有什么用？"
	fmt.Printf("🙋 用户提问: %s\n", userQuestion)

	// 1. 连接 Qdrant
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		log.Fatalf("❌ 无法连接 Qdrant: %v\n\n", err)
	}
	defer client.Close()

	// 2. 将用户问题转化为向量 (使用 llm_helper.go 中的 GetEmbedding)
	queryVec, err := HumanintheLoopAI.GetEmbedding(userQuestion)
	if err != nil {
		log.Fatalf("❌ 生成向量失败: %v\n\n", err)
	}

	// 3. 在 Qdrant 中搜索相关上下文 (Context)
	fmt.Println("🔍 正在知识库中检索...")
	searchResult, err := client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: HumanintheLoopAI.Collection,
		Query:          qdrant.NewQuery(queryVec...),
		Limit:          &[]uint64{3}[0],
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		log.Fatal("搜索失败:", err)
	}

	// 4. 组装 Prompt (Prompt Engineering)
	var contextBuilder strings.Builder
	for _, point := range searchResult {
		// 只有相似度足够高才用 (阈值过滤)
		if point.Score > 0.5 {
			content := point.Payload["content"].String()
			contextBuilder.WriteString(content)
			contextBuilder.WriteString("\n---\n")
		}
	}
	contextText := contextBuilder.String()
	if contextText == "" {
		contextText = "抱歉，知识库中没有找到相关信息。"
	} else {
		fmt.Printf("📖 找到参考资料 (Top match score: %.4f)\n", searchResult[0].Score)
	}

	// 构造最终给 LLM 的指令
	systemPrompt := "你是一个基于知识库的 AI 助手。请根据下方的【参考资料】回答用户问题。如果参考资料里没有答案，请诚实地说不知道。"
	finalPrompt := fmt.Sprintf("【参考资料】:\n%s\n\n【用户问题】: %s", contextText, userQuestion)

	// 5. 生成回答 (Generation)
	fmt.Println("🤖 正在思考...")
	answer, err := HumanintheLoopAI.ChatWithLLM(systemPrompt, finalPrompt)
	if err != nil {
		log.Fatal("生成回答失败:", err)
	}

	fmt.Println("\n💬 AI 回答:")
	fmt.Println("--------------------------------------------------")
	fmt.Println(answer)
	fmt.Println("--------------------------------------------------")
}
