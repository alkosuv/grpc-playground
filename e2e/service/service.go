package service

import (
	"context"
	"regexp"
	"strings"

	apipb "github.com/alkosuv/grpc-playgroud/e2e/api"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	loginRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`)
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type apiService struct {
	apipb.UnimplementedUserServiceServer
}

func New() *apiService {
	return new(apiService)
}

func (_ *apiService) CreateUser(_ context.Context, in *apipb.CreateUserRequest) (*apipb.CreateUserResponse, error) {
	if strings.TrimSpace(in.Name) == "" {
		return nil, status.New(codes.InvalidArgument, "Invalid data name.").Err()
	}

	if !loginRegex.MatchString(in.GetLogin()) {
		return nil, status.New(codes.InvalidArgument, "Invalid data login.").Err()
	}

	if !emailRegex.MatchString(in.GetEmail()) {
		return nil, status.New(codes.InvalidArgument, "Invalid data email.").Err()
	}

	if in.GetAge() < 18 || in.GetAge() > 120 {
		return nil, status.New(codes.InvalidArgument, "Invalid data age. Age must be an integer between 18 and 120.").Err()
	}

	return &apipb.CreateUserResponse{Uuid: uuid.NewString()}, nil
}
