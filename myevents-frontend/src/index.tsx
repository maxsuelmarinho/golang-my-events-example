import * as React from "react";
import * as ReactDOM from "react-dom";
import {HashRouter as Router, Route} from "react-router-dom";

import {EventListContainer} from "./components/event_list_container";
import {Navigation} from "./components/navigantion";
import { EventBookingFormContainer } from "./components/event_booking_form_container";

class App extends React.Component<{}, {}> {
    render() {
        const eventList = () => <EventListContainer eventListURL="http://localhost:8181" />
        const eventBooking = ({match}: any) =>
            <EventBookingFormContainer eventID={match.params.id} 
                eventServiceURL="http://localhost:8181"
                bookingServiceURL="http://localhost:8282" />;

        return <Router>
            <Navigation brandName="MyEvents" />
            <div className="container">
                <h1>MyEvents</h1>
                <Route exact path="/" component={eventList} />
                <Route exact path="/events/:id/book" component={eventBooking} />
            </div>
        </Router>;
    }
}

ReactDOM.render(
    <App />,
    document.getElementById("myevents-app")
);