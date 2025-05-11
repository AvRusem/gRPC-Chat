package main

import (
	"fmt"
	"os"

	"cu.ru/internal/client"
	"cu.ru/internal/client/services"
	"cu.ru/internal/client/ui"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: client (required <host:port>) (required <username>) (optional <password>)")
		os.Exit(1)
	}

	addr := os.Args[1]
	username := os.Args[2]
	password := ""
	if len(os.Args) >= 4 {
		password = os.Args[3]
	}

	authClient := client.NewAuthClient(addr)
	defer authClient.Close()
	token := authClient.Auth(username, password)
	fmt.Println("Token:", token)

	chatClient := client.NewChatClient(addr, username, token)
	defer chatClient.Close()
	chatService := services.NewChatService(chatClient)
	err := ui.StartUI(chatService, username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
