package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"ssd-assignment-api/models"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type SpecificConfigService struct {
	configs map[string]models.SpecificConfig
	mutex   sync.Mutex
	yamlDir string
}

func NewSpecificConfigService(yamlDir string) (*SpecificConfigService, error) {
	service := &SpecificConfigService{
		configs: make(map[string]models.SpecificConfig),
		yamlDir: yamlDir,
	}

	if err := service.loadConfigsFromYAML(); err != nil {
		return nil, fmt.Errorf("failed to load specific configs: %w", err)
	}

	return service, nil
}

func (s *SpecificConfigService) GetMatchingConfigs(host, url, page string) ([]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	priorityMap := make(map[string]int)
	foundMatches := false

	for _, config := range s.configs {
		if host != "" {
			if ids, exists := config.DataSource.Hosts[host]; exists {
				foundMatches = true
				for _, id := range ids {
					priorityMap[id] += 3
				}
			}
		}

		if url != "" {
			if ids, exists := config.DataSource.URLs[url]; exists {
				foundMatches = true
				for _, id := range ids {
					priorityMap[id] += 2
				}
			}
		}

		if page != "" {
			if ids, exists := config.DataSource.Pages[page]; exists {
				foundMatches = true
				for _, id := range ids {
					priorityMap[id] += 1
				}
			}
		}
	}

	if !foundMatches {
		return nil, errors.New("no matching configurations found")
	}

	return s.sortConfigsByPriority(priorityMap), nil
}

func (s *SpecificConfigService) sortConfigsByPriority(priorityMap map[string]int) []string {
	type configPriority struct {
		ID    string
		Score int
	}

	var sorted []configPriority
	for id, score := range priorityMap {
		sorted = append(sorted, configPriority{id, score})
	}

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Score == sorted[j].Score {
			return sorted[i].ID < sorted[j].ID
		}
		return sorted[i].Score > sorted[j].Score
	})

	result := make([]string, len(sorted))
	for i, item := range sorted {
		result[i] = item.ID
	}

	return result
}

func (s *SpecificConfigService) loadConfigsFromYAML() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	files, err := os.ReadDir(s.yamlDir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		filePath := filepath.Join(s.yamlDir, file.Name())
		yamlData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		var config models.SpecificConfig
		if err := yaml.Unmarshal(yamlData, &config); err != nil {
			return fmt.Errorf("YAML parse hatası: %w", err)
		}

		s.configs[config.ID] = config
	}
	return nil
}

func (s *SpecificConfigService) GetAllSpecificConfigs() ([]models.SpecificConfig, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var configList []models.SpecificConfig
	for _, config := range s.configs {
		configList = append(configList, config)
	}
	return configList, nil
}

func (s *SpecificConfigService) GetSpecificConfigByID(id string) (models.SpecificConfig, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	config, exists := s.configs[id]
	if !exists {
		return models.SpecificConfig{}, errors.New("specific config not found")
	}
	return config, nil
}

func (s *SpecificConfigService) AddSpecificConfig(config models.SpecificConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.configs[config.ID]; exists {
		return fmt.Errorf("config with ID '%s' already exists", config.ID)
	}

	filePath := fmt.Sprintf("%s/%s.yaml", s.yamlDir, config.ID)
	if err := s.saveConfigToYAML(config, filePath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	s.configs[config.ID] = config
	return nil
}

func (s *SpecificConfigService) UpdateSpecificConfig(id string, config models.SpecificConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.configs[id]; !exists {
		return errors.New("specific config not found")
	}

	filePath := fmt.Sprintf("%s/%s.yaml", s.yamlDir, id)
	if err := s.saveConfigToYAML(config, filePath); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	s.configs[id] = config
	return nil
}

func (s *SpecificConfigService) DeleteSpecificConfig(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.configs[id]; !exists {
		return errors.New("specific config not found")
	}

	filePath := fmt.Sprintf("%s/%s.yaml", s.yamlDir, id)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete config file: %w", err)
	}

	delete(s.configs, id)
	return nil
}

func (s *SpecificConfigService) saveConfigToYAML(config models.SpecificConfig, path string) error {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(path, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
