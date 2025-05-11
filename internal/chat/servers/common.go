package servers

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"cu.ru/internal/chat/interceptors"
	"google.golang.org/grpc"
)

func buildServer(addr, badWords, authStorage string, authDisabled bool) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.LoggingInterceptor,
			interceptors.AuthUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			interceptors.AuditStreamInterceptor,
			interceptors.AuthStreamInterceptor,
		),
	)
	buildAuthServer(authDisabled, authStorage, grpcServer)
	buildChatServer(badWords, grpcServer)

	return grpcServer, listener
}

func runChat(server *grpc.Server, listener net.Listener) {
	log.Printf("gRPC-server running on %v...", listener.Addr().String())

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func StartChatServer(addr, badWords, authStorage string, authDisabled bool) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	server, listener := buildServer(addr, badWords, authStorage, authDisabled)

	go runChat(server, listener)

	<-stop
	log.Println("Server is shutting down...")
	server.Stop()
	log.Println("Server gracefully stopped")
}
