import React, {Component} from 'react';
import {connect} from "react-redux";
import { checkToken } from "../actions/invitations";
import {EXPIRED, NOT_FOUND, SERVICE_NOT_AVAILABLE} from "../constants/errors";


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
        console.log('invitation', this.props.invitation);

        if (this.props.invitation.loading) {
            return <p>Loading...</p>;
        }

        const error = this.props.invitation.error;

        if (error) {
            if (error.error === SERVICE_NOT_AVAILABLE) {
                return (
                    <div>
                        <p>Service not available.</p>
                        <p>Retry later.</p>
                    </div>
                );
            }

            if (error.error === NOT_FOUND) {
                return (
                    <div>
                        <p>Invitation not found.</p>
                        <p>If you completed registration and verified the email then try to login.</p>
                        <p>Otherwise please ask your scout group coordinator to send you another invitation.</p>
                    </div>
                );
            }

            if (error.error === EXPIRED) {
                return (
                    <div>
                        <p>Your invitation expired.</p>
                        <p>
                            Please contact your scout group coordinator to receive another one.
                        </p>
                    </div>
                );
            }
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
