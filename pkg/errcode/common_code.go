package errcode

// 公共内部错误码
var (
	Success                  = NewError(0, "成功")
	ServerError              = NewError(10000000, "服务器内部错误")
	InvalidParams            = NewError(10000001, "函数入参错误")
	NotFound                 = NewError(10000002, "未找到")

	//UnauthorizedAuthNotExist = NewError(10000003,
	//	"鉴权失败,找不到对应的AppKey和AppSecret")
	//UnauthorizedTokenError    = NewError(10000004, "鉴权失败,Token错误")
	//UnauthorizedTokenTimeout  = NewError(1000005, "鉴权失败,Token超时")
	//UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败,Token生成失败")
	//TooManyRequests           = NewError(10000007, "请求过多")
)
