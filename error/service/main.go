package main

import (
	"context"
	"log/slog"
	"net"
	"os"

	pb "github.com/alkosuv/grpc-playgroud/error/api"
	errorpb "github.com/alkosuv/grpc-playgroud/error/pkg/errors/grpc"

	// "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type apiService struct {
	pb.UnimplementedErrorTestServiceServer
}

func (_ *apiService) GetError(ctx context.Context, in *pb.ErrorRequest) (*pb.ErrorResponse, error) {
	slog.Debug("GetError")

	switch in.ErrorType {
	case "BadRequest":
		badRequest := &errorpb.BadRequest{
			Message: "BadRequest",
		}
		st := status.New(codes.InvalidArgument, "bad request")
		st, err := st.WithDetails(badRequest)
		if err != nil {
			slog.Error("invali with details")
			return nil, status.New(codes.Internal, "error with details").Err()
		}
		return nil, st.Err()
	case "InvalidArgument":
		invalidArgument := &errorpb.InvalidArgument{
			Message: "InvalidArgument",
			Ditales: []*errorpb.Ditale{
				{
					Desciption: "desciption Invalid Argument",
					Field:      "any field",
					Value:      "any value",
				},
			},
		}
		st := status.New(codes.InvalidArgument, "invalid parameter")
		st, err := st.WithDetails(invalidArgument)
		if err != nil {
			slog.Error("invali with details")
			return nil, status.New(codes.Internal, "error with details").Err()
		}
		return nil, st.Err()
	default:
		return nil, status.New(codes.Internal, "error ErrorType").Err()
	}
}

func main() {
	// 1. Открываем TCP-порт для прослушивания
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		slog.Error("Не удалось открыть порт", "error", err)
		os.Exit(0)
	}

	// 2. Создаем экземпляр gRPC сервера
	grpcServer := grpc.NewServer()

	// 3. Регистрируем наш сервис на этом сервере
	pb.RegisterErrorTestServiceServer(grpcServer, &apiService{})

	slog.Info("gRPC сервер запущен на порту :50051...")

	// 4. Запускаем сервер
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Ошибка запуска сервера", "error", err)
		os.Exit(0)
	}
}
