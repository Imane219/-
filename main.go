package main

import (
	"context"
	"contrplatform/global"
	"contrplatform/internal/detectorpool"
	"contrplatform/internal/routers"
	"contrplatform/pkg/logger"
	"contrplatform/pkg/setting"
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	if err:=setupSetting();err!=nil {
		log.Fatalf("init.setupSetting err: %v",err)
	}
	if err:=setupLogger(); err!=nil {
		log.Fatalf("init.setupLogger err: %v",err)
	}
	if err:=setupDetectorPool();err!=nil {
		log.Fatalf("init.setupDetectorPool err: %v",err)
	}
}

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	server := &http.Server{
		Addr: ":"+ global.ServerSetting.HttpPort,
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

func setupSetting() error {
	set, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = set.ReadSection("Server",&global.ServerSetting)
	if err != nil {
		return err
	}
	err = set.ReadSection("App",&global.AppSetting)
	if err != nil {
		return err
	}
	err = set.ReadSection("Log",&global.LogSetting)
	if err != nil {
		return err
	}
	err = set.ReadSection("Pool",&global.PoolSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	if global.LogSetting == nil {
		return errors.New("global.LogSetting is nil")
	}
	fileName := global.LogSetting.SavePath+"/"+global.LogSetting.FileName+
		global.LogSetting.FileExt
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename: fileName,
		MaxSize: global.LogSetting.MaxSize,
		MaxAge: global.LogSetting.MaxAge,
		LocalTime: true,
	},"",log.LstdFlags)
	return nil
}

func setupDetectorPool() error {
	if global.PoolSetting == nil {
		return errors.New("global.PoolSetting is nil")
	}
	global.DetectorPool = detectorpool.New(global.PoolSetting)
	return nil
}