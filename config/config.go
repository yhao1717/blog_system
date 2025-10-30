package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
	App      AppConfig      `yaml:"app"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type JWTConfig struct {
	Secret          string `yaml:"secret"`
	ExpirationHours int    `yaml:"expiration_hours"`
}

type LogConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果未指定配置文件路径，使用默认路径
	if configPath == "" {
		configPath = "config.yaml"
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file %s does not exist", configPath)
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// 解析YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// 设置默认值
	config.setDefaults()

	// 验证配置
	if err := config.validate(); err != nil {
		return nil, err
	}

	log.Printf("Configuration loaded from %s", configPath)
	return &config, nil
}

// setDefaults 设置配置的默认值
func (c *Config) setDefaults() {
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}
	if c.Database.Driver == "" {
		c.Database.Driver = "sqlite"
	}
	if c.Database.DSN == "" {
		c.Database.DSN = "blog.db"
	}
	if c.JWT.Secret == "" {
		c.JWT.Secret = "default_jwt_secret_change_in_production"
	}
	if c.JWT.ExpirationHours == 0 {
		c.JWT.ExpirationHours = 24
	}
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.App.Name == "" {
		c.App.Name = "Blog System"
	}
	if c.App.Version == "" {
		c.App.Version = "1.0.0"
	}
}

// validate 验证配置
func (c *Config) validate() error {
	// 验证服务器模式
	validModes := map[string]bool{
		"debug":   true,
		"release": true,
		"test":    true,
	}
	if !validModes[c.Server.Mode] {
		return fmt.Errorf("invalid server mode: %s. Must be one of: debug, release, test", c.Server.Mode)
	}

	// 验证数据库驱动
	validDrivers := map[string]bool{
		"sqlite": true,
		"mysql":  true,
	}
	if !validDrivers[c.Database.Driver] {
		return fmt.Errorf("invalid database driver: %s. Must be one of: sqlite, mysql", c.Database.Driver)
	}

	// 验证JWT密钥
	if c.JWT.Secret == "default_jwt_secret_change_in_production" {
		log.Println("WARNING: Using default JWT secret. Please change it in production!")
	}

	return nil
}

func (c *DatabaseConfig) GetMySQLDSN() string {
	return c.DSN
}
