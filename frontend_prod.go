// +build embed_frontend

package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed wrtman-frontend/public/*
var embeddedFrontend embed.FS

func MountFrontend(app *fiber.App) {
	subFS, err := fs.Sub(embeddedFrontend, "wrtman-frontend/public")
	if err != nil {
		log.Fatalf("failed to generate subFS: %v", err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(subFS),
		Browse: true,
	}))
	log.Printf("Production frontend mounted")

}
