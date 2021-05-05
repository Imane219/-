package app

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

//验证请求参数的错误
type ValidError struct {
	Key string
	Message string
}

//验证请求参数的错误列表
type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _,err:=range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func (v *ValidErrors) Error() string {
	return strings.Join(v.Errors(),",")
}

//将请求上下文c的数据字段绑定到参数结构体v中,并进行参数验证
func BindParams(c *gin.Context,v interface{}) ValidErrors {
	var errs ValidErrors
	//将上下文中字段数据绑定到接口v中
	if err := c.ShouldBind(v); err != nil {		//绑定数据出错
		//获取validator的翻译器
		v := c.Value("trans")
		trans,_ := v.(ut.Translator)
		//将错误转换为validator的错误
		vErrs, ok:= err.(validator.ValidationErrors)
		//转换失败
		if !ok {
			return errs
		}
		//将所有错误进行翻译
		for k,v:=range vErrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key: k,
				Message: v,
			})
		}
		return errs
	}
	return nil
}