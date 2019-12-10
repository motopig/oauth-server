package router

import (
	"github.com/motopig/oauth-server/app/controller"
	"github.com/motopig/oauth-server/app/middlewares"
	"github.com/motopig/oauth-server/app/utils"
	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
	"net/http"
)

func Load(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/login", controller.LoginHandler)
	api.POST("/register", controller.RegisterHandler)
	api.GET("/auth", controller.AuthHandler)

	auth := api.Group("/oauth2")
	auth.Use(middlewares.LoginCheck())
	{
		auth.GET("/authorize", ginserver.HandleAuthorizeRequest)
		auth.POST("/token", ginserver.HandleTokenRequest)
	}

	cfg := ginserver.Config{
		ErrorHandleFunc: customErrorHandler,
	}

	status := api.Group("/status")
	status.Use(ginserver.HandleTokenVerify(cfg))
	{
		status.GET("/check", checkLogin)
	}
}

func customErrorHandler(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusForbidden, utils.TokenInvalid)
}

func checkLogin(c *gin.Context) {
	ti, exists := c.Get(ginserver.DefaultConfig.TokenKey)
	if exists {
		c.JSON(http.StatusOK, ti)
		return
	}
	c.JSON(http.StatusOK, utils.Err)
}
