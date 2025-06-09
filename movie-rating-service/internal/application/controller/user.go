package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"movie-rating-service/config"
	"movie-rating-service/internal/application/middleware"
	"movie-rating-service/internal/application/models/request"
	"movie-rating-service/internal/application/models/response"
	"movie-rating-service/internal/application/service"
	validate "movie-rating-service/internal/application/validator"
	"time"
)

/*
No need to create unnecessary interface for controllers, because:
- I will not mock controllers functions
- There is no parent layer so creating mocks would be unnecessary, http tests can trigger real functions, we need to mock service and repository operations
- I want to prevent interface pollution
*/

type userController struct {
	userService service.UserService
}

func NewUserController(app *fiber.App, userService service.UserService) {
	authMiddleware := middleware.NewAuthMiddleware(config.Cfg.JWTSecret)

	controller := &userController{userService: userService}

	app.Post("/user", controller.CreateUser)
	app.Get("/user/:id", authMiddleware.AdminHandler, controller.GetUser)

	app.Post("/login", controller.Login)
}

// @Summary Create User
// @Tags User
// @Success 200 {object} response.SuccessResponse{data=response.CreateUser}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Router /user [post]
func (c *userController) CreateUser(ctx *fiber.Ctx) error {
	var req request.CreateUser
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	req.Password = string(hashedPassword)

	// Save user
	res, err := c.userService.Create(ctx.UserContext(), req)
	if err != nil {
		slog.Info("User could not create")
		return err
	}

	slog.Info("User created", "user_id", res.ID)
	return ctx.Status(fiber.StatusCreated).JSON(response.Success(res))
}

// @Summary GetByID User
// @Tags User
// @Param id path string true "User Id"
// @Success 200 {object} response.SuccessResponse{data=response.GetUser}
// @Success 400 {object} response.ErrorResponse
// @Success 500 {object} response.ErrorResponse
// @Router /user/:id [get]
func (c *userController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	req := request.GetUser{ID: cast.ToUint(id)}

	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	res, err := c.userService.Get(ctx.UserContext(), req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.Success(res))
}

// @Summary Login
// @Tags User
// @Success 200 {object} response.SuccessResponse{data=string} // JWT token
// @Success 400 {object} response.ErrorResponse
// @Success 401 {object} response.ErrorResponse
// @Router /login [post]
func (c *userController) Login(ctx *fiber.Ctx) error {
	var req request.Login
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := validate.V.Struct(req)
	if err != nil {
		return err
	}

	user, err := c.userService.IsAuthorized(ctx.UserContext(), req)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"isAdmin":  user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.Cfg.JWTSecret))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token could not be signed"})
	}

	return ctx.JSON(response.Success(response.Login{
		Username: req.Username,
		Token:    tokenString}))
}
