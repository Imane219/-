package routers

import (
	"contrplatform/global"
	"contrplatform/internal/service"
	"contrplatform/internal/tester_pool"
	"contrplatform/pkg/app"
	"contrplatform/pkg/errcode"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Detection struct {}

func NewDetection() Detection {
	return Detection{}
}

func (d Detection) Start(c *gin.Context)  {
	param := service.StartDetectionRequest{}
	response := app.NewResponse(c)

	if errs:= app.BindParams(c,&param); errs != nil {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	state := svc.State(param.ID)
	if state != tester_pool.StateInit {
		details := fmt.Sprintf("id:%s, state:%s - not %s",param.ID,state,
			tester_pool.StateInit)
		errRsp := errcode.ErrorTeseterState.WithDetails(details)
		response.ToErrorResponse(errRsp)
		return
	}

	if err:=svc.StartDetection(&param);err!=nil{
		errRsp := errcode.ErrorStartDetectionFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{
		"id":param.ID,
		"state":global.TesterPool.State(param.ID),
	})
}

func (d Detection) GetResult(c *gin.Context)  {
	param := service.GetResultRequest{}
	response := app.NewResponse(c)
	if errs := app.BindParams(c,&param); errs!=nil{
		errRsp:=errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	state := svc.State(param.ID)
	if state != tester_pool.StateRunning || state != tester_pool.StateStopped {
		details := fmt.Sprintf("id:%s, state:%s - not %s or %s",param.ID,state,
			tester_pool.StateRunning,tester_pool.StateStopped)
		errRsp := errcode.ErrorTeseterState.WithDetails(details)
		response.ToErrorResponse(errRsp)
		return
	}

	results, err:= svc.GetResult(&param)
	if err != nil {
		errRsp := errcode.ErrorGetResultFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{
		"id":param.ID,
		"state":global.TesterPool.State(param.ID),
		"outputs":results,
	})
}

func (d *Detection) Stop(c *gin.Context) {
	param := service.StopDetectionRequest{}
	response := app.NewResponse(c)
	if errs := app.BindParams(c,&param); errs!=nil{
		errRsp:=errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}
	svc := service.New(c)
	state := svc.State(param.ID)
	if state == tester_pool.StateRunning {
		svc.StopDetection(&param)
	}else if state != tester_pool.StateStopped{
		details := fmt.Sprintf("id:%s, state:%s - not %s or %s",param.ID,state,
			tester_pool.StateRunning,tester_pool.StateStopped)
		errRsp := errcode.ErrorTeseterState.WithDetails(details)
		response.ToErrorResponse(errRsp)
		return
	}

	results, err:= svc.GetResult(&service.GetResultRequest{
		ID: param.ID,
	})
	if err != nil {
		errRsp := errcode.ErrorGetResultFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{
		"id":param.ID,
		"state":global.TesterPool.State(param.ID),
		"outputs":results,
	})
}

func (d *Detection) Reset(c *gin.Context) {
	param := service.ResetDetectionRequest{}
	response := app.NewResponse(c)
	if errs := app.BindParams(c,&param); errs!=nil{
		errRsp:=errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}
	svc := service.New(c)
	state := svc.State(param.ID)
	if state != tester_pool.StateStopped{
		details := fmt.Sprintf("id:%s, state:%s - not %s",param.ID,state,
			tester_pool.StateStopped)
		errRsp := errcode.ErrorTeseterState.WithDetails(details)
		response.ToErrorResponse(errRsp)
		return
	}
	if err:= svc.ResetDetection(&param); err != nil {
		errRsp := errcode.ErrorResetDetectionFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse("removed")
}