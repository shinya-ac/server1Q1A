package config_test

import (
	"os"
	"testing"

	configs "github.com/shinya-ac/server1Q1A/configs"
	"github.com/stretchr/testify/assert"
	// "github.com/shinya-ac/TodoAPI/pkg/logging"
	// "github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// logging.InitLogger()

	t.Run("正常系 - 環境変数から設定を読み込む", func(t *testing.T) {
		os.Setenv("DB_USER", "test_user")
		os.Setenv("DB_PASSWORD", "test_password")
		os.Setenv("DB_HOST", "localhost")
		os.Setenv("DB_PORT", "3306")
		os.Setenv("DB_NAME", "test_db")
		os.Setenv("SERVER_ADDRESS", "127.0.0.1")
		os.Setenv("SERVER_PORT", "8080")
		os.Setenv("API_KEY1", "key1")
		os.Setenv("API_KEY2", "key2")
		os.Setenv("API_KEY3", "key3")
		defer func() {
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("SERVER_ADDRESS")
			os.Unsetenv("SERVER_PORT")
			os.Unsetenv("API_KEY1")
			os.Unsetenv("API_KEY2")
			os.Unsetenv("API_KEY3")
		}()

		cfg, err := configs.LoadConfig()
		assert.NoError(t, err)
		assert.Equal(t, "test_user", cfg.DBUser)
		assert.Equal(t, "test_password", cfg.DBPassword)
		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "3306", cfg.DBPort)
		assert.Equal(t, "test_db", cfg.DBName)
		assert.Equal(t, "127.0.0.1", cfg.ServerAddress)
		assert.Equal(t, "8080", cfg.ServerPort)
		assert.Equal(t, "key1", cfg.APIKey1)
		assert.Equal(t, "key2", cfg.APIKey2)
		assert.Equal(t, "key3", cfg.APIKey3)
	})

	t.Run("異常系 - 必要な設定が見つからない場合にエラーを返す", func(t *testing.T) {
		os.Clearenv()

		_, err := configs.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "必要な設定が見つかりません")
	})
}
