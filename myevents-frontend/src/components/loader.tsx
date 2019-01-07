import * as React from "react";

export interface LoaderProps {
    loading: boolean;
    message?: string;
}

export class Loader extends React.Component<LoaderProps, {}> {
    render() {
        const message = this.props.message || "Loading. Please wait...";

        if (this.props.loading) {
            return <div>{message}</div>
        }

        return <div>
            {this.props.children}
        </div>;
    }
}