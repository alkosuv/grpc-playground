package main

import (
	"context"
	"log/slog"
	"os"

	pb "github.com/alkosuv/grpc-playgroud/error/api"
	errorpb "github.com/alkosuv/grpc-playgroud/error/pkg/errors/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("did not connect", "error", err)
		return
	}
	defer conn.Close()

	client := pb.NewErrorTestServiceClient(conn)
	response, err := client.GetError(context.Background(),
		// &pb.ErrorRequest{ErrorType: "BadRequest"},
		&pb.ErrorRequest{ErrorType: "InvalidArgument"},
	)
	slog.Debug("GetError request", "response", response, "error", err)
	if err != nil {
		st, ok := status.FromError(err)
		slog.Debug("status FromError", "status", st, "ok", ok)

		if ok {
			slog.Debug("st.Details", "len", len(st.Details()))
			for _, detail := range st.Details() {
				switch errorType := detail.(type) {
				case *errorpb.BadRequest:
					slog.Warn("BadRequest", "message", errorType.Message, "ditales", errorType.Ditales)
				case *errorpb.InvalidArgument:
					slog.Warn("InvalidArgument", "message", errorType.Message, "ditales", errorType.Ditales)
				default:
					slog.Error("Unexpected", "error", errorType)
				}
			}
		} else {
			slog.Error("not castom error")
			return
		}

		return
	}

	slog.Info("GetError", "response", response)
}
