package app

import (
	"net/http"

	"timer-api-service/app/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
)

func (a *App) InitHTTPServer(host string) error {

	lis := a.connectionListener.Match(cmux.HTTP1Fast())

	//	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	api.InitTimerAPI(a.grpcServer.Timer, r)

	s := &http.Server{
		Handler: r,
	}

	log.WithFields(log.Fields{
		"host": host,
	}).Info("Starting HTTP server on " + host)

	// Starting server
	if err := s.Serve(lis); err != cmux.ErrListenerClosed {
		log.Fatal(err)
		return err
	}

	return nil
}
