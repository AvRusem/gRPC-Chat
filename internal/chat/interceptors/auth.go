package interceptors

import (
	"context"
	"strings"

	"cu.ru/internal/chat/tokens"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	ClientIDKey contextKey = "clientID"
	RoleKey     contextKey = "role"
)

// I have no idea how else to pass the context to the stream
// because the grpc.ServerStream interface doesn't have a way to set the context
// and the context is not passed to the handler
type WrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *WrappedServerStream) Context() context.Context {
	return w.ctx
}

func AuthUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if info.FullMethod == "/auth.AuthService/Auth" {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no auth metadata")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	token := values[0]
	if len(token) > 7 && strings.ToLower(token[:7]) == "bearer " {
		token = token[7:]
	}

	claims, err := tokens.ValidateToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	if info.FullMethod == "/chat.ChatService/BanUser" && claims.Role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "only admins can ban users")
	}

	newCtx := context.WithValue(ctx, ClientIDKey, claims.UserID)
	newCtx = context.WithValue(newCtx, RoleKey, claims.Role)
	return handler(newCtx, req)
}

func AuthStreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	if info.FullMethod == "/auth.AuthService/Auth" {
		return handler(srv, ss)
	}

	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "no auth metadata")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "missing token")
	}

	token := values[0]
	if len(token) > 7 && strings.ToLower(token[:7]) == "bearer " {
		token = token[7:]
	}

	claims, err := tokens.ValidateToken(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	if info.FullMethod == "/chat.ChatService/BanUser" && claims.Role != "admin" {
		return status.Error(codes.PermissionDenied, "only admins can ban users")
	}

	// idk mb you know a better way to do this
	newCtx := context.WithValue(ctx, ClientIDKey, claims.UserID)
	newCtx = context.WithValue(newCtx, RoleKey, claims.Role)
	wrappedSS := &WrappedServerStream{
		ServerStream: ss,
		ctx:          newCtx,
	}
	return handler(srv, wrappedSS)
}
