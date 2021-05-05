package main

import (
	"contrplatform/configs"
	"contrplatform/global"
	"contrplatform/internal/routers"
	"contrplatform/internal/tester_pool"
	"github.com/gin-gonic/gin"
)

func init() {
	global.TesterPool = tester_pool.New()
}

func main() {
	gin.SetMode(configs.RunMode)
	router := routers.NewRouter()
	_ = router.Run(":7777")
}
