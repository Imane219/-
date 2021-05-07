package main

import (
	"context"
	"contrplatform/configs"
	"contrplatform/global"
	"contrplatform/internal/routers"
	"contrplatform/internal/detectorpool"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	global.DetectorPool = detectorpool.New()
}

func main() {
	gin.SetMode(configs.RunMode)
	server := &http.Server{
		Addr: configs.HttpPort,
		Handler:  routers.NewRouter(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server.ListenAndServe err: %v",err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if err:=server.Shutdown(ctx);err!=nil{
		log.Fatalf("Server forced to shutdown: %v",err)
	}
	global.DetectorPool.Delete()
	log.Println("Server exiting")
}
