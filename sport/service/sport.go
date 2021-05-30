package service

import (
	"fmt"
	"time"

	"git.neds.sh/matty/entain/racing/proto/sport"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/net/context"
)

type Sport interface {
	// ListSportEvents will return a collection of sport events.
	ListSportEvents(ctx context.Context, in *sport.ListSportEventsRequest) (*sport.ListSportEventsResponse, error)
}

// sportsService implements the Sport interface.
type sportService struct {
	sportsRepo []*sport.SportEvent
}

// NewSportsService instantiates and returns a new sportService.
func NewSportService() Sport {
	var dummyData []*sport.SportEvent
	ts, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		fmt.Println(err)
	}
	// Keeping it simple with hardcoded value, would be better to get this results from a db
	dummyData = append(dummyData, &sport.SportEvent{
		Id:                  1,
		Name:                "test name",
		AdvertisedStartTime: ts,
	})
	return &sportService{dummyData}
}

func (s *sportService) ListSportEvents(ctx context.Context, in *sport.ListSportEventsRequest) (*sport.ListSportEventsResponse, error) {
	return &sport.ListSportEventsResponse{SportEvents: s.sportsRepo}, nil
}
