package main

import (
	"context"
	"fmt"

	"github.com/oke11o/go-telegram-bot/internal/app"
)

func main() {
	err := app.Run(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("DONE")
}
