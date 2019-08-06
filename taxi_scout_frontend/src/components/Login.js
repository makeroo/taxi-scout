import React, {Component} from 'react';
import "./Login.scss";
import {connect} from "react-redux";
import {signIn} from "../actions/accounts";
import {toastr} from 'react-redux-toastr';
import {NOT_AUTHORIZED} from "../constants/errors";


const mapStateToProps = (state) => {
    return {
        account: state.account,
    };
};

const mapDispatchToProps = (dispatch) => {
    return {
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

    handleSignIn(evt) {
        evt.preventDefault();

        let me = this;

        me.props.signIn(me.emailInput.current.value, me.paswordInput.current.value).then (function () {
            let account = me.props.account;

            if (account.error === null) {
                me.props.history.push('/');
            } else if (account.error.error === NOT_AUTHORIZED) {
                toastr.error('Authentication failed', 'Check password');
            } else {
                // TODO: better error parsing
                toastr.error('Service failure', 'Retry later');
            }
        });
    }

    handleForgotPassword() {
        this.props.history.push('/forgot-password')
    }

    render() {
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
