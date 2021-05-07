package service

import (
	"context"
	"contrplatform/global"
	"contrplatform/internal/detectorpool"
	"github.com/gin-gonic/gin"
)

type Service struct {
	ctx  context.Context
	pool *detectorpool.DetectorPool
	//Session sessions.Session
}

func New(ctx *gin.Context) *Service {
	return &Service{
		ctx:  ctx.Request.Context(),
		pool: global.DetectorPool,
		//Session: sessions.Default(ctx),
	}
}

func (svc *Service) DetectorState(id string) detectorpool.DetectorState {
	return svc.pool.DetectorState(id)
}

func (svc *Service) GeneID() string {
	return svc.pool.GeneID()
}
