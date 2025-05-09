# GitHub Miner

GitHub Miner 是一个自动挖掘 GitHub 流行项目的工具，能定期获取热门开源项目，利用 AI 分析其内容，并通过飞书消息推送，帮助开发者了解开源世界的最新动态。

## 功能特点

- 🔍 **自动发现**：定时抓取 GitHub 趋势榜上的热门项目
- 📄 **内容获取**：自动获取每个项目的 README 文件
- 🤖 **AI 分析**：利用 AI 技术分析项目内容，提取关键信息
- 🔄 **去重处理**：自动过滤已处理过的项目，只关注新内容
- 💾 **数据存储**：将分析结果保存为 JSON 文件，便于后续查询
- 📱 **消息处理**：支持通过注册 callback 函数处理分析结果

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

### 安装

```bash
# 安装依赖
go get github.com/yulai-123/github_miner
```

## 运行效果

![项目运行效果](https://github.com/yulai-123/github_miner/blob/main/project_result.jpg)

## 许可证

MIT License