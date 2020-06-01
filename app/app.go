package app

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/sony/sonyflake"
	"github.com/spf13/viper"
)

type App struct {
	id                 uint64
	flake              *sonyflake.Sonyflake
	connectionListener cmux.CMux
	grpcServer         *GRPCServer
}

func CreateApp() *App {

	// Genereate a unique ID for instance
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		return nil
	}

	return &App{
		id:         id,
		flake:      flake,
		grpcServer: &GRPCServer{},
	}
}

func (a *App) Init() error {

	log.WithFields(log.Fields{
		"a_id": a.id,
	}).Info("Starting application")

	return nil
}

func (a *App) Uninit() {
}

func (a *App) Run() error {

	port := strconv.Itoa(viper.GetInt("service.port"))
	err := a.CreateConnectionListener(":" + port)
	if err != nil {
		return err
	}

	// HTTP
	go func() {
		err := a.InitHTTPServer(":" + port)
		if err != nil {
			log.Error(err)
		}
	}()

	// gRPC
	go func() {
		err := a.InitGRPCServer(":" + port)
		if err != nil {
			log.Error(err)
		}
	}()

	err = a.Serve()
	if err != nil {
		return err
	}

	return nil
}
