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

	app.Post("/rating", controller.CreateRating)
	app.Get("/rating/user", authMiddleware.UserHandler, controller.GetUserRatings)
}

// @Summary Create Rating
// @Tags Rating
// @Success 200 {object} response.SuccessResponse{data=response.CreateRating}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Router /rating [post]
func (c *ratingController) CreateRating(ctx *fiber.Ctx) error {
	var req request.CreateRating
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
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

// @Summary GetUserRatings User
// @Tags Rating
// @Success 200 {object} response.SuccessResponse{data=response.GetUserRatings}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
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
