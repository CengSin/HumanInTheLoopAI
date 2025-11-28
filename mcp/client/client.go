package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
)

func main() {
	ctx := context.Background()

	// 1. 定义要连接的 Server 命令
	// 我们告诉 Client："去运行 server 文件夹下的 main.go"
	// 在实际生产中，这可能是一个 Docker 容器或者远程 SSE 连接
	// 2. 创建 MCP Client (基于 Stdio)
	cli, err := client.NewStdioMCPClient("go", []string{}, "run", "../server/main.go")
	if err != nil {
		log.Fatalf("无法启动 Server: %v", err)
	}
	defer cli.Close()

	// 🔴【关键修复】启动后台监听服务
	// 这会启动一个 goroutine 来读取 Server 的 Stdout，否则就会死锁
	if err = cli.Start(ctx); err != nil {
		log.Fatalf("无法启动客户端监听: %v", err)
	}

	// 3. 握手 (Initialize)
	// 就像 USB 插入时的握手，交换能力信息
	initializeRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			Capabilities:    mcp.ClientCapabilities{},
			ClientInfo: mcp.Implementation{
				Name:    "my-math-server",
				Version: "1.0.0",
			},
		},
	}

	result, err := cli.Initialize(ctx, initializeRequest)
	if err != nil {
		log.Fatalf("握手失败: %v", err)
	}
	fmt.Printf("✅ 连接成功！Server 名称: %s, 版本: %s\n", result.ServerInfo.Name, result.ServerInfo.Version)

	// 4. 工具发现 (List Tools) - 关键步骤！
	// Agent 问："你都有什么本事？"
	toolsList, err := cli.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("获取工具列表失败: %v", err)
	}

	fmt.Printf("🔍 发现 %d 个工具:\n", len(toolsList.Tools))
	for _, t := range toolsList.Tools {
		fmt.Printf("   - 工具名: %s, 描述: %s\n", t.Name, t.Description)
	}

	// 5. 调用工具 (Call Tool)
	// 假设 LLM 决定调用 "add"
	callResult, err := cli.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "add",
			Arguments: map[string]string{
				"a": "10",
				"b": "25",
			},
		},
	})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
	}

	if len(callResult.Content) > 0 {
		text := callResult.Content[0].(mcp.TextContent).Text
		fmt.Printf("🎉 Server 返回结果: %+v\n", text)
	}
}
