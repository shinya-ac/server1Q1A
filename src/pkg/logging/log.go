package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	Logger *slog.Logger
	once   sync.Once
)

func InitLogger() {
	once.Do(func() {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
}
