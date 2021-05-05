package errcode

var (
	ErrorUploadFileFail	= NewError(20010001,"上传文件失败")
	ErrorTeseterState = NewError(20010002, "错误的测试状态")

	ErrorStartDetectionFail   = NewError(20020001,"运行漏洞检测脚本失败")
	ErrorGetResultFail        = NewError(20020002,"获取检测结果失败")
	ErrorResetDetectionFail = NewError(20020003,"重置检测状态失败")
)
