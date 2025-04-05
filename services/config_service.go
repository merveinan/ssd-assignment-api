package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"ssd-assignment-api/models"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

type ConfigService struct {
	configs map[string]models.Config
	mutex   sync.Mutex
	yamlDir string // Only the YAML directory will be stored
}

// NewConfigService now only accepts the YAML directory
func NewConfigService(yamlDir string) (*ConfigService, error) {
	service := &ConfigService{
		configs: make(map[string]models.Config),
		yamlDir: yamlDir,
	}

	if err := service.loadConfigsFromYAML(); err != nil {
		return nil, fmt.Errorf("YAML loading error: %v", err)
	}

	return service, nil
}

// GetAllConfigs retrieves all configurations
func (s *ConfigService) GetAllConfigs() ([]models.Config, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var configList []models.Config
	for _, config := range s.configs {
		configList = append(configList, config)
	}

	return configList, nil
}
func (s *ConfigService) loadConfigsFromYAML() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Read all files in the YAML directory
	files, err := os.ReadDir(s.yamlDir)
	if err != nil {
		return fmt.Errorf("YAML directory could not be read: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue // Process only .yaml files
		}

		// Create the file path
		filePath := filepath.Join(s.yamlDir, file.Name())
		yamlData, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("%s file could not be read: %w", filePath, err)
		}

		// Unmarshal YAML data to Config struct
		var config models.Config
		if err := yaml.Unmarshal(yamlData, &config); err != nil {
			return fmt.Errorf("%s file could not be parsed: %w", filePath, err)
		}

		// Add Config to memory
		s.configs[config.ID] = config
	}

	return nil
}

// GetConfigByID retrieves a configuration by its ID
func (s *ConfigService) GetConfigByID(id string) (models.Config, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	config, exists := s.configs[id]
	if !exists {
		return models.Config{}, errors.New("config not found")
	}
	return config, nil
}

// AddConfig creates a YAML file and adds it to memory
func (s *ConfigService) AddConfig(config models.Config) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// ID check
	if _, exists := s.configs[config.ID]; exists {
		return fmt.Errorf("ID '%s' already exists", config.ID)
	}

	// Create YAML file
	if err := s.createYAMLFile(config); err != nil {
		return fmt.Errorf("YAML file could not be created: %w", err)
	}

	// Add to memory
	s.configs[config.ID] = config
	return nil
}

// UpdateConfig updates the YAML file and changes it in memory
func (s *ConfigService) UpdateConfig(id string, config models.Config) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// ID check
	if _, exists := s.configs[id]; !exists {
		return errors.New("configuration not found")
	}

	// Update YAML file
	if err := s.updateYAMLFile(config); err != nil {
		return fmt.Errorf("YAML could not be updated: %w", err)
	}

	// Update memory
	s.configs[id] = config
	return nil
}

// DeleteConfig deletes the YAML file and removes it from memory
func (s *ConfigService) DeleteConfig(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// ID check
	if _, exists := s.configs[id]; !exists {
		return errors.New("configuration not found")
	}

	// Delete YAML file
	if err := s.deleteYAMLFile(id); err != nil {
		return fmt.Errorf("YAML could not be deleted: %w", err)
	}

	// Remove from memory
	delete(s.configs, id)
	return nil
}

// createYAMLFile creates a YAML file for the new config based on the config's ID
func (s *ConfigService) createYAMLFile(config models.Config) error {
	filePath := fmt.Sprintf("%s/%s.yaml", s.yamlDir, config.ID)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshall config to YAML
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write the YAML data to file
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	return nil
}

// updateYAMLFile updates the YAML file corresponding to the config's ID
func (s *ConfigService) updateYAMLFile(config models.Config) error {
	filePath := fmt.Sprintf("%s/%s.yaml", s.yamlDir, config.ID)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshall config to YAML
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Write the YAML data to file
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	return nil
}

// deleteYAMLFile deletes the YAML file corresponding to the config's ID
func (s *ConfigService) deleteYAMLFile(id string) error {
	filePath := fmt.Sprintf("%s/%s.yaml", s.yamlDir, id)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
