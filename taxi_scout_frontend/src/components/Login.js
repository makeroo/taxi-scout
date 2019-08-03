import React, {Component} from 'react';
import "./Login.scss";
import {connect} from "react-redux";
import {signIn} from "../actions/accounts";


const mapStateToProps = (state) => {
    return {

    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        signIn: (email, password) => {
            dispatch(signIn(email, password))
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

    handleSignIn() {
        this.props.signIn(this.emailInput.current.value, this.paswordInput.current.value)
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
                                type="button"
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
