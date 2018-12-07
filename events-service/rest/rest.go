package rest

import (
	"golang-my-events-example/events-service/persistence"
	"net/http"

	"github.com/gorilla/mux"
)

func ServeAPI(endpoint string, tlsendpoint string, dbHandler persistence.DatabaseHandler) (chan error, chan error) {
	handler := newEventHandler(dbHandler)
	r := mux.NewRouter()
	eventsrouter := r.PathPrefix("/events").Subrouter()
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
