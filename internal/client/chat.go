package client

import (
	"context"
	"fmt"
	"io"

	"cu.ru/api/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type ChatClient struct {
	conn   *grpc.ClientConn
	client pb.ChatServiceClient
	login  string
	token  string
	stream pb.ChatService_StartChatClient
}

func NewChatClient(addr, login, token string) *ChatClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Sprintf("Can't connect to chat server: %v", err))
	}

	client := pb.NewChatServiceClient(conn)

	return &ChatClient{
		conn:   conn,
		client: client,
		login:  login,
		token:  token,
		stream: nil,
	}
}

func (c *ChatClient) StartStream(ctx context.Context) error {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + c.token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := c.client.StartChat(ctx)
	if err != nil {
		return fmt.Errorf("failed to start stream: %w", err)
	}
	c.stream = stream
	return nil
}

func (c *ChatClient) SendMessage(text string) error {
	msg := &pb.ChatMessage{
		Text: text,
	}
	return c.stream.Send(msg)
}

func (c *ChatClient) BanUser(ctx context.Context, targetedUser string) error {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + c.token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	msg := &pb.BanMessage{
		TargetLogin: targetedUser,
	}
	_, err := c.client.BanUser(ctx, msg)
	return err
}

func (c *ChatClient) ReceiveMessages(out chan<- *pb.ChatMessage) {
	for {
		msg, err := c.stream.Recv()
		if err == io.EOF {
			close(out)
			return
		}
		if err != nil {
			close(out)
			// fmt.Printf("error getting message: %v\n", err)
			return
		}
		out <- msg
	}
}

func (c *ChatClient) Close() {
	_ = c.conn.Close()
}
