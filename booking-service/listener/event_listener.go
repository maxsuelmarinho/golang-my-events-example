package listener

import (
	"fmt"
	"log"

	"github.com/maxsuelmarinho/golang-my-events-example/contracts"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/persistence"

	"gopkg.in/mgo.v2/bson"
)

type EventProcessor struct {
	EventListener msgqueue.EventListener
	Database      persistence.DatabaseHandler
}

func (p *EventProcessor) ProcessEvents() error {
	log.Println("Listening to events...")

	received, errors, err := p.EventListener.Listen("event.created")
	if err != nil {
		return err
	}

	for {
		select {
		case evt := <-received:
			fmt.Printf("got event %T: %s\n", evt, evt)
			p.handleEvent(evt)
		case err = <-errors:
			log.Printf("received error while processing message: %s", err)
		}
	}
}

func (p *EventProcessor) handleEvent(event msgqueue.Event) {
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		log.Printf("event %s created: %s", e.ID, e)

		if !bson.IsObjectIdHex(e.ID) {
			log.Printf("event %v did not contain valid object ID", e)
		}

		p.Database.AddEvent(persistence.Event{ID: bson.ObjectId(e.ID), Name: e.Name})
	case *contracts.LocationCreatedEvent:
		log.Printf("location %s created: %v", e.ID, e)

		if !bson.IsObjectIdHex(e.ID) {
			log.Printf("event %v did not contain valid object ID", e)
		}

		p.Database.AddLocation(persistence.Location{ID: bson.ObjectId(e.ID)})
	default:
		log.Printf("unknown event: %t", e)
	}
}
