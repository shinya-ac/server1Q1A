package main

import (
	"fmt"

	config "github.com/shinya-ac/server1Q1A/configs"
)

func main() {
	// logging.InitLogger()
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("configエラー：%+v", err)
		// logging.Logger.Error("configの読み込みに失敗", "error", err)
	}
	fmt.Printf(cfg.APIKey1)

	// db.NewMainDB(cfg)

	// server.Run(ctx, cfg)
}
