package main

import (
	"context"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	HumanintheLoopAI "humanintheloop"
	"log"
)

func RealDataQueryWithAI() {
	c, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		log.Fatalln("无法创建 Qdrant Client", err)
	}
	defer c.Close()

	ctx := context.Background()
	collectionName := "news_knowledge_base"

	// --- 2. 创建集合 (注意维度改为 1536) ---
	// 每次运行前先删除旧的，方便测试
	c.DeleteCollection(ctx, collectionName)
	err = c.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1536,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		log.Fatalln("无法创建集合", err)
	}
	fmt.Println("✅ 集合创建成功！")

	// --- 3. 准备真实数据 ---
	documents := []string{
		"Temporal 是一种工作流编排平台，保证代码的持久化运行。", // 知识点 A
		"Go 语言的高并发特性非常适合构建后台服务。",            // 知识点 B
		"今天的晚餐吃的是红烧牛肉面，味道很不错。",              // 干扰项
	}

	var points []*qdrant.PointStruct
	for i, doc := range documents {
		fmt.Printf("正在生成向量 (%d/%d): %s...\n", i+1, len(documents), doc[:10])
		// 调用 OpenRouter 生成向量
		vec, err := HumanintheLoopAI.GetEmbedding(doc)
		if err != nil {
			log.Fatalf("生成向量失败: %v", err)
		}

		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(uint64(i + 1)),
			Vectors: qdrant.NewVectors(vec...),
			Payload: qdrant.NewValueMap(map[string]any{"content": doc}),
		})
	}

	// 插入数据
	wait := true
	_, err = c.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
		Wait:           &wait,
	})
	if err != nil {
		log.Fatalln("无法插入数据", err)
	}
	fmt.Println("✅ 数据插入成功！")

	// --- 4. 真实语义搜索 ---
	// 用户提问 (注意：这句话里没有 "Temporal" 也没有 "持久化" 这些词，只有 "可靠性")
	queryText := "如何保证后端任务的可靠性？"
	fmt.Printf("\n🔍 用户提问: %s\n", queryText)

	query, err := HumanintheLoopAI.GetEmbedding(queryText)
	if err != nil {
		log.Fatalf("生成查询向量失败: %v", err)
	}
	searchResult, err := c.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(query...),
		Limit:          &[]uint64{3}[0],
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}

	if len(searchResult) > 0 {
		bestMatch := searchResult[0]
		fmt.Printf("🏆 最佳匹配: %s\n", bestMatch.Payload["content"])
		fmt.Printf("   相似度: %.4f\n", bestMatch.Score)
	} else {
		fmt.Println("未找到匹配结果")
	}
}
