package middleware

import (
	"bytes"
	"contrplatform/global"
	"contrplatform/pkg/logger"
	"github.com/gin-gonic/gin"
	"time"
)

//访问日志
type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

//Write 方法重写(达到写入响应时顺便将响应内容记录到body中)
func (w AccessLogWriter) Write(p []byte) (int, error) {
	//将切片p中数据写入w的Buffer中
	if n,err:=w.body.Write(p); err != nil {
		return n, err
	}
	//对响应进行写操作
	return w.ResponseWriter.Write(p)
}

//访问日志中间件
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWriter := &AccessLogWriter{
			body: bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		//由于AccessLogWriter中有成员ResponseWriter,因此可以实现接口赋值
		//而AccessLogWriter中重写了Write方法,因此调用时会调用该方法
		c.Writer = bodyWriter	//接口的替换

		//请求开始时间
		beginTime := time.Now().Unix()
		c.Next()
		//请求结束时间
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"request":  c.Request.PostForm.Encode(),	//请求
			"response": bodyWriter.body.String(),		//响应
		}
		//记录日志

		global.Logger.WithFields(fields).Infof(
			c,
			"access log: method: %s, status_code: %d, " +
				"begin_time: %d, end_time: %d",
			c.Request.Method,
			bodyWriter.Status(),
			beginTime,
			endTime,
		)
	}
}


