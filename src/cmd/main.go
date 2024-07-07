package main

import (
	"context"
	"fmt"

	config "github.com/shinya-ac/server1Q1A/configs"
	"github.com/shinya-ac/server1Q1A/infrastructure/mysql/db"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
	"github.com/shinya-ac/server1Q1A/server"
)

func main() {
	logging.InitLogger()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logging.Logger.Info("log初期化")

	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Error("configの読み込みに失敗", "error", err)
	}
	fmt.Printf(cfg.APIKey1)

	db.NewMainDB(cfg)

	server.Run(ctx, cfg)
}
