package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/proto/shortener"
	"github.com/ypxd99/yandex-practicm/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer представляет gRPC-сервер для сервиса сокращения URL.
type GRPCServer struct {
	server *grpc.Server
	port   uint
}

// NewGRPCServer создает и возвращает новый экземпляр GRPCServer.
func NewGRPCServer(service service.LinkService) *GRPCServer {
	// Создаем gRPC сервер с опциями
	grpcServer := grpc.NewServer()
	
	// Создаем обработчик
	handler := NewGRPCHandler(service)
	
	// Регистрируем сервис
	shortener.RegisterShortenerServiceServer(grpcServer, handler)
	
	// Включаем reflection для отладки
	reflection.Register(grpcServer)
	
	return &GRPCServer{
		server: grpcServer,
		port:   util.GetConfig().Server.GRPCPort,
	}
}

// Run запускает gRPC-сервер и начинает прослушивание запросов.
// Возвращает ошибку, если сервер не удалось запустить.
func (s *GRPCServer) Run() error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	
	util.GetLogger().Infof("gRPC server starting on port %d", s.port)
	return s.server.Serve(listener)
}

// Stop корректно останавливает gRPC-сервер.
// Ожидает завершения существующих соединений перед остановкой.
func (s *GRPCServer) Stop(ctx context.Context) error {
	util.GetLogger().Infof("shutting down gRPC server...")
	s.server.GracefulStop()
	return nil
}
