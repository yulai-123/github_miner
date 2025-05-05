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

	prompt := fmt.Sprintf(`你是一位技术分析师，请基于以下GitHub项目的README内容，提供简洁可靠的分析摘要。

项目信息:
- 名称: %s/%s
- 描述: %s
- 主要语言: %s
- README内容:
%s

请严格根据README中实际包含的信息，提供以下分析（如果某项信息README中没有明确提及，请标注"未提供"而不要猜测）：

## 📝 项目概述
[1-2段简洁描述，概括项目的主要目的和特点]

## 📋 核心功能
[3-5点简明列表，说明这个项目解决什么问题/提供什么功能]

## ⚙️ 技术实现
[简要说明项目的技术栈和关键实现方式，只列出README中明确提到的技术]

## 💡 开发者价值
[1-2段简洁描述，说明这个项目对独立开发者的参考价值]

## 🚀 潜在机会
[如果有的话，简述这个项目反映的市场需求或发展趋势]

请保持内容简洁、客观，每个部分不超过3-5行，只包含README中明确提供的信息。`,
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
