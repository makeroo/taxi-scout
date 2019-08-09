import React, {Component} from 'react';
import "./Login.scss";
import {connect} from "react-redux";
import {signIn} from "../actions/accounts";
import {toastr} from 'react-redux-toastr';
import {fetchMyAccount} from "../actions/accounts";
import {NOT_AUTHORIZED, SERVICE_NOT_AVAILABLE, SERVER_ERROR} from "../constants/errors";


const mapStateToProps = (state) => {
    return {
        account: state.account,
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        fetchMyAccount: () => {
            dispatch(fetchMyAccount())
        },

        signIn: (email, password) => {
            return dispatch(signIn(email, password));
        },
    };
};


class Login extends Component {
    constructor(props) {
        super(props);

        this.handleSignIn = this.handleSignIn.bind(this);
        this.handleForgotPassword = this.handleForgotPassword.bind(this);
        this.emailInput = React.createRef();
        this.paswordInput = React.createRef();
    }

    componentDidMount() {
        this.props.fetchMyAccount();
    }

    handleSignIn(evt) {
        evt.preventDefault();

        let me = this;

        me.props.signIn(me.emailInput.current.value, me.paswordInput.current.value).then (function () {
            me.props.history.push('/');

        }).catch (function (error) {
            // TODO: move error messages to some error indexed map elsewhere
            if (error.error === NOT_AUTHORIZED) {
                toastr.error('Authentication failed', 'Check password');
            } else if (error.error === SERVICE_NOT_AVAILABLE) {
                toastr.error('Service failure', 'Retry later');
            } else if (error.error === SERVER_ERROR) {
                toastr.error('Service failure', 'Please contact site admins');
            } else {
                toastr.error('Application error', 'This is probably a bug. Please contact the developers.');
            }
        });
    }

    handleForgotPassword() {
        this.props.history.push('/forgot-password')
    }

    render() {
        const account = this.props.account;

        if (account.data) {
            this.props.history.push("/");
            return null;
        }

        // note: signin errors are handled in handleSignIn method

        return (
            <div className="Login">
                <div className="container">
                    <h1>Taxi Scout!</h1>
                    <form className="form-signin">
                        <h2>Please sign in</h2>
                        <label htmlFor="signin_email" className="sr-only">Email</label>
                        <input type="email" id="signin_email"
                               ref={this.emailInput}
                               placeholder="Email address"
                               className="form-control"
                        />
                        <label htmlFor="signin_password" className="sr-only">Password</label>
                        <input type="password" id="signin_password"
                               ref={this.paswordInput}
                               placeholder="Password"
                               className="form-control"
                        />
                        <button className="btn btn-lg btn-primary btn-block mb-3"
                                type="submit"
                                onClick={this.handleSignIn}
                        >Sign in</button>
                        <button className="btn btn-lg btn-link btn-block"
                                type="button"
                                onClick={this.handleForgotPassword}
                        >Forgot password?</button>
                    </form>
                </div>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(Login);
