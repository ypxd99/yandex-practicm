package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ypxd99/yandex-practicm/internal/repository"
	"github.com/ypxd99/yandex-practicm/internal/repository/postgres"
	"github.com/ypxd99/yandex-practicm/internal/repository/storage"
	"github.com/ypxd99/yandex-practicm/internal/server"
	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/internal/transport/grpc"
	"github.com/ypxd99/yandex-practicm/internal/transport/handler"
	"github.com/ypxd99/yandex-practicm/util"
)

// buildVersion содержит версию сборки, может быть переопределена через ldflags.
// buildDate содержит дату сборки, может быть переопределена через ldflags.
// buildCommit содержит хеш коммита, может быть переопределён через ldflags.
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	fmt.Printf("Build version: %s\n", ifNA(buildVersion))
	fmt.Printf("Build date: %s\n", ifNA(buildDate))
	fmt.Printf("Build commit: %s\n", ifNA(buildCommit))

	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	// go util.GenerateRSA()
	logger := util.GetLogger()
	logger.Info("start shortener service")

	if cfg.Postgres.MakeMigration && cfg.Postgres.UsePostgres {
		go makeMegrations()
	}

	var (
		repo repository.LinkRepository
		err  error
	)
	if cfg.Postgres.UsePostgres {
		repo, err = postgres.Connect(context.Background())
		if err != nil {
			logger.Errorf("Failed to initialize Postgres: %v", err)
			return
		}
	} else {
		repo, err = storage.InitStorage(cfg.FileStoragePath)
		if err != nil {
			logger.Errorf("Failed to initialize Storage: %v", err)
			return
		}
	}
	defer repo.Close()

	service := service.InitService(repo)
	h := handler.InitHandler(service)

	router := gin.Default()
	h.InitRoutes(router)
	httpServer := server.NewServer(router)

	grpcServer := grpc.NewGRPCServer(service)

	go func() {
		util.GetLogger().Infof("HTTP server listening at: %s", cfg.Server.ServerAddress)
		if err := httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("error occurred while running HTTP server: %s\n", err.Error())
		}
	}()

	go func() {
		util.GetLogger().Infof("gRPC server listening on port: %d", cfg.Server.GRPCPort)
		if err := grpcServer.Run(); err != nil {
			logger.Errorf("error occurred while running gRPC server: %s\n", err.Error())
		}
	}()

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logger.Errorf("HTTP server forced to shutdown: %s", err.Error())
	}

	if err := grpcServer.Stop(ctx); err != nil {
		logger.Errorf("gRPC server forced to shutdown: %s", err.Error())
	}

	logger.Info("All servers stopped")
}

func makeMegrations() {
	// migrate UP
	util.GetLogger().Info("start migrations")
	if err := postgres.MigrateDBUp(context.Background()); err != nil {
		util.GetLogger().Error(err)
		return
	}
	util.GetLogger().Info("migrations up")

	// migrate DOWN
	// if err := postgres.MigrateDBDown(context.Background()); err != nil {
	//	util.GetLogger().Error(err)
	//	return
	// }
}

// ifNA возвращает значение или "N/A", если оно пустое.
func ifNA(val string) string {
	if val == "" {
		return "N/A"
	}
	return val
}
