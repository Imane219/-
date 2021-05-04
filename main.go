package main

import (
	"contrplatform/configs"
	"contrplatform/internal/routers"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(configs.RunMode)
	router := routers.NewRouter()
	_ = router.Run(":7777")
}
