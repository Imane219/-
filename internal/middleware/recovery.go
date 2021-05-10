package middleware

import (
	"contrplatform/global"
	"contrplatform/pkg/app"
	"contrplatform/pkg/errcode"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err:=recover(); err != nil {
				global.Logger.WithCallerFrames().Errorf(c, "panic recover err: %v", err)
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}