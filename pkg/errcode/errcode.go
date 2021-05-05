package errcode

import (
	"fmt"
	"net/http"
)

//内部错误码
type Error struct {
	code int	//错误码号
	msg string	//错误信息
	details []string	//错误详情
}

//错误码映射表
var codes = map[int]string{}

//新建错误码
func NewError(code int, msg string) *Error {
	if _,ok:=codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已存在,请更换一个",code))
	}
	codes[code]=msg
	return &Error{
		code: code,
		msg: msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码:%d, 错误信息: %s",e.code,e.msg)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

//格式化输出错误信息
func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg,args...)
}

func (e *Error) Details() []string {
	return e.details
}

//创建带错误详情的错误码
func (e *Error) WithDetails(details ...string) *Error {
	newError := *e	//创建的错误码包含原错误码的码号和详情信息
	newError.details = []string{}
	for _,d := range details {
		newError.details = append(newError.details, d)
	}
	return &newError
}

//将内部错误码转换为HTTP状态码
func (e *Error) StatusCode() int {
	switch e.code {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	//case UnauthorizedAuthNotExist.Code():
	//	fallthrough
	//case UnauthorizedTokenError.Code():
	//	fallthrough
	//case UnauthorizedTokenGenerate.Code():
	//	fallthrough
	//case UnauthorizedTokenTimeout.Code():
	//	return http.StatusUnauthorized
	//case NotFound.Code():
	//	return http.StatusNotFound
	//case TooManyRequests.Code():
	//	return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
