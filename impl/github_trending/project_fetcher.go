package github_trending

import (
	"github.com/andygrunwald/go-trending"
	"github.com/sirupsen/logrus"
	"github.com/yulai-123/github_miner/model"
	"strings"
	"time"
)

type ProjectFetcher struct {
	trendingClient *trending.Trending
	interval       time.Duration
	callbackList   []model.CallbackFunc
	readmeClient   *GitHubReadmeClient
	analyzerClient model.ReadmeAnalyzer
	storage        model.ProjectStorage
}

// 修改构造函数
func NewProjectFetcher(interval time.Duration,
	readmeClient *GitHubReadmeClient,
	analyzerClient model.ReadmeAnalyzer,
	storage model.ProjectStorage) *ProjectFetcher {
	return &ProjectFetcher{
		trendingClient: trending.NewTrending(),
		interval:       interval,
		readmeClient:   readmeClient,
		analyzerClient: analyzerClient,
		storage:        storage,
	}
}

func (p *ProjectFetcher) AddCallback(callback model.CallbackFunc) {
	p.callbackList = append(p.callbackList, callback)
}

func (p *ProjectFetcher) Start() error {
	go p.fetchPeriodically()
	return nil
}

func (p *ProjectFetcher) fetchPeriodically() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	// 立即执行一次
	p.fetchAndProcess()

	for range ticker.C {
		p.fetchAndProcess()
	}
}

// 修改 fetchAndProcess 方法
func (p *ProjectFetcher) fetchAndProcess() {
	projects, err := p.fetchTrendingProjects()
	if err != nil {
		logrus.Errorf("获取GitHub流行项目失败: %v", err)
		return
	}

	logrus.Infof("[fetchAndProcess] 成功拉取 %v 个项目", len(projects))

	// 过滤已处理过的项目
	var newProjects []model.MinedProject
	for _, project := range projects {
		if !p.storage.IsProjectProcessed(project.Owner, project.Name) {
			newProjects = append(newProjects, project)
		}
	}

	if len(newProjects) == 0 {
		logrus.Info("没有新的项目需要处理")
		return
	}

	logrus.Infof("发现 %d 个新项目需要处理", len(newProjects))

	// 获取README并分析
	processedProjects := make([]model.MinedProject, 0, len(newProjects))
	for i := range newProjects {
		if err := p.readmeClient.FetchReadme(&newProjects[i]); err != nil {
			logrus.Errorf("获取项目 %s README失败: %v", newProjects[i].Name, err)
			continue
		}

		if err := p.analyzerClient.Analyze(&newProjects[i]); err != nil {
			logrus.Errorf("分析项目 %s README失败: %v", newProjects[i].Name, err)
			continue
		}

		processedProjects = append(processedProjects, newProjects[i])
	}

	// 保存处理过的项目
	if len(processedProjects) > 0 {
		if err := p.storage.SaveProjects(processedProjects); err != nil {
			logrus.Errorf("保存项目数据失败: %v", err)
		}
	}

	// 执行回调
	if len(processedProjects) > 0 {
		for _, callback := range p.callbackList {
			for _, project := range processedProjects {
				if err := callback(project); err != nil {
					logrus.Errorf("执行回调失败: %v", err)
				}
			}
		}
	}
}

func (p *ProjectFetcher) fetchTrendingProjects() ([]model.MinedProject, error) {
	projects, err := p.trendingClient.GetProjects("", "")
	if err != nil {
		return nil, err
	}

	result := make([]model.MinedProject, 0, len(projects))
	for _, proj := range projects {
		// proj 返回的 name 格式为 owner/repo 的格式，我需要解析出 repo来
		ownerRepo := strings.Split(proj.Name, "/")
		if len(ownerRepo) != 2 {
			logrus.Errorf("项目名称格式错误: %s", proj.Name)
			continue
		}
		result = append(result, model.MinedProject{
			Name:        ownerRepo[1],
			Owner:       proj.Owner,
			Description: proj.Description,
			URL:         proj.URL.String(),
			Language:    proj.Language,
			Stars:       proj.Stars,
		})
	}

	return result, nil
}
