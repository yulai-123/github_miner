package ai_analyzer

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"github.com/yulai-123/github_miner/model"
)

type OpenAIAnalyzer struct {
	client *openai.Client
	model  string
}

func NewOpenAIAnalyzer(param model.OpenAIConfig) *OpenAIAnalyzer {
	config := openai.DefaultConfig(param.APIKey)
	config.BaseURL = param.BaseURL

	return &OpenAIAnalyzer{
		client: openai.NewClientWithConfig(config),
		model:  param.Model,
	}
}

func (o *OpenAIAnalyzer) Analyze(project *model.MinedProject) error {
	ctx := context.Background()

	prompt := fmt.Sprintf(`作为独立开发者的顾问，请分析以下GitHub项目README文件，并提供结构化的分析报告。

项目信息:
- 名称: %s/%s
- 描述: %s
- 主要语言: %s
- README内容:
%s

请提供以下格式的分析报告：

## 📌 项目概述
[简明扼要地描述项目的核心目的和解决的问题]

## 💡 核心功能
[以要点形式列出项目的主要功能]

## 🔧 技术实现
[分析项目使用的技术栈、架构设计和关键实现方法]

## 🎯 适用场景
[分析项目适合的使用场景和目标用户]

## 👨‍💻 开发者价值
[从独立开发者角度，分析项目可能给你带来的启发、可复用的组件或解决方案]

## 💼 市场洞察
[项目可能反映的市场需求或发展趋势]

## 🔍 亮点与特色
[项目中值得关注的创新点或独特优势]

请确保分析内容准确、有深度，并基于README内容给出有理有据的分析。内容应当精炼实用，避免冗余信息。`,
		project.Owner, project.Name, project.Description, project.Language, project.Readme)

	resp, err := o.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: o.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 5000,
		},
	)

	if err != nil {
		// 添加相关日志，记录请求的项目、输入的Token和输出的Token，错误原因
		logrus.Errorf("分析项目失败: 项目 %s/%s, 错误原因: %v", project.Owner, project.Name, err)
		return err
	}

	project.Analysis = resp.Choices[0].Message.Content

	// 添加相关日志，记录请求的项目、输入的Token和输出的Token
	logrus.Infof("分析完成: 项目 %s/%s, 输入Token: %d, 输出Token: %d, 总Token: %d",
		project.Owner, project.Name, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	return nil
}
