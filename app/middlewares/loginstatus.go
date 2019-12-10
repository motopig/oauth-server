package middlewares

import (
	"github.com/motopig/oauth-server/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
	"net/http"
)

func LoginCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		store, err := session.Start(nil, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			c.Abort()
		}
		// 检查登录是否成功 没有成功则跳转登录页面
		if store != nil {
			if _, ok := store.Get("LoggedInUserID"); !ok {
				c.JSON(http.StatusForbidden, utils.LoginFailed)
				c.Abort()
			} else {
				c.Next()
			}
		} else {
			c.JSON(http.StatusForbidden, utils.LoginFailed)
			c.Abort()
		}
	}
}
