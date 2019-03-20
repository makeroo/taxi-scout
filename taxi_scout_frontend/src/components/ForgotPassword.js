import React, {Component} from 'react';
import "./Login.scss";
import {connect} from "react-redux";
import {sendInvitation} from "../actions/invitations";


const mapStateToProps = (state) => {
    return {

    };
};

const mapDispatchToProps = (dispatch) => {
    return {
        sendInvitation: (email, password) => {
            dispatch(sendInvitation(email, password))
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

    handleSendInvitation() {
        this.props.sendInvitation(this.emailInput.current.value);
    }

    handleLogin() {
        this.props.history.push('/login');
    }

    render() {
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
                                onClick={this.handleForgotPassword}
                        >Forgot password?</button>
                    </form>
                </div>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(ForgotPassword);
