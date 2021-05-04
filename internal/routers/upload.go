package routers

import (
	"contrplatform/internal/service"
	"contrplatform/pkg/app"
	"contrplatform/pkg/errcode"
	"contrplatform/pkg/upload"
	"github.com/gin-gonic/gin"
)

func UploadContract(c *gin.Context)  {
	response:= app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		errRsp := errcode.InvalidParams.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}

	if fileHeader==nil {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	fileType := upload.TypeContract
	svc := service.New(c)

	if svc.GetSessionID() == nil {
		if err:=svc.InitSession();err != nil {
			response.ToErrorResponse(errcode.SessionError.WithDetails(err.Error()))
			return
		}
	}

	if err := svc.UploadFile(upload.FileType(fileType),file,fileHeader); err!= nil {
		errRsp:=errcode.ErrorUploadFileFail.WithDetails(err.Error())
		response.ToErrorResponse(errRsp)
		return
	}
	if err := svc.SetSessionFile(fileHeader.Filename); err != nil {
		response.ToErrorResponse(errcode.SessionError.WithDetails(err.Error()))
		return
	}
	response.ToResponse(gin.H{
		"uploaded_file":svc.GetSessionFiles(),
	})
}
