package handler

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nantapop-kj/go-auth-service/config"
	"github.com/nantapop-kj/go-auth-service/dto"
	"github.com/nantapop-kj/go-auth-service/middleware"
	"github.com/nantapop-kj/go-auth-service/service"
	"github.com/nantapop-kj/go-auth-service/utils"
)

type AuthHandler struct {
	cfg     *config.Config
	service *service.AuthService
	tokens  *service.TokenService
}

func NewAuthHandler(cfg *config.Config, s *service.AuthService, tokens *service.TokenService) *AuthHandler {
	return &AuthHandler{cfg: cfg, service: s, tokens: tokens}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if msg := bindAndValidate(c, &req); msg != "" {
		return utils.ResponseMessageError(c, fiber.StatusBadRequest, msg)
	}
	res, err := h.service.Register(c.UserContext(), req, clientInfo(c))
	if err != nil {
		return h.fail(c, err)
	}
	h.setAuthCookies(c, res)
	return utils.ResponseData(c, fiber.StatusCreated, "User registered successfully.", res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if msg := bindAndValidate(c, &req); msg != "" {
		return utils.ResponseMessageError(c, fiber.StatusBadRequest, msg)
	}
	res, err := h.service.Login(c.UserContext(), req, clientInfo(c))
	if err != nil {
		return h.fail(c, err)
	}
	h.setAuthCookies(c, res)
	return utils.ResponseData(c, fiber.StatusOK, "Login successful.", res)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	res, err := h.service.Refresh(c.UserContext(), h.refreshToken(c))
	if err != nil {
		h.clearAuthCookies(c)
		return h.fail(c, err)
	}
	h.setAuthCookies(c, res)
	return utils.ResponseData(c, fiber.StatusOK, "Token refreshed.", res)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	if err := h.service.Logout(c.UserContext(), h.refreshToken(c), clientInfo(c)); err != nil {
		return h.fail(c, err)
	}
	h.clearAuthCookies(c)
	return utils.ResponseMessageSuccess(c, fiber.StatusOK, "Logged out.")
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userID, _ := c.Locals(middleware.LocalUserID).(uuid.UUID)
	user, err := h.service.GetUser(c.UserContext(), userID)
	if err != nil {
		return h.fail(c, err)
	}
	return utils.ResponseData(c, fiber.StatusOK, "Profile retrieved successfully.", user)
}

func (h *AuthHandler) refreshToken(c *fiber.Ctx) string {
	var req dto.RefreshRequest
	_ = c.BodyParser(&req)
	if req.RefreshToken == "" && h.cfg.CookieEnabled {
		return c.Cookies(h.cfg.CookieRefreshName)
	}
	return req.RefreshToken
}

func (h *AuthHandler) fail(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	switch {
	case errors.Is(err, service.ErrEmailTaken):
		status = fiber.StatusConflict
	case errors.Is(err, service.ErrInvalidCredentials), errors.Is(err, service.ErrInvalidToken):
		status = fiber.StatusUnauthorized
	case errors.Is(err, service.ErrUserInactive):
		status = fiber.StatusForbidden
	case errors.Is(err, service.ErrUserNotFound):
		status = fiber.StatusNotFound
	case errors.Is(err, service.ErrTooManyAttempts):
		status = fiber.StatusTooManyRequests
		c.Set(fiber.HeaderRetryAfter, strconv.Itoa(int(h.cfg.LoginLockDuration.Seconds())))
	}

	if status == fiber.StatusInternalServerError {
		log.Printf("%s %s: %v", c.Method(), c.Path(), err)
		return utils.ResponseMessageError(c, status, "internal server error")
	}
	return utils.ResponseMessageError(c, status, err.Error())
}

func (h *AuthHandler) setAuthCookies(c *fiber.Ctx, res *dto.TokenResponse) {
	if !h.cfg.CookieEnabled {
		return
	}
	now := time.Now()
	c.Cookie(h.cookie(h.cfg.CookieName, res.AccessToken, now.Add(h.tokens.AccessTTL())))
	c.Cookie(h.cookie(h.cfg.CookieRefreshName, res.RefreshToken, now.Add(h.tokens.RefreshTTL())))
}

func (h *AuthHandler) clearAuthCookies(c *fiber.Ctx) {
	if !h.cfg.CookieEnabled {
		return
	}
	expired := time.Unix(0, 0)
	c.Cookie(h.cookie(h.cfg.CookieName, "", expired))
	c.Cookie(h.cookie(h.cfg.CookieRefreshName, "", expired))
}

func (h *AuthHandler) cookie(name, value string, expires time.Time) *fiber.Cookie {
	sameSite := fiber.CookieSameSiteLaxMode
	switch h.cfg.CookieSameSite {
	case "strict":
		sameSite = fiber.CookieSameSiteStrictMode
	case "none":
		sameSite = fiber.CookieSameSiteNoneMode
	}
	return &fiber.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   h.cfg.CookieDomain,
		Expires:  expires,
		HTTPOnly: true,
		Secure:   h.cfg.CookieSecure || sameSite == fiber.CookieSameSiteNoneMode,
		SameSite: sameSite,
	}
}

func bindAndValidate(c *fiber.Ctx, req any) string {
	if err := c.BodyParser(req); err != nil {
		return "invalid request body"
	}
	if err := utils.Validate.Struct(req); err != nil {
		return utils.ValidationErrorToString(err)
	}
	return ""
}

func clientInfo(c *fiber.Ctx) service.ClientInfo {
	return service.ClientInfo{IP: c.IP(), UserAgent: c.Get(fiber.HeaderUserAgent)}
}
