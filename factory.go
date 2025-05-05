package github_miner

import (
	"github.com/yulai-123/github_miner/impl/ai_analyzer"
	"github.com/yulai-123/github_miner/impl/github_trending"
	"github.com/yulai-123/github_miner/impl/storage"
	"github.com/yulai-123/github_miner/model"
	"time"
)

type ProjectMinerFactory struct {
	fetchers []model.ProjectFetcher
	storage  model.ProjectStorage
}

func NewProjectMinerFactory(storagePath string) (*ProjectMinerFactory, error) {
	// 创建JSON存储
	jsonStorage, err := storage.NewJSONStorage(storagePath)
	if err != nil {
		return nil, err
	}

	return &ProjectMinerFactory{
		fetchers: make([]model.ProjectFetcher, 0),
		storage:  jsonStorage,
	}, nil
}

func (f *ProjectMinerFactory) RegisterGitHubTrending(
	period string,
	language string,
	interval time.Duration,
	githubAPIKey string,
	aiConfig model.OpenAIConfig,
	callbacks []model.CallbackFunc,
) error {
	readmeClient := github_trending.NewGitHubReadmeClient(githubAPIKey)
	aiAnalyzer := ai_analyzer.NewOpenAIAnalyzer(aiConfig)

	// 传入storage
	fetcher := github_trending.NewProjectFetcher(period, language, interval, readmeClient, aiAnalyzer, f.storage)

	for _, callback := range callbacks {
		fetcher.AddCallback(callback)
	}

	f.fetchers = append(f.fetchers, fetcher)
	return nil
}

func (f *ProjectMinerFactory) Start() error {
	for _, fetcher := range f.fetchers {
		if err := fetcher.Start(); err != nil {
			return err
		}
	}
	return nil
}
