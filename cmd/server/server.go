package main

import (
	"fmt"
	"os"

	"cu.ru/internal/chat/servers"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: server (required <host:port>) (required <path/to/the/file/withBadWords>) (optional <authStorage>)")
		os.Exit(1)
	}

	addr := os.Args[1]
	badWords := os.Args[2]
	authStorage := ""
	disabled := true
	if len(os.Args) >= 4 {
		authStorage = os.Args[3]
		disabled = false
	}

	servers.StartChatServer(addr, badWords, authStorage, disabled)
}
