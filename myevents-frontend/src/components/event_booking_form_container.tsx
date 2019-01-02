import * as React from "react";
import {EventBookingForm} from "./event_booking_form";
import {Event} from "../models/event";

export interface EventBookingFormContainerProps {
    eventID: string;
    eventServiceURL: string;
    bookingServiceURL: string;
}

export interface EventBookingFormContainerState {
    state: "loading"|"ready"|"saving"|"done"|"error";
    event?: Event
}

export class EventBookingFormContainer extends React.Component<EventBookingFormContainerProps, EventBookingFormContainerState> {
    constructor(p: EventBookingFormContainerProps) {
        super(p);

        this.state = {state: "loading"};

        fetch(p.eventServiceURL + "/events/" + p.eventID)
            .then<Event>(response => response.json())
            .then(event => {
                this.setState({
                    state: "ready",
                    event: event
                })
            })
    }

    render() {
        if (this.state.state === "loading") {
            return <div>Loading...</div>;
        }

        if (this.state.state === "saving") {
            return <div>Saving...</div>;
        }

        if (this.state.state === "done") {
            return <div className="alert alert-success">Booking completed! Thank you!</div>;
        }

        if (this.state.state === "error" || !this.state.event) {
            return <div className="alert alert-danger">Unknown error!</div>;
        }

        return <EventBookingForm event={this.state.event} onSubmit={seats => this.handleSubmit(seats)} />
    }

    private handleSubmit(seats: number) {
        const url = this.props.bookingServiceURL + "/events/" + this.props.eventID + "/bookings";
        const payload = {seats: seats};

        this.setState({
            event: this.state.event,
            state: "saving"
        });

        fetch(url, {method: "POST", body: JSON.stringify(payload)})
            .then(response => {
                this.setState({
                    event: this.state.event,
                    state: response.ok ? "done" : "error"
                });
            })
    }
}