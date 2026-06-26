package main

import (
	"log/slog"
	"net"
	"os"

	apipb "github.com/alkosuv/grpc-playgroud/e2e/api"
	"github.com/alkosuv/grpc-playgroud/e2e/service"
	"google.golang.org/grpc"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	// 1. Открываем TCP-порт для прослушивания
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("Не удалось открыть порт", "error", err)
		os.Exit(0)
	}

	// 2. Создаем экземпляр gRPC сервера
	grpcServer := grpc.NewServer()

	// 3. Регистрируем наш сервис на этом сервере
	apipb.RegisterUserServiceServer(grpcServer, service.New())

	slog.Info("gRPC сервер запущен на порту :50051...")

	// 4. Запускаем сервер
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Ошибка запуска сервера", "error", err)
		os.Exit(0)
	}
}
