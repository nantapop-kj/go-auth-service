package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/dto"
	"github.com/nantapop-kj/go-auth-service/enums"
	"github.com/nantapop-kj/go-auth-service/service"
	"github.com/nantapop-kj/go-auth-service/utils"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Service: s,
	}
}

func (h *AuthHandler) HandlerRegister(c *fiber.Ctx) error {
	ctx := c.Context()
	isProd := config.GetEnv("APP_ENV", "") == string(enums.EnvProd)
	cookieExp := config.GetEnv("COOKIE_EXPIRE_HOURS", "")

	hours, err := strconv.Atoi(cookieExp)
	if err != nil || hours <= 0 {
		hours = 24
	}
	var cookiExpireHours int = hours

	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	file, err := c.FormFile("profile_image")
	if err != nil {
		return utils.ResponseMessageError(c, fiber.StatusBadRequest, "Profile image is required")
	}

	req := dto.RegisterRequest{
		Email:        email,
		Password:     password,
		FirstName:    firstName,
		LastName:     lastName,
		ProfileImage: file,
	}

	if err := utils.Validate.Struct(req); err != nil {
		return utils.ResponseMessageError(c, fiber.StatusBadRequest, utils.ValidationErrorToString(err))
	}

	_, token, err := h.Service.ServiceRegister(ctx, req)
	if err != nil {
		return utils.ResponseMessageError(c, fiber.StatusBadRequest, err.Error())
	}

	sameSite := fiber.CookieSameSiteLaxMode
	if isProd {
		sameSite = fiber.CookieSameSiteNoneMode
	}

	c.Cookie(&fiber.Cookie{
		Name:     config.GetEnv("COOKIE_NAME", ""),
		Value:    token,
		Expires:  time.Now().Add(time.Duration(cookiExpireHours) * time.Hour),
		HTTPOnly: true,
		Secure:   isProd,
		SameSite: sameSite,
		Path:     "/",
	})

	return utils.ResponseMessageSuccess(c, fiber.StatusOK, "User registered successfully.")
}
