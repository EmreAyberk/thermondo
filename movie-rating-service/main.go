package main

import (
	"context"
	"fmt"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/swagger"
	"log/slog"
	"movie-rating-service/config"
	_ "movie-rating-service/docs"
	"movie-rating-service/internal/application/controller"
	"movie-rating-service/internal/application/service"
	"movie-rating-service/internal/common"
	"movie-rating-service/internal/infrastructure/db"
	"movie-rating-service/internal/infrastructure/db/seeder"
	"movie-rating-service/internal/infrastructure/repository"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// @title           movieratingservice
// @version         1.0
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath  /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := config.Init()
	if err != nil {
		panic(err)
	}

	database, err := db.Connect()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 && strings.EqualFold(os.Args[1], "seeder") {
		s := seeder.NewSeeder(database)
		err = s.Seed()
		if err != nil {
			slog.Info("Seeder error", "error", err)
		}
		return
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: common.ErrorHandler()})
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/swagger/*", swagger.HandlerDefault)

	prometheus := fiberprometheus.New("movie-rating-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Get("/monitor", monitor.New())

	userRepository := repository.NewUserRepository(database)
	userService := service.NewUserService(userRepository)
	controller.NewUserController(app, userService)

	movieRepository := repository.NewMovieRepository(database)
	movieCacheRepository := repository.NewCachedMovieRepository(movieRepository, time.Second*30)

	movieService := service.NewMovieService(movieCacheRepository)
	controller.NewMovieController(app, movieService)

	ratingRepository := repository.NewRatingRepository(database)
	ratingService := service.NewRatingService(ratingRepository, movieRepository)
	controller.NewRatingController(app, ratingService)

	go func() {
		if err = app.Listen(fmt.Sprintf(":%d", config.Cfg.Port)); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("Server gracefully stopped")
	}
}
