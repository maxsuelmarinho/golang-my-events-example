package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang-my-events-example/contracts"
	"golang-my-events-example/lib/msgqueue"
	"golang-my-events-example/lib/persistence"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type eventServiceHandler struct {
	dbhandler    persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func (eh *eventServiceHandler) findEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{error: No search criteria found, you can either search by id via /id/4 to search by name via /name/coldplayconcert}`)
		return
	}

	searchKey, ok := vars["search"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{error: No search keys found, you can either search by id via /id/4 to search by name via /name/coldplayconcert}`)
		return
	}

	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbhandler.FindEventByName(searchKey)
	case "id":
		id, err := hex.DecodeString(searchKey)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error occured %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allEventHandler(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{error: Error occurer while trying to find all available events %s}", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{error: Error occured while trying to find all available events to JSON %s}", err)
	}
}

func (eh *eventServiceHandler) oneEventHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["eventID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "missing route parameter 'eventID'")
		return
	}

	eventIDBytes, _ := hex.DecodeString(eventID)
	event, err := eh.dbhandler.FindEvent(eventIDBytes)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "event with id %s was not found", eventID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) newEventHandler(w http.ResponseWriter, r *http.Request) {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{error: error occured while decoding event data %s}", err)
		return
	}

	id, err := eh.dbhandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "{error: error occured while persisting event %d %s}", id, err)
		return
	}

	message := contracts.EventCreatedEvent{
		ID:         hex.EncodeToString(id),
		Name:       event.Name,
		LocationID: string(event.Location.ID),
		Start:      time.Unix(event.StartDate, 0),
		End:        time.Unix(event.EndDate, 0),
	}

	eh.eventEmitter.Emit(&message)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) allLocationHandler(w http.ResponseWriter, r *http.Request) {
	locations, err := eh.dbhandler.FindAllLocations()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "could not load locations: %s", err)
		return
	}

	json.NewEncoder(w).Encode(locations)
}

func (eh *eventServiceHandler) newLocationHandler(w http.ResponseWriter, r *http.Request) {
	location := persistence.Location{}

	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "request body could not be unserialized to location: %s", err)
		return
	}

	persistedLocation, err := eh.dbhandler.AddLocation(location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not persist location: %s", err)
		return
	}

	msg := contracts.LocationCreatedEvent{
		ID:      string(persistedLocation.ID),
		Name:    persistedLocation.Name,
		Address: persistedLocation.Address,
		Country: persistedLocation.Country,
		Halls:   persistedLocation.Halls,
	}

	eh.eventEmitter.Emit(&msg)

	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&persistedLocation)
}

func newEventHandler(databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler:    databaseHandler,
		eventEmitter: eventEmitter,
	}
}
