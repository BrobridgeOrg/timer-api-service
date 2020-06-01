package timer

import (
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	app "timer-api-service/app/interface"
	pb "timer-api-service/pb"
)

type Service struct {
	app app.AppImpl
}

func CreateService(a app.AppImpl) *Service {

	// Preparing service
	service := &Service{
		app: a,
	}

	return service
}

func (service *Service) CreateTimer(ctx context.Context, in *pb.CreateTimerRequest) (*pb.CreateTimerReply, error) {

	// TODO: Create timer
	timerID := uuid.NewV1().String()

	switch in.Mode.Mode {
	case "appointment":
		// in.Mode.Timestamp
	case "countdown":
		// in.Mode.Interval
	}

	log.WithFields(log.Fields{
		"id": timerID,
	}).Info("Created timer")

	return &pb.CreateTimerReply{
		TimerID: timerID,
	}, nil
}

func (service *Service) DeleteTimer(ctx context.Context, in *pb.DeleteTimerRequest) (*pb.DeleteTimerReply, error) {

	// TODO: Delete timer

	log.WithFields(log.Fields{
		"id": in.TimerID,
	}).Info("Deleted timer")

	return &pb.DeleteTimerReply{}, nil
}
