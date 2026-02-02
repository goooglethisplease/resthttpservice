package main

import (
	"log"

	"github.com/joho/godotenv"

	"restservice/cmd/app/run"
)

func main() {
	_ = godotenv.Load()

	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
