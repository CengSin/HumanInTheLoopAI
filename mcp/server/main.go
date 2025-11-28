package main

import (
	"context"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 1. 创建 MCP Server 实例
	// 参数是 Server 的名字和版本
	s := server.NewMCPServer("my-math-server", "1.0.0")

	// 2. 注册一个工具 (Tool)
	// 这就像我们之前给 LLM 写 JSON Schema，但这里是用 Go 代码写的
	tool := mcp.NewTool("add",
		mcp.WithDescription("计算两个数字的和"),
		mcp.WithString("a", mcp.Required(), mcp.Description("第一个数字")),
		mcp.WithString("b", mcp.Required(), mcp.Description("第二个数字")),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, ok := request.GetArguments()["a"].(string)
		if !ok {
			return mcp.NewToolResultError("参数 a 必须是字符串"), nil
		}

		b, ok := request.GetArguments()["b"].(string)
		if !ok {
			return mcp.NewToolResultError("参数 a 必须是字符串"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("结果是: %s + %s", a, b)), nil
	})

	// 4. 启动服务 (通过 Stdio 通信)
	// 这样 Client 就可以通过命令行启动这个进程并与之对话
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
