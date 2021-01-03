package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	listenAddr = os.Getenv("LISTEN_ADDR")
	quiet      bool
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
	flag.BoolVar(&quiet, "q", false, "Quiet mode (don't print HTTP log)")
	flag.Parse()
}

func main() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	if !quiet {
		app.Use("/api", logger.New())
	}

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

	app.Post("/api/rustBans", func(c *fiber.Ctx) error {
		var ban Ban
		err := c.BodyParser(&ban)
		if err != nil {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}
		err = addBan(ban)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "created"})
	})

	log.Println("Listening on " + listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
