## MCP

### 启动方式

```shell
go run client/client.go
```

### MCP的特性说明

稀奇之处在于“解耦”和“标准化转换”：

动态性：如果你在 server/main.go 里加了一个 subtract 工具，Client 的代码一行都不用改。Client 再次运行时，ListTools 就会自动发现新工具。

对接 LLM 的桥梁： 在真实的 Agent 代码中，我们会把 cli.ListTools() 返回的结果，自动转换为 OpenAI 格式的 tools: [...] JSON。

核心代码片段（伪代码）：

```go
// 自动把 MCP 工具转为 OpenAI 工具
openAITools := []openai.Tool{}
mcpTools, _ := cli.ListTools(...)

for _, t := range mcpTools.Tools {
    openAITools = append(openAITools, openai.Tool{
        Type: "function",
        Function: &openai.FunctionDefinition{
            Name: t.Name,
            Description: t.Description,
            Parameters: t.InputSchema, // MCP 的 Schema 和 OpenAI 是完全兼容的！
        },
    })
}

// 然后传给 LLM
client.CreateChatCompletion(..., Tools: openAITools)
```

LLM 如何调用 tools 参考如下代码：

[tools调用方式](../tools_demo/tool_demo.go)

### 总结

MCP 的核心价值在于：Write Once, Use Everywhere.

不管是连接本地的 SQLite，还是连接 GitHub API，还是连接 Google Drive。

只要有人把那个服务封装成了 MCP Server。

你的 Agent 只需要实现一次 MCP Client 逻辑，就可以动态加载这成百上千种工具，而不需要每次都去写适配器代码。