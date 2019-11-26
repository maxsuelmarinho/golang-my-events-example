package contracts

type EventBookedEvent struct {
	EventID string `json:"eventId"`
	UserID  string `json:"userId"`
}

func (c *EventBookedEvent) EventName() string {
	return "event.booked"
}

func (e *EventBookedEvent) PartitionKey() string {
	return e.EventID
}
