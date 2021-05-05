package app

import (
	"contrplatform/pkg/errcode"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data interface{})  {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK,data)
}

//设置错误响应
func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{
		"code": err.Code(),
		"msg":  err.Msg(),
	}
	details := err.Details()
	//若有详情信息则添加到响应中
	if len(details) > 0 {
		response["details"] = details
	}
	r.Ctx.JSON(err.StatusCode(), response)
}