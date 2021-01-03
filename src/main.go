package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	listenAddr = os.Getenv("LISTEN_ADDR")
	db         *gorm.DB
)

func init() {
	if listenAddr == "" {
		log.Fatal("LISTEN_ADDR env var not defined.")
	}
	var err error
	db, err = gorm.Open(sqlite.Open("rust-bans.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&Ban{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	log.Println("Listening on " + listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
