# GitHub Miner

GitHub Miner 是一个小型信息处理工具，用来定期获取 GitHub Trending 项目，拉取 README，交给 LLM 做结构化摘要，并通过 callback 交给后续流程处理。

它的目标不是替代人工判断，而是把“每天扫一遍热门项目”这件事变成可持续的信息流。

## 功能特点

- **自动发现**：定时抓取 GitHub Trending 上的热门项目
- **内容获取**：自动读取项目 README
- **AI 分析**：用 OpenAI-compatible API 生成项目摘要
- **去重处理**：过滤已经处理过的项目
- **数据存储**：把结果保存为 JSON，方便后续查询
- **消息处理**：通过 callback 把分析结果接入飞书、邮件或其他通知系统

## 项目架构

项目采用模块化设计，主要组件包括：

- `ProjectFetcher`：负责从 GitHub 获取流行项目信息
- `GitHubReadmeClient`：负责获取项目的 README 文件
- `OpenAIAnalyzer`：使用 AI 分析 README 内容
- `ProjectStorage`：处理项目数据的存储和检索

## 安装使用

### 前置条件

- Go 1.16+
- GitHub API 密钥（用于获取 README）
- OpenAI 或 DeepSeek API 密钥（用于 AI 分析）

### 安装依赖

```bash
go mod tidy
```

### 使用方式

这个仓库目前更像一个 Go package / building block。你可以在自己的程序里组合：

- `ProjectFetcher`
- `GitHubReadmeClient`
- `OpenAIAnalyzer`
- `ProjectStorage`
- callback handler

然后把分析结果接到自己的通知或存储系统里。

## 运行效果

![项目运行效果](https://github.com/yulai-123/github_miner/blob/main/project_result.jpg)

## Notes

这个项目是个人 AI 信息流实验的一部分，重点在于信息获取、摘要、去重和 callback 流程。它不是一个完整的 SaaS 产品，也不会保证覆盖所有 GitHub Trending 边界情况。

## 许可证

MIT License
