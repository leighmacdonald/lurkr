package webhook

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/leighmacdonald/lurkr/internal/parser"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	router *gin.Engine
	srv    *http.Server
)

type request struct {
}

func Start(newAnnounce chan parser.Announce) {
	router := gin.Default()
	router.GET("/webhook", func(context *gin.Context) {

	})
	httpSrv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	srv = httpSrv
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("Error listening on http webhook: %v", err)
		}
	}()
}

func Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server shutdown error: %v", err)
		return
	}
	log.Debugf("HTTP server shutdown")
}
