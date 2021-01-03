package main

import (
	"errors"
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

	app.Get("/api/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/api/rustBans/:steamID64", func(c *fiber.Ctx) error {
		ban, err := getBan(c.Params("steamID64"))
		if err != nil {
			if errors.Is(err, errNotFound) {
				return c.SendStatus(http.StatusNotFound)
			}
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(ban)
	})

	log.Println("Listening on " + listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
