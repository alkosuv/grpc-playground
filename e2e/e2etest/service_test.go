//go:build e2etest

package e2etest

import (
	"context"
	"net"
	"testing"

	apipb "github.com/alkosuv/grpc-playgroud/e2e/api"
	"github.com/alkosuv/grpc-playgroud/e2e/service"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// Размер буфер должен быть, чем самое большое передаваемое сообщение
const bufSize = 1024 * 1024

func serverAndClientSetup(t *testing.T) (apipb.UserServiceClient, func()) {
	// 1. Создаем in-memory слушатель
	// Заменяет сетевое соединение
	listener := bufconn.Listen(bufSize)

	// 2. Инициализируем и регистрируем реальный gRPC сервер
	grpcServer := grpc.NewServer(
	// Здесь можно подключить ваши реальные интерцепторы (Auth, Logging и т.д.)
	)

	// Реальный инстанс бизнес-логики (можно передать тестовую БД)
	userServer := service.New()
	apipb.RegisterUserServiceServer(grpcServer, userServer)

	// 3. Запускаем сервер в отдельной горутине
	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("Server exited with error: %v", err)
		}
	}()

	// 4. Настраиваем кастомный Dial для клиента, перенаправляющий трафик в bufconn
	conn, err := grpc.NewClient("passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	// 5. Создаем клиент
	client := apipb.NewUserServiceClient(conn)

	// Возвращаем клиент и функцию очистки ресурсов (teardown)
	cleanup := func() {
		conn.Close()
		grpcServer.GracefulStop()
	}

	return client, cleanup
}

func TestE2E_CreateUser(t *testing.T) {
	client, cleanup := serverAndClientSetup(t)
	defer cleanup()

	t.Run("CreateUser – валидный запрос", func(t *testing.T) {
		resp, err := client.CreateUser(context.Background(),
			&apipb.CreateUserRequest{
				Name:  "Ivan",
				Login: "Ivan123",
				Email: "ivan123@gmail.com",
				Age:   20,
			},
		)
		require.NoError(t, err)
		require.NotZero(t, resp.Uuid)
	})

	t.Run("CreateUser – Invalid data name", func(t *testing.T) {
		resp, err := client.CreateUser(context.Background(),
			&apipb.CreateUserRequest{
				Name:  "    ",
				Login: "Ivan123",
				Email: "ivan123@gmail.com",
				Age:   20,
			},
		)
		require.Error(t, err)
		require.Zero(t, resp)

		st, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Equal(t, "Invalid data name.", st.Message())
	})

	t.Run("CreateUser – Invalid data login", func(t *testing.T) {
		resp, err := client.CreateUser(context.Background(),
			&apipb.CreateUserRequest{
				Name:  "Ivan",
				Login: "123Ivan123",
				Email: "ivan123@gmail.com",
				Age:   20,
			},
		)
		require.Error(t, err)
		require.Zero(t, resp)

		st, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Equal(t, "Invalid data login.", st.Message())
	})

	t.Run("CreateUser – Invalid data email", func(t *testing.T) {
		resp, err := client.CreateUser(context.Background(),
			&apipb.CreateUserRequest{
				Name:  "Ivan",
				Login: "Ivan123",
				Email: "ivan123",
				Age:   20,
			},
		)
		require.Error(t, err)
		require.Zero(t, resp)

		st, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Equal(t, "Invalid data email.", st.Message())
	})

	t.Run("CreateUser – Invalid data age", func(t *testing.T) {
		resp, err := client.CreateUser(context.Background(),
			&apipb.CreateUserRequest{
				Name:  "Ivan",
				Login: "Ivan123",
				Email: "ivan@gmail.com",
				Age:   121,
			},
		)
		require.Error(t, err)
		require.Zero(t, resp)

		st, _ := status.FromError(err)
		require.Equal(t, codes.InvalidArgument, st.Code())
		require.Equal(t, "Invalid data age. Age must be an integer between 18 and 120.", st.Message())
	})
}
