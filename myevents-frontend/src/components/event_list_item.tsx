import {Event} from "../models/event";
import * as React from  "react";
import { Link } from "react-router-dom";

export interface EventListItemProps {
    event: Event;
    selected?: boolean;

    onBooked: () => any;
}

export class EventListItem extends React.Component<EventListItemProps, {}> {
    render() {
        const start = new Date(this.props.event.startDate * 1000);
        const end = new Date(this.props.event.endDate * 1000);

        const locationName = this.props.event.location ? this.props.event.location.name : "unknown";

        console.log(this.props.event);

        return <tr>
                <td>{this.props.event.name}</td>
                <td>{locationName}</td>
                <td>{start.toLocaleDateString()}</td>
                <td>{end.toLocaleDateString()}</td>
                <td>
                    <Link to={`/api/events/${this.props.event.id}/book`} className="btn btn-primary">Book now!</Link>
                </td>
            </tr>
    }
}
