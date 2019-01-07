package rest

import (
	"golang-my-events-example/lib/msgqueue"
	"golang-my-events-example/lib/persistence"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

func ServeAPI(listenAddress string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) {
	r := mux.NewRouter()
	r.Methods("post").Path("/api/events/{eventID}/bookings").Handler(&CreateBookingHandler{eventEmitter, database})

	srv := http.Server{
		Handler:      handlers.CORS()(r),
		Addr:         listenAddress,
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	srv.ListenAndServe()
}
