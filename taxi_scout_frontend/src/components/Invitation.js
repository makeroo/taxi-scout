import React, {Component} from 'react';
import {connect} from "react-redux";
import {Link} from "react-router-dom";
import { checkToken } from "../actions/invitations";
import {EXPIRED, NOT_FOUND, SERVICE_NOT_AVAILABLE} from "../constants/errors";


const mapStateToProps = (state /*, props*/) => {
    return {
        invitation: state.invitation,
        account: state.account,
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
        const invitation = this.props.invitation;

        if (invitation.loading) {
            return <p>Loading...</p>;
        }

        const error = invitation.error;

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
                        <p>If you completed registration then try to <Link to="/login/">login</Link>.</p>
                        <p>Otherwise please ask your scout group coordinator to send you another invitation.</p>
                    </div>
                );
            }

            if (error.error === EXPIRED) {
                return (
                    <div>
                        <p>Your invitation expired.</p>
                        <p>Please contact your scout group coordinator to receive another one.</p>
                        <p>If you have a valid account then try to <Link to="/login/">login</Link>.</p>
                    </div>
                );
            }

            return (
                <div>
                    <p>Unexpected error.</p>
                    <p>Something went wrong. Sometimes just reloading the page resolves the issue. If the problem
                        persists then please contact system administrator.</p>
                </div>
            );
        }

        if (invitation.data === null) {
            return <p>Please, wait...</p>;
        }

        if (invitation.data.authenticated) {
            if (invitation.data.new_account) {
                return (
                    <div>
                        <p>Wellcome to Taxi Scout!</p>
                        <p>You successfully received an invitation and an account has been created.</p>
                        <p>Next steps:
                            <ul>
                                <li>complete your profile: visit the <Link to="/account/">account page</Link></li>
                                <li>start coordinating with the other users in the group for the next program activity.</li>
                            </ul>
                        </p>
                    </div>
                );
            } else if (invitation.data.scout_group) {
                return (
                    <div>
                        <p>Congratulations! You successfully joined another group.</p>
                        <p><Link to="/">Back to the homepage</Link></p>
                    </div>
                );
            } else {
                return (
                    <div>
                        <p>You successfully signed in by verifying your email address.</p>
                        <p>Please, <a href="/change-password/">change your password</a> now.</p>
                    </div>
                );
            }
        } else {
            return (
                <div>
                    <p>The invitation has been already processed. Please go to the <Link to="/">homepage</Link> and update your bookmarks.</p>
                </div>
            );
        }
    }
}

export default connect(mapStateToProps, mapDispatchToProps)(Invitation);
