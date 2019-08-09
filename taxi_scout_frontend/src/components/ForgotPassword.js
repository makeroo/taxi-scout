import React, {Component} from 'react';
import "./Login.scss";
import {connect} from "react-redux";
import {sendInvitation} from "../actions/invitations";
import {fetchMyAccount} from "../actions/accounts";
import {toastr} from 'react-redux-toastr';
import {FORBIDDEN, SERVICE_NOT_AVAILABLE, SERVER_ERROR} from "../constants/errors";


const mapStateToProps = (state) => {
    return {
        account: state.account,
        invitation: state.invitation,
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchMyAccount: () => {
            dispatch(fetchMyAccount())
        },

        sendInvitation: (email, password) => {
            return dispatch(sendInvitation(email, password))
        },
    };
};


class ForgotPassword extends Component {
    constructor(props) {
        super(props);

        this.handleSendInvitation = this.handleSendInvitation.bind(this);
        this.handleLogin = this.handleLogin.bind(this);
        this.emailInput = React.createRef();
    }

    componentDidMount() {
        this.props.fetchMyAccount();
    }

    handleSendInvitation(evt) {
        evt.preventDefault();

        this.props.sendInvitation(this.emailInput.current.value).then(function () {
            // console.log('invitation sent', arguments);

        }).catch(function (error) {
            //console.log('invitation error', error, arguments);
            if (error.error === FORBIDDEN) {
                toastr.error('Password reset failed', 'The email was not recognized');
            } else if (error.error === SERVICE_NOT_AVAILABLE) {
                toastr.error('Service failure', 'Retry later');
            } else if (error.error === SERVER_ERROR) {
                toastr.error('Service failure', 'Please contact site admins');
            } else {
                toastr.error('Application error', 'This is probably a bug. Please contact the developers.');
            }
        });
    }

    handleLogin() {
        this.props.history.push('/login');
    }

    render() {
        const account = this.props.account;

        if (account.data) {
            this.props.history.push("/");
            return null;
        }

        const invitation = this.props.invitation;

        if (invitation.data) {
            return (
                <div className="Login">
                    <div className="container">
                        <h1>Taxi Scout!</h1>
                        <h2>Mail sent</h2>
                        <p>
                            You should recive an email containing a link.
                            Please follow it to sign in and then change your
                            password.
                        </p>
                    </div>
                </div>
            );
        }

        return (
            <div className="Login">
                <div className="container">
                    <h1>Taxi Scout!</h1>
                    <form className="form-signin">
                        <h2>Reset password</h2>
                        <p>
                            Enter your email address: you will receive
                            an email, click on the link inside and insert
                            a new password when requested.
                        </p>
                        <label htmlFor="signin_email" className="sr-only">Email</label>
                        <input type="email" id="signin_email"
                               ref={this.emailInput}
                               placeholder="Email address"
                               className="form-control"
                        />

                        <button className="btn btn-lg btn-primary btn-block mb-3"
                                type="submit"
                                onClick={this.handleSendInvitation}
                        >Reset password</button>
                        <button className="btn btn-lg btn-link btn-block"
                                onClick={this.handleLogin}
                        >Sign in</button>
                    </form>
                </div>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(ForgotPassword);
