package stat

import (
	"LinkShorty/pkg/event"
	"log"
)

type StatService struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

type StatServiceDeps struct {
	EventBus       *event.EventBus
	StatRepository *StatRepository
}

func NewStatService(deps *StatServiceDeps) *StatService {
	return &StatService{
		EventBus:       deps.EventBus,
		StatRepository: deps.StatRepository,
	}
}

func (svc *StatService) AddClick() {
	for msg := range svc.EventBus.Subscribe() {
		if msg.Type == event.EventLinkVisited {
			id, ok := msg.Data.(uint)
			if !ok {
				log.Fatalln("Wrong data type", msg.Data)
				continue
			}
			svc.StatRepository.AddClick(id)
		}
	}
}
