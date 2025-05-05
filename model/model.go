package model

// MinedProject 表示一个挖掘的 GitHub 项目
type MinedProject struct {
	Name        string
	Owner       string
	Description string
	URL         string
	Language    string
	Stars       int
	Readme      string
	Analysis    string
}

// ProjectFetcher 获取GitHub流行项目的接口
type ProjectFetcher interface {
	// AddCallback 增加回调函数，每次获取到项目时调用
	AddCallback(callback CallbackFunc)

	// Start 开始执行
	Start() error
}

// CallbackFunc 回调函数类型定义
type CallbackFunc = func(projects []MinedProject) error

// ReadmeAnalyzer 分析README的接口
type ReadmeAnalyzer interface {
	Analyze(project *MinedProject) error
}

type OpenAIConfig struct {
	BaseURL string
	Model   string
	APIKey  string
}

// ProjectStorage 项目存储接口
type ProjectStorage interface {
	// 保存项目列表
	SaveProjects(projects []MinedProject) error

	// 加载所有历史项目
	LoadAllProjects() ([]MinedProject, error)

	// 检查项目是否已处理过
	IsProjectProcessed(owner, name string) bool
}
