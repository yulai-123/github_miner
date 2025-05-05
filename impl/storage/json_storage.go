package storage

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yulai-123/github_miner/model"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type JSONStorage struct {
	storagePath       string
	processedProjects map[string]time.Time // 改为存储项目和其处理时间
	mutex             sync.RWMutex
}

func NewJSONStorage(storagePath string) (*JSONStorage, error) {
	// 确保存储目录存在
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, err
	}

	storage := &JSONStorage{
		storagePath:       storagePath,
		processedProjects: make(map[string]time.Time),
	}

	// 加载历史数据
	if err := storage.loadHistoricalData(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *JSONStorage) loadHistoricalData() error {
	files, err := filepath.Glob(filepath.Join(s.storagePath, "*.json"))
	if err != nil {
		return err
	}

	for _, file := range files {
		// 从文件名中解析日期
		baseName := filepath.Base(file)
		var year, month, day int
		_, err := fmt.Sscanf(baseName, "projects_%4d-%2d-%2d.json", &year, &month, &day)
		if err != nil {
			logrus.Warnf("无法从文件名解析日期: %s, %v", baseName, err)
			continue
		}
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)

		// 解析处理日期
		processDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			logrus.Warnf("无法解析日期字符串: %s, %v", dateStr, err)
			continue
		}

		data, err := os.ReadFile(file)
		if err != nil {
			logrus.Errorf("读取文件失败 %s: %v", file, err)
			continue
		}

		var projects []model.MinedProject
		if err := json.Unmarshal(data, &projects); err != nil {
			logrus.Errorf("解析JSON文件失败 %s: %v", file, err)
			continue
		}

		for _, project := range projects {
			key := fmt.Sprintf("%s/%s", project.Owner, project.Name)
			// 如果同一个项目在多个日期处理过，保留最新日期
			if existingDate, exists := s.processedProjects[key]; !exists || processDate.After(existingDate) {
				s.processedProjects[key] = processDate
			}
		}
	}

	logrus.Infof("已加载 %d 个历史处理过的项目", len(s.processedProjects))
	return nil
}

func (s *JSONStorage) SaveProjects(projects []model.MinedProject) error {
	if len(projects) == 0 {
		return nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 生成当前日期的文件名
	today := time.Now().Format("2006-01-02")
	filename := filepath.Join(s.storagePath, fmt.Sprintf("projects_%s.json", today))
	processDate, _ := time.Parse("2006-01-02", today)

	// 如果文件已存在，先读取合并
	var existingProjects []model.MinedProject
	if _, err := os.Stat(filename); err == nil {
		data, err := os.ReadFile(filename)
		if err == nil {
			if err := json.Unmarshal(data, &existingProjects); err != nil {
				logrus.Warnf("解析现有JSON文件失败: %v", err)
			}
		}
	}

	// 合并项目列表并去重
	projectMap := make(map[string]model.MinedProject)
	for _, p := range existingProjects {
		key := fmt.Sprintf("%s/%s", p.Owner, p.Name)
		projectMap[key] = p
	}

	for _, p := range projects {
		key := fmt.Sprintf("%s/%s", p.Owner, p.Name)
		projectMap[key] = p
		s.processedProjects[key] = processDate
	}

	// 转换回列表
	mergedProjects := make([]model.MinedProject, 0, len(projectMap))
	for _, p := range projectMap {
		mergedProjects = append(mergedProjects, p)
	}

	// 保存到文件
	data, err := json.MarshalIndent(mergedProjects, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (s *JSONStorage) LoadAllProjects() ([]model.MinedProject, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var allProjects []model.MinedProject
	files, err := filepath.Glob(filepath.Join(s.storagePath, "*.json"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			logrus.Errorf("读取文件失败 %s: %v", file, err)
			continue
		}

		var projects []model.MinedProject
		if err := json.Unmarshal(data, &projects); err != nil {
			logrus.Errorf("解析JSON文件失败 %s: %v", file, err)
			continue
		}

		allProjects = append(allProjects, projects...)
	}

	return allProjects, nil
}

func (s *JSONStorage) IsProjectProcessed(owner, name string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	key := fmt.Sprintf("%s/%s", owner, name)
	processTime, exists := s.processedProjects[key]
	if !exists {
		return false
	}

	// 检查是否在一周内处理过
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	return processTime.After(oneWeekAgo)
}
