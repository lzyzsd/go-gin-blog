package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lzyzsd/go-gin-blog/models"
	"github.com/lzyzsd/go-gin-blog/pkg/gredis"
	"github.com/lzyzsd/go-gin-blog/pkg/logging"
	"github.com/lzyzsd/go-gin-blog/pkg/setting"
	"github.com/lzyzsd/go-gin-blog/routers"
)

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()

	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	router := routers.InitRouter()
	s := &http.Server{
		Addr:           endPoint,
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
