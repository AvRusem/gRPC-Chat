package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: server (required <authStorage>) (required <login>) (required <role>) (required <password>)")
		os.Exit(1)
	}

	authStorage := os.Args[1]
	login := os.Args[2]
	role := os.Args[3]
	password := os.Args[4]

	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}
	file, err := os.OpenFile(authStorage, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%s:%s:%s\n", login, hashedPassword, role))
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("User %s with role %s added successfully\n", login, role)
}
