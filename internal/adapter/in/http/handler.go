package httpin

import (
	"net/http"
	"strings"
	"time"

	"auth/internal/adapter/in/http/dto"
	"auth/internal/domain"
	"auth/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.Service
}

func NewAuthHandler(svc *service.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func bearer(c *gin.Context) string {
	h := c.GetHeader("Authorization")
	if len(h) < 7 {
		return ""
	}
	if strings.ToLower(h[:7]) != "bearer " {
		return ""
	}
	return strings.TrimSpace(h[7:])
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request"})
		return
	}

	tp, err := h.svc.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == domain.ErrAlreadyExists {
			c.AbortWithStatusJSON(http.StatusConflict, dto.ErrorResponse{Error: "email_exists"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal"})
		return
	}

	c.JSON(http.StatusCreated, dto.TokenResponse{
		AccessToken:  tp.AccessToken,
		RefreshToken: tp.RefreshToken,
	})
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request"})
		return
	}

	tp, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCreds:
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid_credentials"})
		case domain.ErrBlockedUser:
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{Error: "blocked"})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  tp.AccessToken,
		RefreshToken: tp.RefreshToken,
	})
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request"})
		return
	}

	tp, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid_refresh"})
		return
	}

	c.JSON(http.StatusOK, dto.TokenResponse{
		AccessToken:  tp.AccessToken,
		RefreshToken: tp.RefreshToken,
	})
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	_ = c.ShouldBindJSON(&req)
	_ = h.svc.Logout(c.Request.Context(), req.RefreshToken)
	c.Status(http.StatusNoContent)
}

// POST /auth/password/reset/start
func (h *AuthHandler) StartReset(c *gin.Context) {
	var req dto.StartResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request"})
		return
	}

	code, err := h.svc.StartPasswordResetOTP(c.Request.Context(), req.Email)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}
	// lol, for testing purposes
	c.JSON(http.StatusOK, gin.H{"dev_code": code})
}

// POST /auth/password/reset/confirm
func (h *AuthHandler) ConfirmReset(c *gin.Context) {
	var req dto.ConfirmResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "bad_request"})
		return
	}

	if err := h.svc.ConfirmPasswordResetOTP(c.Request.Context(), req.Email, req.Code, req.NewPassword); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid_code"})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /auth/validate
func (h *AuthHandler) Validate(c *gin.Context) {
	access := bearer(c)
	if access == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "no_bearer"})
		return
	}

	sub, claims, err := h.svc.ValidateAccess(c.Request.Context(), access)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid_token"})
		return
	}

	var exp time.Time
	if v, ok := claims["exp"].(float64); ok {
		exp = time.Unix(int64(v), 0).UTC()
	}

	iss, _ := claims["iss"].(string)
	c.JSON(http.StatusOK, dto.ValidateResponse{
		Sub:       sub.String(),
		Issuer:    iss,
		ExpiresAt: exp,
		Claims:    claims,
	})
}
