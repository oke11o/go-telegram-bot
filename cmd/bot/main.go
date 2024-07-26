package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/oke11o/go-telegram-bot/internal/app"
)

var (
	Version = "dev"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	err := app.Run(ctx, Version)
	if err != nil {
		fmt.Printf("\nSTOP with error: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("DONE")
}
