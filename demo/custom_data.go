package main

import (
	"context"
	"fmt"
	"github.com/qdrant/go-client/qdrant"
	"log"
)

func CustomDataQuery() {
	c, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		log.Fatalln("无法创建 Qdrant Client", err)
	}
	defer c.Close()

	ctx := context.Background()
	collectionName := "vector_test"

	// 2. 创建集合 (Create Collection)
	// 我们定义向量的维度为 2 (方便我们在脑海中想象 X, Y 坐标)
	// DistanceCosine 表示使用余弦相似度计算距离
	err = c.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     2,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		log.Println("无法创建集合", err)
	} else {
		fmt.Println("✅ 集合创建成功！")
	}

	// 3. 插入数据 (Upsert)
	// 我们手动制造 3 个向量
	points := []*qdrant.PointStruct{
		{
			Id: qdrant.NewIDNum(1),
			// 向量 [0.9, 0.1] -> 代表 "北京" (偏向 X 轴)
			Vectors: qdrant.NewVectors(0.9, 0.1),
			Payload: qdrant.NewValueMap(map[string]any{"city": "北京", "type": "城市"}),
		},
		{
			Id: qdrant.NewIDNum(2),
			// 向量 [0.8, 0.2] -> 代表 "上海" (也很偏向 X 轴，离北京很近)
			Vectors: qdrant.NewVectors(0.8, 0.2),
			Payload: qdrant.NewValueMap(map[string]any{"city": "上海", "type": "城市"}),
		},
		{
			Id: qdrant.NewIDNum(3),
			// 向量 [0.1, 0.9] -> 代表 "苹果" (偏向 Y 轴，离城市很远)
			Vectors: qdrant.NewVectors(0.1, 0.9),
			Payload: qdrant.NewValueMap(map[string]any{"fruit": "苹果", "type": "水果"}),
		},
	}

	wait := true // 等待数据落盘
	_, err = c.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
		Wait:           &wait,
	})
	if err != nil {
		log.Fatalln("无法插入数据", err)
	}
	fmt.Println("✅ 数据插入成功！")

	// 4. 语义搜索 (Search)
	// 假设我们现在搜 "天津"
	// 虽然数据库里没有 "天津"，但我们知道 "天津" 的向量应该离 X 轴很近
	// 我们构造一个查询向量: [0.95, 0.05]
	queryVector := []float32{0.95, 0.05}
	searchResult, err := c.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(queryVector...),
		Limit:          &[]uint64{3}[0],             // 返回前3条结果
		WithPayload:    qdrant.NewWithPayload(true), // 把元数据也拿回来
	})
	if err != nil {
		log.Fatalf("搜索失败: %v", err)
	}

	for i, point := range searchResult {
		payload := point.Payload
		// 这里的 Score 就是相似度 (越接近 1 越相似)
		fmt.Printf("%d. 找到: %v (ID: %v), 相似度 Score: %.4f\n",
			i+1, payload["city"], point.Id, point.Score)
	}
}
