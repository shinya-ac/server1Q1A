package config

import (
	"fmt"
	"os"

	"github.com/shinya-ac/server1Q1A/pkg/logging"
	"gopkg.in/ini.v1"
)

type ConfigList struct {
	DBUser               string
	DBPassword           string
	DBHost               string
	DBPort               string
	DBName               string
	CACertPath           string
	ServerAddress        string
	ServerPort           string
	APIKey1              string
	APIKey2              string
	APIKey3              string
	AuthDomain           string
	AuthClientID         string
	SigningMethod        string
	DBTLSMode            string
	ChatGptApiKey        string
	MicrocmsApiKey       string
	MicrocmsOcrContentID string
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
		logging.Logger.Warn("config.ini の読み込みに失敗。環境変数から設定を読み込む。", "error", err)
	}

	missingConfig := []string{}

	Config = ConfigList{
		DBUser:               getEnv("DB_USER", getINIValue(cfg, "db", "user", "")),
		DBPassword:           getEnv("DB_PASSWORD", getINIValue(cfg, "db", "password", "")),
		DBHost:               getEnv("DB_HOST", getINIValue(cfg, "db", "host", "")),
		DBPort:               getEnv("DB_PORT", getINIValue(cfg, "db", "port", "")),
		DBName:               getEnv("DB_NAME", getINIValue(cfg, "db", "name", "")),
		CACertPath:           getEnv("CA_CERT_PATH", getINIValue(cfg, "db", "ca_cert_path", "")),
		ServerAddress:        getEnv("SERVER_ADDRESS", getINIValue(cfg, "server", "address", "")),
		ServerPort:           getEnv("SERVER_PORT", getINIValue(cfg, "server", "port", "")),
		APIKey1:              getEnv("API_KEY1", getINIValue(cfg, "api", "key1", "")),
		APIKey2:              getEnv("API_KEY2", getINIValue(cfg, "api", "key2", "")),
		APIKey3:              getEnv("API_KEY3", getINIValue(cfg, "api", "key3", "")),
		AuthDomain:           getEnv("AUTH_DOMAIN", getINIValue(cfg, "auth0", "auth_domain", "")),
		AuthClientID:         getEnv("AUTH_CLIENT_ID", getINIValue(cfg, "auth0", "auth_client_id", "")),
		SigningMethod:        getEnv("JWT_SIGNING_METHOD", getINIValue(cfg, "auth0", "signing_method", "")),
		DBTLSMode:            getEnv("DB_TLS_MODE", getINIValue(cfg, "db", "tls_mode", "")),
		ChatGptApiKey:        getEnv("CHAT_GPT_API_KEY", getINIValue(cfg, "chatgpt", "api_key", "")),
		MicrocmsApiKey:       getEnv("MICROCMS_API_KEY", getINIValue(cfg, "microcms", "api_key", "")),
		MicrocmsOcrContentID: getEnv("MICROCMS_OCR_CONTENT_ID", getINIValue(cfg, "microcms", "ocr_content_id", "")),
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
	if Config.CACertPath == "" {
		missingConfig = append(missingConfig, "CA_CERT_PATH")
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
	if Config.AuthDomain == "" {
		missingConfig = append(missingConfig, "AUTH_DOMAIN")
	}
	if Config.AuthClientID == "" {
		missingConfig = append(missingConfig, "AUTH_CLIENT_ID")
	}
	if Config.SigningMethod == "" {
		missingConfig = append(missingConfig, "JWT_SIGNING_METHOD")
	}
	if Config.DBTLSMode == "" {
		missingConfig = append(missingConfig, "DB_TLS_MODE")
	}
	if Config.ChatGptApiKey == "" {
		missingConfig = append(missingConfig, "CHAT_GPT_API_KEY")
	}
	if Config.MicrocmsApiKey == "" {
		missingConfig = append(missingConfig, "MICROCMS_API_KEY")
	}
	if Config.MicrocmsOcrContentID == "" {
		missingConfig = append(missingConfig, "MICROCMS_OCR_CONTENT_ID")
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
