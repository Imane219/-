package service

import (
	"context"
	"contrplatform/global"
	"contrplatform/internal/tester_pool"
	"github.com/gin-gonic/gin"
)

type Service struct {
	ctx context.Context
	pool *tester_pool.TesterPool
	//Session sessions.Session
}

func New(ctx *gin.Context) *Service {
	return &Service{
		ctx: ctx.Request.Context(),
		pool: global.TesterPool,
		//Session: sessions.Default(ctx),
	}
}

func (svc *Service) State(id string) tester_pool.TesterState {
	return svc.pool.State(id)
}

func (svc *Service) GeneID() string {
	return svc.pool.GeneID()
}
