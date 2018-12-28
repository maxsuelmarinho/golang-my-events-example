package rest

import (
	"golang-my-events-example/events-service/msgqueue"
	"golang-my-events-example/events-service/persistence"
	"log"
	"net/http"

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
	eventsrouter.Methods("POST").Path("").HandlerFunc(handler.newEventHandler)

	httpErrChan := make(chan error)
	httptlsErrChan := make(chan error)
	go func() {
		httptlsErrChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", r)
	}()

	go func() {
		httpErrChan <- http.ListenAndServe(endpoint, r)
	}()

	return httpErrChan, httptlsErrChan

}
