package routers

import (
	"contrplatform/configs"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(20*time.Minute),
	})
	r.Use(sessions.Sessions("my_session",store))
	r.StaticFS("/static",http.Dir(configs.UploadSavePath))
	r.StaticFile("/","static/welcome.html")
	r.StaticFile("/main","static/main.html")
	r.POST("/upload", UploadContract)

	return r
}
