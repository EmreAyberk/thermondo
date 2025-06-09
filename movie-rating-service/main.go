package main

import (
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
	"movie-rating-service/internal/infrastructure/repository"
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

	err := config.Init()
	if err != nil {
		panic(err)
	}

	database, err := db.Connect()
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: common.ErrorHandler()})
	app.Get("/swagger/*", swagger.HandlerDefault)

	prometheus := fiberprometheus.New("movie-rating-service")
	prometheus.RegisterAt(app, "/metrics")

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

	err = app.Listen(fmt.Sprintf(":%d", config.Cfg.Port))
	if err != nil {
		panic(err)
	}

}
