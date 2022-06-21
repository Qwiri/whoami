package main

import (
	"github.com/Qwiri/gobby/pkg/gobby"
	"github.com/Qwiri/gobby/pkg/handlers"
	"github.com/apex/log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	g := gobby.New(app)

	g.HandleAll(gobby.Handlers{
		"CHAT": handlers.Chat,
		"LIST": handlers.List,
	})

	if err := app.Listen(":8081"); err != nil {
		log.WithError(err).Warn("cannot serve")
	}
}
