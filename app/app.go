package app

import (
	"strconv"
	"time"

	"github.com/BrobridgeOrg/vibration-api-service/app/eventbus"
	app "github.com/BrobridgeOrg/vibration-api-service/app/interface"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
)

type App struct {
	id                 uint64
	flake              *sonyflake.Sonyflake
	eventbus           *eventbus.EventBus
	connectionListener cmux.CMux
	grpcServer         *GRPCServer
	httpServer         *HTTPServer
	isReady            bool
}

func CreateApp() *App {

	// Genereate a unique ID for instance
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return nil
	}

	idStr := strconv.FormatUint(id, 16)

	// exposed port
	port := strconv.Itoa(viper.GetInt("service.port"))

	a := &App{
		id:         id,
		flake:      flake,
		grpcServer: NewGRPCServer(":" + port),
		httpServer: NewHTTPServer(":" + port),
		isReady:    false,
	}

	a.eventbus = eventbus.CreateEventBus(
		viper.GetString("service.event_server"),
		viper.GetString("service.event_cluster_id"),
		idStr,
		func(natsConn *nats.Conn) {

			for {
				log.Warn("re-connect to event server")

				// Connect to NATS Streaming
				err := a.eventbus.Connect()
				if err != nil {
					log.Error("Failed to connect to event server")
					time.Sleep(time.Duration(1) * time.Second)
					continue
				}

				a.isReady = true

				break
			}
		},
		func(natsConn *nats.Conn) {
			a.isReady = false
		},
	)

	return a
}

func (a *App) Init() error {

	log.WithFields(log.Fields{
		"a_id": a.id,
	}).Info("Starting application")

	// Initializing connection listener
	port := strconv.Itoa(viper.GetInt("service.port"))
	err := a.CreateConnectionListener(":" + port)
	if err != nil {
		return err
	}

	// Connect to event server
	err = a.eventbus.Connect()
	if err != nil {
		return err
	}

	// Initialize GRPC server
	err = a.grpcServer.Init(a)
	if err != nil {
		return err
	}

	// Initialize HTTP server
	err = a.httpServer.Init(a)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Uninit() {
}

func (a *App) Run() error {

	// gRPC
	go func() {
		err := a.grpcServer.Serve()
		if err != nil {
			log.Error(err)
		}
	}()

	// HTTP
	go func() {
		err := a.httpServer.Serve()
		if err != nil {
			log.Error(err)
		}
	}()

	err := a.Serve()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) GetEventBus() app.EventBusImpl {
	return app.EventBusImpl(a.eventbus)
}
