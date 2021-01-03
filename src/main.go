package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	listenPort = os.Getenv("LISTEN_PORT")
	db         *gorm.DB
)

func init() {
	if listenPort == "" {
		log.Fatal("LISTEN_PORT env var not defined.")
	}
	var err error
	db, err = gorm.Open(sqlite.Open("rust-bans.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()

	log.Println("Listening on " + listenPort)
	log.Fatal(app.Listen(":"))
}
