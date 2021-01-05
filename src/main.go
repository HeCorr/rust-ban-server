package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const version string = "v1.1.0"

var (
	listenAddr string
	quiet      bool
	db         *gorm.DB
	sIDval     *regexp.Regexp
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
	sIDval, err = regexp.Compile(`^[0-9]{17}$`)
	if err != nil {
		log.Fatal(err)
	}
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
		if !sIDval.Match([]byte(c.Params("steamID64"))) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid SteamID64."})
		}
		ban, err := getBan(c.Params("steamID64"))
		if err != nil {
			if errors.Is(err, errNotFound) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "SteamID64 not found."})
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
		if !sIDval.Match([]byte(ban.SteamID)) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid SteamID64."})
		}
		err = addBan(ban)
		if err != nil {
			if errors.Is(err, errNotInserted) {
				return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "SteamID64 already banned."})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "SteamID64 banned."})
	})

	app.Delete("/api/rustBans/:steamID64", func(c *fiber.Ctx) error {
		if !sIDval.Match([]byte(c.Params("steamID64"))) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid SteamID64."})
		}
		err := delBan(c.Params("steamID64"))
		if err != nil {
			if errors.Is(err, errNotDeleted) {
				return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "SteamID64 not banned."})
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "SteamID64 unbanned."})
	})

	log.Println("Listening on " + listenAddr)
	log.Fatal(app.Listen(listenAddr))
}
