package config

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
	// "github.com/shinya-ac/TodoAPI/pkg/logging"
)

type ConfigList struct {
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	ServerAddress string
	ServerPort    string
	APIKey1       string
	APIKey2       string
	APIKey3       string
}

var Config ConfigList

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func LoadConfig() (ConfigList, error) {
	var cfg *ini.File
	var err error

	cfg, err = ini.Load("config.ini")
	if err != nil {
		// logging.Logger.Warn("config.ini の読み込みに失敗。環境変数から設定を読み込む。", "error", err)
	}

	missingConfig := []string{}

	Config = ConfigList{
		DBUser:        getEnv("DB_USER", getINIValue(cfg, "db", "user", "")),
		DBPassword:    getEnv("DB_PASSWORD", getINIValue(cfg, "db", "password", "")),
		DBHost:        getEnv("DB_HOST", getINIValue(cfg, "db", "host", "")),
		DBPort:        getEnv("DB_PORT", getINIValue(cfg, "db", "port", "")),
		DBName:        getEnv("DB_NAME", getINIValue(cfg, "db", "name", "")),
		ServerAddress: getEnv("SERVER_ADDRESS", getINIValue(cfg, "server", "address", "")),
		ServerPort:    getEnv("SERVER_PORT", getINIValue(cfg, "server", "port", "")),
		APIKey1:       getEnv("API_KEY1", getINIValue(cfg, "api", "key1", "")),
		APIKey2:       getEnv("API_KEY2", getINIValue(cfg, "api", "key2", "")),
		APIKey3:       getEnv("API_KEY3", getINIValue(cfg, "api", "key3", "")),
	}

	if Config.DBUser == "" {
		missingConfig = append(missingConfig, "DB_USER")
	}
	if Config.DBPassword == "" {
		missingConfig = append(missingConfig, "DB_PASSWORD")
	}
	if Config.DBHost == "" {
		missingConfig = append(missingConfig, "DB_HOST")
	}
	if Config.DBPort == "" {
		missingConfig = append(missingConfig, "DB_PORT")
	}
	if Config.DBName == "" {
		missingConfig = append(missingConfig, "DB_NAME")
	}
	if Config.ServerAddress == "" {
		missingConfig = append(missingConfig, "SERVER_ADDRESS")
	}
	if Config.ServerPort == "" {
		missingConfig = append(missingConfig, "SERVER_PORT")
	}
	if Config.APIKey1 == "" {
		missingConfig = append(missingConfig, "API_KEY1")
	}
	if Config.APIKey2 == "" {
		missingConfig = append(missingConfig, "API_KEY2")
	}
	if Config.APIKey3 == "" {
		missingConfig = append(missingConfig, "API_KEY3")
	}

	if len(missingConfig) > 0 {
		errMsg := fmt.Sprintf("必要な設定が見つかりません: %v", missingConfig)
		// logging.Logger.Error(errMsg)
		return Config, fmt.Errorf(errMsg)
	}

	return Config, nil
}

func getINIValue(cfg *ini.File, section string, key string, fallback string) string {
	if cfg == nil {
		return fallback
	}
	return cfg.Section(section).Key(key).String()
}
