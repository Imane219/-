package routers

import (
	"contrplatform/configs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFS("/static",http.Dir(configs.UploadSavePath))
	r.StaticFile("/","static/welcome.html")
	r.StaticFile("/main","static/main.html")
	r.POST("/upload", UploadContracts)

	detection := NewDetection()
	r.POST("/detect",detection.Start)
	r.POST("/result",detection.GetResult)
	r.POST("/stop",detection.Stop)
	r.DELETE("/reset",detection.Reset)
	return r
}
