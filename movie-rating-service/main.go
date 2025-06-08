package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"log/slog"
	"os"
	"time"
)

// @title           movieratingservice
// @version         1.0
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get("/monitor", monitor.New())

	err := app.Listen(":8080")
	if err != nil {
		panic(err)
	}

}
