package app

import (
	"net"
	"net/http"

	"github.com/BrobridgeOrg/vibration-api-service/app/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
)

type HTTPServer struct {
	server   *http.Server
	listener net.Listener
	host     string
}

func NewHTTPServer(host string) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{},
		host:   host,
	}
}

func (hs *HTTPServer) Init(a *App) error {

	// Preparing listener
	lis := a.connectionListener.Match(cmux.HTTP1Fast())
	hs.listener = lis

	r := gin.Default()

	// APIs
	api.InitTimerAPI(a.grpcServer.Timer, r)

	hs.server.Handler = r

	return nil
}

func (hs *HTTPServer) Serve() error {

	log.WithFields(log.Fields{
		"host": hs.host,
	}).Info("Starting HTTP server")

	// Starting server
	if err := hs.server.Serve(hs.listener); err != cmux.ErrListenerClosed {
		log.Error(err)
		return err
	}

	return nil
}
