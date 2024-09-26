package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	config "github.com/shinya-ac/server1Q1A/configs"
	"github.com/shinya-ac/server1Q1A/middlewares/auth0"
	"github.com/shinya-ac/server1Q1A/pkg/logging"
	"github.com/shinya-ac/server1Q1A/presentation/settings"
	"github.com/shinya-ac/server1Q1A/server/route"
)

func Run(ctx context.Context, conf config.ConfigList) {
	api := settings.NewEchoEngine()

	jwks, err := auth0.FetchJWKS(conf.AuthDomain)
	if err != nil {
		logging.Logger.Error("JWTキー取得エラー", ":", err)
	}

	jwtMiddleware, err := auth0.NewMiddleware(conf.AuthDomain, conf.AuthClientID, jwks)
	if err != nil {
		logging.Logger.Error("jwtミドルウェア初期化エラー", ":", err)
	}

	api.Use(echo.WrapMiddleware(auth0.WithJWTMiddleware(jwtMiddleware)))

	route.InitRoute(api)

	address := conf.ServerAddress + ":" + conf.ServerPort
	logging.Logger.Error("サーバー起動中...", "address:", address)
	srv := &http.Server{
		Addr:              address,
		Handler:           api,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Minute,
		WriteTimeout:      10 * time.Minute,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Error("サーバーの配信に失敗", "error", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.Logger.Error("サーバー停止", "error", err)
		os.Exit(1)
	}
}
