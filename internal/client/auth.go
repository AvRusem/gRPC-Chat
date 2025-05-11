package client

import (
	"context"
	"fmt"
	"time"

	"cu.ru/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

func NewAuthClient(addr string) *AuthClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("Can't connect to auth server: %v", err))
	}

	client := pb.NewAuthServiceClient(conn)
	return &AuthClient{conn: conn, client: client}
}

func (a *AuthClient) Auth(login, password string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.AuthMessage{
		Login:    login,
		Password: password,
	}

	resp, err := a.client.Auth(ctx, req)
	if err != nil {
		fmt.Printf("Auth Error: %v\n", err)
		panic("failed to log in")
	}

	return resp.Token
}

func (a *AuthClient) Close() {
	_ = a.conn.Close()
}
