package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := gin.New()
	r.Use(gin.Logger())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "it still work?? : "+port+" "+os.Getenv("ENV")+" "+os.Getenv("DEV_MONGODB_URI"))
	})

	// r.Run(":" + port)

	srv := http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("server start at : " + srv.Addr)

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("HTTP server Shutdown: " + err.Error())
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("HTTP server ListenAndServe: " + err.Error())
		return
	}

	<-idleConnsClosed
}
