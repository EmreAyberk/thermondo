package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"log/slog"
	"movie-rating-service/config"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/application/service"
	validate "movie-rating-service/internal/application/validator"
)

type movieController struct {
	movieService service.MovieService
}

func NewMovieController(app *fiber.App, movieService service.MovieService) {
	authMiddleware := middleware.NewAuthMiddleware(config.Cfg.JWTSecret)

	controller := &movieController{movieService: movieService}

	app.Post("/movie", authMiddleware.AdminHandler, controller.CreateMovie)
	app.Get("/movie/:id", controller.GetMovie)

}

// @Summary GetByID Movie
// @Tags Movie
// @Param id path string true "Movie Id"
// @Success 200 {object} response.SuccessResponse{data=response.GetMovie}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Router /movie/:id [get]
func (c *movieController) GetMovie(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	req := request.GetMovie{ID: cast.ToUint(id)}
	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	res, err := c.movieService.Get(ctx.UserContext(), req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success(res))
}

// @Summary Create Movie
// @Tags Movie
// @Success 200 {object} response.SuccessResponse{data=response.CreateMovie}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Router /movie [post]
func (c *movieController) CreateMovie(ctx *fiber.Ctx) error {
	var req request.CreateMovie
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	res, err := c.movieService.Create(ctx.UserContext(), req)
	if err != nil {
		slog.Info("Movie could not create")
		return err
	}

	slog.Info("Movie created", "movie_id", res.ID)
	return ctx.Status(fiber.StatusCreated).JSON(response.Success(res))
}
