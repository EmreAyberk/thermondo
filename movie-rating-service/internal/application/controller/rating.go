package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"log/slog"
	"movie-rating-service/config"
	"movie-rating-service/internal/application/middleware"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/application/service"
	validate "movie-rating-service/internal/application/validator"
)

type ratingController struct {
	ratingService service.RatingService
}

func NewRatingController(app *fiber.App, ratingService service.RatingService) {
	authMiddleware := middleware.NewAuthMiddleware(config.Cfg.JWTSecret)

	controller := &ratingController{ratingService: ratingService}

	// Good improvement will be separating it from different groups, maybe separating it to movie and user.
	// It depends on the team choice

	app.Post("movie/:id/rating", authMiddleware.UserHandler, controller.CreateRating)
	app.Patch("movie/:id/rating", authMiddleware.UserHandler, controller.UpdateRating)
	app.Delete("movie/:id/rating", authMiddleware.UserHandler, controller.DeleteRating)
	app.Get("user/rating", authMiddleware.UserHandler, controller.GetUserRatings)
}

// @Summary Create Rating
// @Tags Rating
// @Param id path string true "Movie Id"
// @Param body body request.CreateRating true "Rating create payload"
// @Success 200 {object} response.SuccessResponse{data=response.CreateRating}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /movie/{id}/rating [post]
func (c *ratingController) CreateRating(ctx *fiber.Ctx) error {
	var req request.CreateRating
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	id := ctx.Params("id")
	req.MovieID = cast.ToUint(id)

	claims := ctx.Locals("user").(jwt.MapClaims)

	req.UserID = cast.ToUint(claims["user_id"])

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	res, err := c.ratingService.Create(ctx.UserContext(), req)
	if err != nil {
		slog.Info("Rating could not create")
		return err
	}

	slog.Info("Rating created", "user_id", res.ID)
	return ctx.Status(fiber.StatusCreated).JSON(response.Success(res))
}

// @Summary Update Rating
// @Tags Rating
// @Param id path string true "Movie Id"
// @Param body body request.UpdateRating true "Rating update payload"
// @Success 200 {object} response.SuccessResponse{data=response.UpdateRating}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /movie/{id}/rating [patch]
func (c *ratingController) UpdateRating(ctx *fiber.Ctx) error {
	var req request.UpdateRating
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	id := ctx.Params("id")
	req.MovieID = cast.ToUint(id)

	claims := ctx.Locals("user").(jwt.MapClaims)

	req.UserID = cast.ToUint(claims["user_id"])

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	res, err := c.ratingService.Update(ctx.UserContext(), req)
	if err != nil {
		slog.Info("Rating could not update")
		return err
	}

	slog.Info("Rating updated", "user_id", res.ID)
	return ctx.Status(fiber.StatusOK).JSON(response.Success(res))
}

// @Summary Delete Rating
// @Tags Rating
// @Param id path string true "Movie Id"
// @Success 200 {object} response.SuccessResponse
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /movie/{id}/rating [delete]
func (c *ratingController) DeleteRating(ctx *fiber.Ctx) error {
	var req request.DeleteRating

	id := ctx.Params("id")
	req.MovieID = cast.ToUint(id)

	claims := ctx.Locals("user").(jwt.MapClaims)

	req.UserID = cast.ToUint(claims["user_id"])

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	err = c.ratingService.Delete(ctx.UserContext(), req)
	if err != nil {
		slog.Info("Rating could not update")
		return err
	}

	slog.Info("Rating deleted")
	return ctx.Status(fiber.StatusOK).JSON(response.Success(struct{}{}))
}

// @Summary GetUserRatings User
// @Tags Rating
// @Success 200 {object} response.SuccessResponse{data=response.GetUserRatings}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Security BearerAuth
// @Router /rating/user [get]
func (c *ratingController) GetUserRatings(ctx *fiber.Ctx) error {
	claims := ctx.Locals("user").(jwt.MapClaims)

	req := request.GetUserRatings{UserID: cast.ToUint(claims["user_id"])}

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	res, err := c.ratingService.GetRatingsByUserID(ctx.UserContext(), req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success(res))
}
