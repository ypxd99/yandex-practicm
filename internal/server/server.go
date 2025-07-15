package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ypxd99/yandex-practicm/util"
)

// Server представляет HTTP-сервер для сервиса сокращения URL.
// Обертывает стандартный http.Server с дополнительной конфигурацией
// и предоставляет методы для запуска и остановки сервера.
type Server struct {
	httpServer *http.Server
}

// Run запускает HTTP-сервер и начинает прослушивание запросов.
// Возвращает ошибку, если сервер не удалось запустить.
func (s *Server) Run() error {
	cfg := util.GetConfig().Server
	if cfg.EnableHTTPS {
		if cfg.TLSCertPath == "" || cfg.TLSKeyPath == "" {
			return errors.New("TLS enabled but cert or key path is empty")
		}
		return s.httpServer.ListenAndServeTLS(cfg.TLSCertPath, cfg.TLSKeyPath)
	}
	return s.httpServer.ListenAndServe()
}

// Stop корректно останавливает сервер с учетом заданного контекста.
// Ожидает завершения существующих соединений перед остановкой.
// Возвращает ошибку, если остановка не удалась.
func (s *Server) Stop(ctx context.Context) error {
	util.GetLogger().Infof("shutting down server...")
	return s.httpServer.Shutdown(ctx)
}

// NewServer создает и возвращает новый экземпляр Server с предоставленным обработчиком.
// Настраивает сервер с параметрами из конфигурации приложения,
// включая адрес, таймаут чтения и таймаут записи.
func NewServer(handler http.Handler) *Server {
	cfg := util.GetConfig().Server
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.ServerAddress,
			ReadTimeout:  time.Duration(cfg.RTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WTimeout) * time.Second,
			Handler:      handler,
		},
	}
}
