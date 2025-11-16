package httpin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewGinRouter(h *AuthHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	a := r.Group("/auth/v1")
	{
		a.POST("/register", h.Register)
		a.POST("/login", h.Login)
		a.POST("/refresh", h.Refresh)
		a.POST("/logout", h.Logout)

		pr := a.Group("/password/reset")
		{
			pr.POST("/start", h.StartReset)
			pr.POST("/confirm", h.ConfirmReset)
		}

		a.GET("/validate", h.Validate)
	}

	return r
}
