package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cu.ru/internal/client/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StartUI(service *services.ChatService, username string) error {
	err := service.StartReceiving()
	if err != nil {
		return fmt.Errorf("failed to connect to the chat: %w", err)
	}

	go func() {
		for msg := range service.Messages() {
			if msg.GetIsError() {
				switch msg.GetText() {
				case "banned":
					fmt.Println("ОШИБКА: Ты больше не можешь писать в этот чат")
				case "censored":
					fmt.Println("ОШИБКА: Нельзя ругаться")
				default:
					fmt.Printf("ОШИБКА: %s\n", msg.Text)
				}
				continue
			}
			fmt.Printf("[%s]: %s\n", msg.Login, msg.Text)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if !scanner.Scan() {
			break
		}
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			continue
		}

		if strings.HasPrefix(text, "/") {
			switch {
			case text == "/exit":
				fmt.Println("Выход из чата...")
				return nil
			case strings.HasPrefix(text, "/ban "):
				target := strings.TrimSpace(strings.TrimPrefix(text, "/ban"))
				if err := service.BanUser(target); err != nil {
					st, _ := status.FromError(err)
					if st.Code() == codes.PermissionDenied {
						fmt.Println("ОШИБКА: Ты не администратор")
					} else {
						fmt.Printf("ОШИБКА: Не удалось забанить %s: %v\n", target, err)
					}
				}
			default:
				fmt.Println("ОШИБКА: Неизвестная команда")
			}
			continue
		}

		if err := service.Send(text); err != nil {
			fmt.Printf("ОШИБКА: Не удалось отправить сообщение: %v\n", err)
		}
	}

	return scanner.Err()
}
