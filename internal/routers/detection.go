package routers

import (
	"contrplatform/internal/service"
	"contrplatform/pkg/app"
	"contrplatform/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Detection struct{}

func NewDetection() Detection {
	return Detection{}
}

func (d Detection) Start(c *gin.Context) {
	param := service.StartDetectionRequest{}
	response := app.NewResponse(c)

	if errs := app.BindParams(c, &param); errs != nil {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	if err := svc.StartDetection(&param); err != nil {
		errRsp := errcode.ErrorStartDetectionFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{
		"id":    param.ID,
		"state": svc.DetectorState(param.ID),
	})
}

func (d Detection) GetResult(c *gin.Context) {
	param := service.GetResultRequest{}
	response := app.NewResponse(c)
	if errs := app.BindParams(c, &param); errs != nil {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	result, err := svc.GetResult(&param)
	if err != nil {
		errRsp := errcode.ErrorGetResultFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{
		"id":      param.ID,
		"state":   svc.DetectorState(param.ID),
		"outputs": result,
	})
}

func (d *Detection) Stop(c *gin.Context) {
	param := service.StopDetectionRequest{}
	response := app.NewResponse(c)
	if errs := app.BindParams(c, &param); errs != nil {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}

	svc := service.New(c)
	if err:= svc.StopDetection(&param);err!=nil{
		errRsp := errcode.ErrorStopDetectionFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	results, err := svc.GetResult(&service.GetResultRequest{
		ID: param.ID,
	})
	if err != nil {
		errRsp := errcode.ErrorGetResultFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse(gin.H{
		"id":      param.ID,
		"state":   svc.DetectorState(param.ID),
		"outputs": results,
	})
}

func (d *Detection) Reset(c *gin.Context) {
	param := service.ResetDetectionRequest{}
	response := app.NewResponse(c)
	if errs := app.BindParams(c, &param); errs != nil {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrorResponse(errRsp)
		return
	}
	svc := service.New(c)
	if err := svc.ResetDetection(&param); err != nil {
		errRsp := errcode.ErrorResetDetectionFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	response.ToResponse("removed")
}
