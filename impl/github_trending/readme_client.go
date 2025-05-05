package github_trending

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/yulai-123/github_miner/model"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type GitHubReadmeClient struct {
	client *github.Client
}

func NewGitHubReadmeClient(token string) *GitHubReadmeClient {
	var httpClient *http.Client

	if token != "" {
		// 使用OAuth2认证
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		httpClient = oauth2.NewClient(context.Background(), ts)
	}

	return &GitHubReadmeClient{
		client: github.NewClient(httpClient),
	}
}

func (g *GitHubReadmeClient) FetchReadme(project *model.MinedProject) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	readme, _, err := g.client.Repositories.GetReadme(ctx, project.Owner, project.Name, nil)
	if err != nil {
		return err
	}

	content, err := readme.GetContent()
	if err != nil {
		return err
	}

	project.Readme = content
	// 添加日志，记录请求的项目，readme长度
	logrus.Infof("成功获取项目 %s/%s 的 README，长度为 %d 字符", project.Owner, project.Name, len(content))
	return nil
}
