package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oke11o/go-telegram-bot/internal/app"
)

var (
	Version = "dev"
)

func main() {
	err := app.Run(context.Background(), Version)
	if err != nil {
		fmt.Printf("\nSTOP with error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}
