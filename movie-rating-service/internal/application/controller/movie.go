package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"log/slog"
	"movie-rating-service/config"
	"movie-rating-service/internal/application/middleware"
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
	app.Put("/movie/:id", authMiddleware.AdminHandler, controller.UpdateMovie)
	app.Delete("/movie/:id", authMiddleware.AdminHandler, controller.DeleteMovie)
	app.Get("/movie/:id", controller.GetMovie)

}

// @Summary GetByID Movie
// @Tags Movie
// @Param id path string true "Movie Id"
// @Success 200 {object} response.SuccessResponse{data=response.GetMovie}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /movie/{id} [get]
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

// @Summary Update Movie
// @Tags Movie
// @Param id   path     int  true  "Movie ID"
// @Param body body request.UpdateMovie true "Movie update payload"
// @Success 200 {object} response.SuccessResponse
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /movie/{id} [put]
func (c *movieController) UpdateMovie(ctx *fiber.Ctx) error {
	var req request.UpdateMovie
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	id := ctx.Params("id")
	req.ID = cast.ToUint(id)

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	err = c.movieService.Update(ctx.UserContext(), req)
	if err != nil {
		slog.Info("Movie could not updated")
		return err
	}

	slog.Info("Movie updated")
	return ctx.Status(fiber.StatusCreated).JSON(response.Success(struct{}{}))
}

// @Summary Delete Movie
// @Tags Movie
// @Param id   path      int  true  "Movie ID"
// @Success 200 {object} response.SuccessResponse
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /movie/{id} [delete]
func (c *movieController) DeleteMovie(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	req := request.DeleteMovie{ID: cast.ToUint(id)}

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	err = c.movieService.Delete(ctx.UserContext(), req)
	if err != nil {
		slog.Info("Movie could not deleted")
		return err
	}

	slog.Info("Movie deleted")
	return ctx.Status(fiber.StatusCreated).JSON(response.Success(struct{}{}))
}

// @Summary Create Movie
// @Tags Movie
// @Param body body request.CreateMovie true "Movie create payload"
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
