# Temporal

一个简易的使用Qdrant数据库配合Temporal工作流引擎的RAG项目

## 分支说明

- main 分支：包含完整的RAG功能实现，使用Qdrant进行向量存储和检索，Temporal进行数据处理和保存。
- [tool_use](https://github.com/CengSin/HumanInTheLoopAI/tree/tool_use) 分支：包含了llms使用tools以及mcp的示例代码

## 目录结构

[client](./client)是temporal的client示例

[worker](./worker)文件夹下是temporal的worker示例

[workflow](./workflow.go) 中写了workflow的详细逻辑，并且调用了activity

[activity](./activity.go) 中写了activity的详细逻辑

[activity_rag](./activity_rag.go) 中写了RAG相关的activity逻辑

[llm_helper](./llm_helper.go) 中封装了调用大模型的逻辑, 使用openrouter调用embedding和chat模型，进行向量数据的生成和RAG检索对话。

[rag_search](./rag_search) 中写了一个简单的RAG检索逻辑，即通过把用户的问题生成embedding，然后在Qdrant中进行向量检索，最后把检索到的上下文信息给LLM，LLM根据上下文信息生成回答并返回。

## 运行步骤

1. 安装并运行Qdrant数据库，可以参考[Docker运行Qdrant](./数据库.md)。
2. 安装并运行Temporal服务器，可以参考[Temporal安装文档](https://docs.temporal.io/docs/server/quick-install/)。
3. 配置环境变量：
    ```bash
       export OPENROUTER_API_KEY=your_openrouter_api_key
    ```
4. 启动worker：
    ```bash
       go run ./worker/main.go
    ```
5. 启动client，提交workflow： 
    ```bash
       go run ./client/main.go
    ```
6. 此时qdrant数据库中会保存被抓去到的新闻的向量数据
7. 启动rag_search，进行RAG检索对话：
    ```bash
       go run ./client/rag_search.go
    ```
   
## 编辑器

推荐使用[GoLand](https://www.jetbrains.com/go/)作为编辑器，支持Go语言的开发和调试。