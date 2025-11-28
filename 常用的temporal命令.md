1. 启动temporal

```shell
temporal server start-dev --db-filename your_temporal.db --ui-port 8080
```

2. 发送信号

```shell
temporal workflow signal --query 'WorkflowType="NewsAgentWorkflow"' --name ReviewSignal --input '{"Action": "APPROVE"}'
```

3. 查看namespace的配置

```shell
temporal operator namespace describe default
```

4. 修改namespace的配置

```shell
temporal operator namespace update --retention 7d default
```