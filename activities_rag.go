package HumanintheLoopAI

import (
	"context"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	"go.temporal.io/sdk/activity"
	"time"
)

// 定义 Qdrant 连接配置 (实际项目中建议作为结构体字段注入)
const (
	QdrantHost = "localhost"
	QdrantPort = 6334
	Collection = "news_knowledge_base"
)

func StoreNewsToKnowledgeBase(ctx context.Context, newsList []NewsItem) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("正在存储新闻到知识库", "count", len(newsList))

	// 初始化客户端
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: QdrantHost,
		Port: QdrantPort,
	})
	if err != nil {
		return "", fmt.Errorf("无法连接 Qdrant: %v", err)
	}
	defer client.Close()

	// 2. 确保 Collection 存在 (通常这步在部署脚本里做，但为了演示方便放在这里)
	// 注意：确保 Size 和你 Qwen 返回的维度一致 (1536)
	_ = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: Collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1536,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	// 3. 批量生成向量并准备 Point
	var points []*qdrant.PointStruct
	for i, news := range newsList {
		textToEmbed := fmt.Sprintf("%s\n%s", news.Title, news.Content)

		// 调用你的 Qwen Embedding 函数
		vec, err := GetEmbedding(textToEmbed)
		if err != nil {
			logger.Error("生成向量失败", "title", news.Title, "error", err)
			continue // 跳过这一条，继续下一条
		}

		// 生成唯一 ID (实际项目中可以用 UUID，这里为了简单用毫秒时间戳+序号)
		id := uint64(time.Now().UnixNano()) + uint64(i)
		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(id),
			Vectors: qdrant.NewVectors(vec...),
			Payload: qdrant.NewValueMap(map[string]any{
				"title":   news.Title,
				"content": news.Content,
				"date":    time.Now().Format("2006-01-02"),
			}),
		})
	}

	if len(points) > 0 {
		wait := true
		_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
			CollectionName: Collection,
			Points:         points,
			Wait:           &wait,
		})
		if err != nil {
			return "", fmt.Errorf("插入向量失败: %v", err)
		}
	}

	return fmt.Sprintf("成功存入 %d 条新闻到知识库", len(points)), nil
}
