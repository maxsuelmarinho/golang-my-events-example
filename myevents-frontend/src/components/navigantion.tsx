import * as React from "react";
import {Link} from "react-router-dom";

export interface NavigationProps {
    brandName: string;
}

export class Navigation extends React.Component<NavigationProps, {}> {
    render() {
        return <nav className="navbar navbar-default">
            <div className="container">
                <Link to="/" className="navbar-brand">
                    {this.props.brandName}
                </Link>
            </div>

            <ul className="nav navbar-nav">
                <li><Link to="/">Events</Link></li>
            </ul>
        </nav>
    }
}