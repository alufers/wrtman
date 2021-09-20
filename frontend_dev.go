// +build !embed_frontend 

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func MountFrontend(app *fiber.App) {
	// app.Use("/", filesystem.New(filesystem.Config{
	// 	Root:   http.FS(embeddedFrontend),
	// 	Browse: true,

	// }))
	app.Use(func(c *fiber.Ctx) error {
		if err := proxy.Do(c, "http://localhost:5000"+c.OriginalURL()); err != nil {
			return err
		}
		return nil
	})
}
