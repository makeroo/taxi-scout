import React, {Component} from 'react';
import {connect} from "react-redux";
import {Link} from "react-router-dom";
import {jsonFetch} from "../utils/json_fetch";
import { fetchMyAccount } from '../actions/accounts';
import { BAD_REQUEST, NOT_AUTHORIZED, SERVICE_NOT_AVAILABLE } from '../constants/errors';


const mapStateToProps = (state) => {
    return {
        account: state.account,
    }
};


const mapDispatchToProps = (dispatch) => {
    return {
        fetchMyAccount: () => {
            dispatch(fetchMyAccount())
        },
    };
};


class ChangePassword extends Component {
    constructor(props) {
        super(props);

        this.handleBack = this.handleBack.bind(this);
        this.handlePasswordChange = this.handlePasswordChange.bind(this);
        this.handlePasswordCheckChange = this.handlePasswordCheckChange.bind(this);
        this.handleSave = this.handleSave.bind(this);

        this.state = {
            newPassword: '',
            passwordCheck: '',
            passwordMatch: true,
        };
    }

    componentDidMount() {
        this.props.fetchMyAccount();
    }

    handleBack() {
        this.props.history.push("/account/");
    }

    handlePasswordChange(evt) {
        const p = evt.target.value;

        this.setState({
            newPassword: p,
            passwordMatch: p === this.state.passwordCheck,
        });
    }

    handlePasswordCheckChange(evt) {
        const p = evt.target.value;

        this.setState({
            passwordCheck: p,
            passwordMatch: p === this.state.newPassword,
        });
    }

    handleSave(evt) {
        this.setState({saving: true});

        const account = this.props.account;
        const newPassword = this.state.newPassword;

        let me = this;

        jsonFetch(`/account/${account.data.id}/password`, { 'p': newPassword }).then(function (res) {
            me.setState({
                saved: true,
                saving: false,
            });

            me.props.history.push("/account");
        }).catch (function (error) {
            me.setState({
                saveError: error,
            });
        });
    }

    render() {
        const account = this.props.account;

        if (account.loading) {
            return <p>Loading...</p>;
        }

        const error = account.error;

        if (error) {
            if (error.error === BAD_REQUEST || error.error === NOT_AUTHORIZED) {
                return (
                    <div>
                        <p>Session expired (probably)</p>
                        <p>Please <Link to="/login">login again</Link></p>
                    </div>
                );
            }

            if (error.error === SERVICE_NOT_AVAILABLE) {
                return (
                    <div>
                        <p>Service not available.</p>
                        <p>Retry later.</p>
                    </div>
                )
            }

            return (
                <div>
                    <p>Unexpected error.</p>
                    <p>Something went wrong. Sometimes just reloading the page resolves the issue. If the problem
                        persists then please contact system administrator.</p>
                </div>
            );
        }

        if (!account.data) {
            return <p>...</p>;
        }

        const props = this.state;

        return (
            <div className="container">
                <h2>Change your password</h2>
                <form>
                    <div className="form-group">
                        <label htmlFor="newPassword">Enter new password</label>
                        <input type="password"
                                className="form-control"
                                id="newPassword"
                                placeholder="Choose it carefully"
                                aria-describedby="newPasswordHelpBlock"
                                value={props.newPassword}
                                onChange={this.handlePasswordChange}
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="passwordCheck">Repeat it</label>
                        <input type="password"
                                className="form-control"
                                id="passwordCheck"
                                placeholder="Enter the new password again"
                                aria-describedby="passwordCheckHelpBlock"
                                value={props.passwordCheck}
                                onChange={this.handlePasswordCheckChange}
                        />
                    </div>

                    <button type="button"
                            className="btn btn-primary mb-2"
                            disabled={!props.newPassword || !props.passwordMatch}
                            onClick={this.handleSave}
                    >Update and return to homepage</button>

                    <button type="button"
                            className="btn btn-primary"
                            onClick={this.handleBack}
                    >Oh well, maybe I'll keep the old one</button>
                </form>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(ChangePassword);
