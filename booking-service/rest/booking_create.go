package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang-my-events-example/contracts"
	"golang-my-events-example/lib/msgqueue"
	"golang-my-events-example/lib/persistence"
	"time"
)

type eventRef struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type createBookingRequest struct {
	Seats int `json:"seats`
}

type createBookinhResponse struct {
	ID    string   `json:"id"`
	Event eventRef `json:"event"`
}

type CreateBookingHandler struct {
	eventEmitter msgqueue.EventEmitter
	database     persistence.DatabaseHandler
}

func (h *CreateBookingHandler) ServeHTTP(w http.ResponseWritter, r *http.Request) {
	routeVars := mux.Vars(r)
	eventID, ok := routeVars["eventID"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "missing route parameter 'eventID'")
		return
	}

	eventIDMongo, _ := hex.DecodeString(eventID)
	event, err := h.database.FindEvent(eventIDMongo)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "event %s could not be loaded: %s", eventID, err)
	}

	bookingRequest := createBookingRequest{}
	err = json.NewDecoder(r.Body).Decode(&bookingRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "could not decode JSON body: %s", err)
		return
	}

	if bookingRequest.Seats <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "seat number must be positive (was %d)", bookingRequest.Seats)
		return
	}

	eventIDAsBytes, _ := event.ID.MarshalText()
	booking := persistence.Booking{
		Date:    time.Now().Unix(),
		EventID: eventIDAsBytes,
		Seats:   bookingRequest.Seats,
	}

	msg := contracts.EventBookedEvent{
		EventID: event.ID.Hex(),
		UserID:  "someUserID",
	}
	h.eventEmitter.Emit(&msg)

	h.database.AddBookingForUser([]byte("someUserID"), booking)

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(http.StatusCreated)

	json.NewDecoder(w).Encode(&booking)
}
