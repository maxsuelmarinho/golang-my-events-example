package contracts

import "golang-my-events-example/events-service/persistence"

type LocationCreatedEvent struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Address string             `json:"address"`
	Country string             `json:"country"`
	Halls   []persistence.Hall `json:"halls"`
}

func (c *LocationCreatedEvent) EventName() string {
	return "location.created"
}
