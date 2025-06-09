package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"movie-rating-service/config"
	"strings"
)

type AuthMiddleware interface {
	AdminHandler(ctx *fiber.Ctx) error
	UserHandler(ctx *fiber.Ctx) error
}

type authMiddleware struct {
	Secret string
}

func NewAuthMiddleware(secret string) AuthMiddleware {
	return &authMiddleware{Secret: secret}
}

func (a *authMiddleware) AdminHandler(ctx *fiber.Ctx) error {
	claims, err := authBase(ctx)
	if err != nil {
		return err
	}

	if val, ok := claims["isAdmin"]; ok && val.(bool) == false {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not allowed to access this resource"})
	}
	ctx.Locals("user", claims)

	return ctx.Next()
}

func (a *authMiddleware) UserHandler(ctx *fiber.Ctx) error {
	claims, err := authBase(ctx)
	if err != nil {
		return err
	}

	ctx.Locals("user", claims)

	return ctx.Next()
}

func authBase(ctx *fiber.Ctx) (jwt.MapClaims, error) {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	tokenString := removeBearer(authHeader)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(config.Cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}
	return claims, nil
}

func removeBearer(tokenStr string) string {
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.Split(tokenStr, " ")[1]
	}
	return tokenStr
}
