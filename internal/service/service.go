package service

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Service struct {
	ctx context.Context
	Session sessions.Session
}

func New(ctx *gin.Context) *Service {
	return &Service{
		ctx: ctx.Request.Context(),
		Session: sessions.Default(ctx),
	}
}
