package servers

import (
	"context"
	"io"
	"log"

	"cu.ru/api/pb"
	appErrors "cu.ru/internal/chat/errors"
	"cu.ru/internal/chat/interceptors"
	"cu.ru/internal/chat/repositories"
	"cu.ru/internal/chat/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	chatService services.ChatService
}

func NewChatServer(chatService services.ChatService) *ChatServer {
	return &ChatServer{
		chatService: chatService,
	}
}

func (c *ChatServer) StartChat(stream pb.ChatService_StartChatServer) error {
	ctx := stream.Context()
	clientID := ctx.Value(interceptors.ClientIDKey).(string)
	if clientID == "" {
		return status.Error(codes.Unauthenticated, "client ID not found in context")
	}

	if err := c.chatService.RegisterClient(clientID, stream); err != nil {
		if err == appErrors.AlreadyExistsError {
			return status.Error(codes.AlreadyExists, "client already exists")
		}
		return status.Error(codes.Internal, "failed to register client")
	}
	defer func() {
		if err := c.chatService.UnregisterClient(clientID); err != nil {
			log.Printf("failed to unregister client %s: %v", clientID, err)
		}
	}()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "error in the stream: %v", err)
		}

		if msg == nil {
			return status.Error(codes.InvalidArgument, "received nil message")
		}

		censored, err := c.chatService.ModerateMessage(clientID, msg.GetText())
		if err != nil {
			if err == appErrors.BannedError {
				if err := stream.Send(&pb.ChatMessage{
					Text:    "banned",
					IsError: true,
				}); err != nil {
					log.Printf("failed to send banned message: %v", err)
					return status.Error(codes.Internal, "failed to send banned message")
				}
				continue
			}
			log.Printf("failed to moderate message: %v", err)
			return status.Error(codes.Internal, "failed to moderate message")
		}
		if censored {
			if err := stream.Send(&pb.ChatMessage{
				Text:    "censored",
				IsError: true,
			}); err != nil {
				log.Printf("failed to send censored message: %v", err)
				return status.Error(codes.Internal, "failed to send censored message")
			}
			continue
		}

		if err := c.chatService.BroadcastMessage(clientID, msg.GetText()); err != nil {
			if err == appErrors.NotFoundError {
				return status.Error(codes.NotFound, "client not found")
			}
			if err == appErrors.BannedError {
				return status.Error(codes.PermissionDenied, "client is banned")
			}
			return status.Error(codes.Internal, "failed to broadcast message")
		}
	}
}

func (c *ChatServer) BanUser(ctx context.Context, banMessage *pb.BanMessage) (*pb.Empty, error) {
	clientID := ctx.Value(interceptors.ClientIDKey).(string)
	if clientID == "" {
		return nil, status.Error(codes.Unauthenticated, "client ID not found in context")
	}

	if err := c.chatService.BanClient(banMessage.GetTargetLogin()); err != nil {
		if err == appErrors.NotFoundError {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "failed to ban user")
	}

	return &pb.Empty{}, nil
}

func getChatService(file string) services.ChatService {
	profanityRepo := repositories.NewProfanityRepositoryInMemory(file)
	chatRepo := repositories.NewChatRepositoryInMemory()

	return services.NewChatServiceDefault(profanityRepo, chatRepo)
}

func buildChatServer(badWords string, grpcServer *grpc.Server) {
	pb.RegisterChatServiceServer(grpcServer, NewChatServer(getChatService(badWords)))
}
