package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"

	app "github.com/BrobridgeOrg/vibration-api-service/app/interface"
	pb "github.com/BrobridgeOrg/vibration-api-service/pb"
	timer "github.com/BrobridgeOrg/vibration-api-service/services/timer"
)

type GRPCServer struct {
	Timer *timer.Service
}

func (a *App) InitGRPCServer(host string) error {
	/*
		// Start to listen on port
		lis, err := net.Listen("tcp", host)
		if err != nil {
			log.Fatal(err)
			return err
		}
	*/
	lis := a.connectionListener.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
	)

	log.WithFields(log.Fields{
		"host": host,
	}).Info("Starting gRPC server on " + host)

	// Create gRPC server
	s := grpc.NewServer()

	// Register data source adapter service
	timerService := timer.CreateService(app.AppImpl(a))
	a.grpcServer.Timer = timerService
	pb.RegisterTimerServer(s, timerService)
	//reflection.Register(s)

	log.WithFields(log.Fields{
		"service": "Timer",
	}).Info("Registered service")

	// Starting server
	if err := s.Serve(lis); err != cmux.ErrListenerClosed {
		log.Fatal(err)
		return err
	}

	return nil
}
