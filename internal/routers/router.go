package routers

import (
	"contrplatform/global"
	"contrplatform/internal/middleware"
	"contrplatform/internal/routers/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	//r.Use(middleware.AccessLog())
	r.Use(middleware.Recovery())

	r.StaticFS("/static", http.Dir(global.AppSetting.FrontendPath))
	r.StaticFile("/", "static/home.html")

	detection := api.NewDetection()
	r.POST("/upload", detection.Upload)
	r.POST("/detect", detection.Start)
	r.POST("/result", detection.GetResult)
	r.POST("/stop", detection.Stop)
	r.DELETE("/reset", detection.Reset)
	return r
}
