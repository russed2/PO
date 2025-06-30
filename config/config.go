package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Структура для YAML конфигурации
type Config struct {
	App struct {
		Env     string `yaml:"env"`
		Port    int    `yaml:"port"`
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"app"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Logging struct {
		Level  string `yaml:"level"`
		File   string `yaml:"file"`
		Format string `yaml:"format"`
	} `yaml:"logging"`

	JWT struct {
		Secret      string `yaml:"secret"`
		ExpiryHours int    `yaml:"expiry_hours"`
	} `yaml:"jwt"`

	Server struct {
		ReadTimeout    string `yaml:"read_timeout"`
		WriteTimeout   string `yaml:"write_timeout"`
		MaxHeaderBytes int    `yaml:"max_header_bytes"`
	} `yaml:"server"`
}

// LoadConfig загружает переменные окружения из .env файла
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println("Loaded configuration from .env file")
	}
}

// GetEnv получает переменную окружения с значением по умолчанию
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// LoadYAMLConfig загружает конфигурацию из YAML файла
func LoadYAMLConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
