package common

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"movie-rating-service/internal/application/models/response"
)

const uniqueValidationErr = "23505"

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == uniqueValidationErr
	}
	return false
}

func ErrorHandler() func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		if errors.As(err, &validator.ValidationErrors{}) {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.Error("Bad request.", err.Error()))
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(response.Error("Record not found.", err.Error()))
		}
		if isUniqueViolation(err) {
			return ctx.Status(fiber.StatusBadRequest).JSON(response.Error("Cannot use same values.", err.Error()))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(response.Error("Internal server error.", err.Error()))
	}
}
