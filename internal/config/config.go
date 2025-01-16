package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey      string   `json:"api_key"`
	IPWhitelist []string `json:"ip_whitelist,omitempty"`
}

const (
	configDir  = "/etc/hy2agent"
	configFile = "config.json"
)

// 生成随机API Key
func generateAPIKey() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 加载或创建配置
func LoadConfig() (*Config, error) {
	configPath := filepath.Join(configDir, configFile)

	// 检查配置目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return nil, err
		}
	}

	// 尝试读取配置文件
	if data, err := os.ReadFile(configPath); err == nil {
		var config Config
		if err := json.Unmarshal(data, &config); err == nil {
			return &config, nil
		}
	}

	// 如果配置文件不存在或无效，创建新配置
	config := &Config{
		APIKey: generateAPIKey(),
	}

	// 保存配置
	if err := SaveConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// 保存配置到文件
func SaveConfig(cfg *Config) error {
	configPath := filepath.Join(configDir, configFile)

	// 将配置转换为JSON
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(configPath, data, 0644)
}
