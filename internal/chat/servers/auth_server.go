package servers

import (
	"context"

	"cu.ru/api/pb"
	appErrors "cu.ru/internal/chat/errors"
	"cu.ru/internal/chat/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	authService services.AuthService
}

func NewAuthServer(authService services.AuthService) *AuthServer {
	return &AuthServer{
		authService: authService,
	}
}

func (a *AuthServer) Auth(ctx context.Context, authData *pb.AuthMessage) (*pb.AuthResponse, error) {
	token, err := a.authService.GenerateToken(authData.Login, authData.Password)
	if err != nil {
		if err == appErrors.ErrNotAuthorized {
			return nil, status.Error(codes.Unauthenticated, "not authorized")
		}
		return nil, err
	}

	return &pb.AuthResponse{
		Token: token,
	}, nil
}

func getAuthService(disabled bool, _ string) services.AuthService {
	if disabled {
		return services.NewAuthServiceDisabled()
	}
	panic("getService: not implemented yet")
}

func buildAuthServer(disabled bool, authStorage string, grpcServer *grpc.Server) {
	pb.RegisterAuthServiceServer(grpcServer, NewAuthServer(getAuthService(disabled, authStorage)))
}
