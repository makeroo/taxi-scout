import React, {Component} from 'react';
import {connect} from "react-redux";


const mapStateToProps = (state /*, props*/) => {
    return {
        invitation: state.invitation,
    };
};


const mapDispatchToProps = (dispatch) => {
    return {
        checkToken: (token) => {
            dispatch(checkToken(token))
        },
    };
};


class Invitation extends Component {
    componentDidMount() {
        this.props.checkToken(this.props.match.params.token);
    }

    render() {
        if (this.props.invitation.loading) {
            return <p>Loading...</p>;
        }

        if (this.props.invitation.error) {
            return (
                <div>
                    <p>Sorry! There was an error: {this.props.invitation.error}.</p>
                    <p>
                        TODO: try reloading
                        or
                        link to home...
                    </p>
                </div>
            );
        }

        return (
            <div>
                TODO: complete registration form
                or edit profile until email is validated
                or redirect to Home if profile is valid
            </div>
        );
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Invitation);
