package app

import (
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	app "github.com/BrobridgeOrg/vibration-api-service/app/interface"
	pb "github.com/BrobridgeOrg/vibration-api-service/pb"
	timer "github.com/BrobridgeOrg/vibration-api-service/services/timer"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
	host     string
	Timer    *timer.Service
}

func NewGRPCServer(host string) *GRPCServer {
	return &GRPCServer{
		server: grpc.NewServer(),
		host:   host,
	}
}

func (gs *GRPCServer) Init(a *App) error {

	// Preparing listener
	lis := a.connectionListener.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
	)

	gs.listener = lis

	log.WithFields(log.Fields{
		"host": gs.host,
	}).Info("Preparing gRPC server")

	// Services
	timerService := timer.CreateService(app.AppImpl(a))
	a.grpcServer.Timer = timerService
	pb.RegisterTimerServer(gs.server, timerService)

	log.WithFields(log.Fields{
		"service": "Timer",
	}).Info("Registered service")

	return nil
}

func (gs *GRPCServer) Serve() error {

	log.WithFields(log.Fields{
		"host": gs.host,
	}).Info("Starting gRPC server")

	// Starting server
	if err := gs.server.Serve(gs.listener); err != cmux.ErrListenerClosed {
		log.Error(err)
		return err
	}

	return nil
}
