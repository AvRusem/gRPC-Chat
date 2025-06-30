package services

import (
	"context"

	"cu.ru/api/pb"
	"cu.ru/internal/client"
)

type ChatService struct {
	client   *client.ChatClient
	messages chan *pb.ChatMessage
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewChatService(client *client.ChatClient) *ChatService {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChatService{
		client:   client,
		messages: make(chan *pb.ChatMessage),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *ChatService) StartReceiving() error {
	err := s.client.StartStream(s.ctx)
	if err != nil {
		return err
	}
	go s.client.ReceiveMessages(s.messages)
	return nil
}

func (s *ChatService) Send(text string) error {
	return s.client.SendMessage(text)
}

func (s *ChatService) Messages() <-chan *pb.ChatMessage {
	return s.messages
}

func (s *ChatService) BanUser(login string) error {
	return s.client.BanUser(s.ctx, login)
}
