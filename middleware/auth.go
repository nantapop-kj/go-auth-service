package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/service"
	"github.com/nantapop-kj/go-auth-service/utils"
)

const LocalUserID = "user_id"

func RequireAuth(tokens *service.TokenService, cookieName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := bearerToken(c.Get(fiber.HeaderAuthorization))
		if token == "" && cookieName != "" {
			token = c.Cookies(cookieName)
		}
		if token == "" {
			return utils.ResponseMessageError(c, fiber.StatusUnauthorized, "missing access token")
		}

		claims, err := tokens.ParseAccessToken(token)
		if err != nil {
			return utils.ResponseMessageError(c, fiber.StatusUnauthorized, err.Error())
		}
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return utils.ResponseMessageError(c, fiber.StatusUnauthorized, service.ErrInvalidToken.Error())
		}

		c.Locals(LocalUserID, userID)
		return c.Next()
	}
}

func bearerToken(header string) string {
	const prefix = "bearer "
	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}
	return ""
}
