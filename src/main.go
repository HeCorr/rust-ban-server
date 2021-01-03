package main

import (
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const version string = "v1.0.0"

var (
	listenAddr string
	quiet      bool
	db         *gorm.DB
)

func init() {
	log.Println("Rust Ban Server " + version)
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
	flag.StringVar(&listenAddr, "l", ":4000", "Listen address (default: ':4000')")
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
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "SteamID not found."})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
			if errors.Is(err, errNotInserted) {
				return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "SteamID already banned."})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "SteamID banned."})
	})

	app.Delete("/api/rustBans/:steamID64", func(c *fiber.Ctx) error {
		err := delBan(c.Params("steamID64"))
		if err != nil {
			if errors.Is(err, errNotDeleted) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "SteamID not banned."})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "SteamID unbanned."})
	})

	log.Println("Listening on " + listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
