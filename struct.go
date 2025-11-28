package HumanintheLoopAI

// 协议：Researcher Agent 的输入输出
type ResearchRequest struct {
	Topic string
	Depth string // "deep" or "quick"
}

type ResearchResult struct {
	Facts   []string
	Sources []string
}

type WriteRequest struct {
	Topic string
	Facts []string
	Tone  string // "formal" or "casual"
}

type WriteResult struct {
	Content string // HTML string
}
