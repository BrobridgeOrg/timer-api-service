package timer

import (
	"time"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"

	app "github.com/BrobridgeOrg/timer-api-service/app/interface"
	pb "github.com/BrobridgeOrg/timer-api-service/pb"
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

	// Generate timer ID
	timerID := uuid.NewV1().String()

	event := &pb.TimerCreation{
		TimerID: timerID,
		Info: &pb.TimerInfo{
			Payload:  in.Payload,
			Callback: in.Callback,
		},
	}

	// Setup time
	switch in.Mode.Mode {
	case "appointment":
		// in.Mode.Timestamp
		event.Timestamp = in.Mode.Timestamp
	case "countdown":
		event.Timestamp = uint64(time.Now().Unix()) + uint64(in.Mode.Interval)
	}

	log.WithFields(log.Fields{
		"id": timerID,
	}).Info("Created timer")

	// Convert message to bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return &pb.CreateTimerReply{}, err
	}

	// Send message to specific room
	ebClient := service.app.GetEventBus().GetClient()
	err = ebClient.Publish("timer.timerCreated", data)
	if err != nil {
		return &pb.CreateTimerReply{}, err
	}

	return &pb.CreateTimerReply{
		TimerID: timerID,
	}, nil
}

func (service *Service) DeleteTimer(ctx context.Context, in *pb.DeleteTimerRequest) (*pb.DeleteTimerReply, error) {

	event := &pb.TimerDeletion{
		TimerID: in.TimerID,
	}

	log.WithFields(log.Fields{
		"id": in.TimerID,
	}).Info("Deleted timer")

	// Convert message to bytes
	data, err := proto.Marshal(event)
	if err != nil {
		return &pb.DeleteTimerReply{}, err
	}

	// Send message to specific room
	ebClient := service.app.GetEventBus().GetClient()
	err = ebClient.Publish("timer.timerDeleted", data)
	if err != nil {
		return &pb.DeleteTimerReply{}, err
	}

	return &pb.DeleteTimerReply{}, nil
}
