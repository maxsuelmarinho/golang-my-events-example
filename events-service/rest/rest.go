package rest

import (
	"log"
	"net/http"

	"github.com/maxsuelmarinho/golang-my-events-example/lib/msgqueue"
	"github.com/maxsuelmarinho/golang-my-events-example/lib/persistence"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ServeAPI(endpoint string, tlsendpoint string, dbHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := newEventHandler(dbHandler, eventEmitter)

	r := mux.NewRouter()
	var api = r.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	api.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			/*
				if r.Header.Get("x-auth-token") != "admin" {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			*/

			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})

	eventsrouter := api.PathPrefix("/events").Subrouter()

	eventsrouter.Methods("GET").Path("/{SearchCriteria}/{search}").HandlerFunc(handler.findEventHandler)
	eventsrouter.Methods("GET").Path("").HandlerFunc(handler.allEventHandler)
	eventsrouter.Methods("GET").Path("/{eventID}").HandlerFunc(handler.oneEventHandler)
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	locationRouter := api.PathPrefix("/locations").Subrouter()
	locationRouter.Methods("GET").Path("").HandlerFunc(handler.allLocationHandler)
	locationRouter.Methods("POST").Path("").HandlerFunc(handler.newLocationHandler)

	server := handlers.CORS()(r)

	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)
	go func() {
		httptlsErrChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", server)
	}()

	go func() {
		httpErrChan <- http.ListenAndServe(endpoint, server)
	}()

	return httpErrChan, httptlsErrChan

}
