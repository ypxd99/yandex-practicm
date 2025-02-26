package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/repository/postgres"
	"github.com/ypxd99/yandex-practicm/internal/server"
	"github.com/ypxd99/yandex-practicm/internal/transport/handler"
	"github.com/ypxd99/yandex-practicm/util"
)

func main() {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	//go util.GenerateRSA()
	logger := util.GetLogger()
	logger.Info("start shortener service")

	if cfg.Postgres.MakeMigration {
		go makeMegrations()
	}

	router := gin.Default()
	handler.InitRoutes(router)

	srv := server.NewServer(router)
	go func() {
		util.GetLogger().Infof("SHORTENER server listeing at: %s:%d",
			cfg.Server.Address,
			cfg.Server.Port)

		err := srv.Run()
		if !errors.Is(err, http.ErrServerClosed) {
			util.GetLogger().Fatalf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		util.GetLogger().Fatalf("Server forced to shutdown: %s", err.Error())
	}
	util.GetLogger().Log(4, "HTTP SHORTENER service stopped")
}

func makeMegrations() {
	// migrate UP
	util.GetLogger().Info("start migrations")
	err := postgres.MigrateDBUp(context.Background())
	if err != nil {
		util.GetLogger().Error(err)
		return
	}
	util.GetLogger().Info("migrations up")

	// migrate DOWN
	//err = postgres.MigrateDBDown(context.Background())
	//if err != nil {
	//	util.GetLogger().Error(err)
	//	return
	//}
}
