# LLM-as-a-Judge模式

我们将实现一个简单的 Ragas (Retrieval Augmented Generation Assessment) 核心逻辑。我们要评估两个指标：

Faithfulness (忠实度)：AI 的回答是否只基于检索到的上下文？（有没有瞎编？）

Answer Relevance (相关度)：AI 的回答是否解决了用户的问题？