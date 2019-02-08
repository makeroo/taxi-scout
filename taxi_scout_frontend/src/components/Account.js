import React, {Component} from 'react';
import Modal from 'react-modal';
import {connect} from "react-redux";
import {accountUpdateAddress, accountUpdateName, fetchMyAccount, saveAccount} from "../actions/accounts";
import {BASE_URL} from "../constants/rest_api";
import {SERVER_ERROR, SERVICE_NOT_AVAILABLE} from "../constants/errors";


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

        updateName: (name) => {
            dispatch(accountUpdateName(name))
        },

        updateAddress: (address) => {
            dispatch(accountUpdateAddress(address))
        },

/*        saveAccount: (account) => {
            return dispatch(saveAccount(account))
        },*/
    };
};


const customStyles = {
    content : {
        top                   : '50%',
        left                  : '50%',
        right                 : 'auto',
        bottom                : 'auto',
        marginRight           : '-50%',
        transform             : 'translate(-50%, -50%)'
    }
};


class Account extends Component {
    constructor(props) {
        super(props);

        this.state = {
            saving: false
        };

        this.handleNameChange = this.handleNameChange.bind(this);
        this.handleAddressChange = this.handleAddressChange.bind(this);
        this.handleEditScouts = this.handleEditScouts.bind(this);
        this.handleSave = this.handleSave.bind(this);
    }

    componentDidMount() {
        this.props.fetchMyAccount();
    }

    handleNameChange(evt) {
        this.props.updateName(evt.target.value);
    }

    handleAddressChange(evt) {
        this.props.updateAddress(evt.target.value);
    }

    handleEditScouts(evt) {
        this.props.history.push("/children/");
    }

    handleSave(evt) {
        evt.preventDefault();

        const account = this.props.account.data;

        this.setState({
            saving: true
        });

        fetch(`${BASE_URL}/account/${account.id}`, {
            credentials: 'same-origin',
            method: 'POST',
            body: JSON.stringify(account)
        })
            .then((response) => {
                switch (response.status) {
                    case 502:
                        throw { error: SERVICE_NOT_AVAILABLE };
                    case 200:
                        return {};
                    default:
                        return response.json();
                }
            })
            .then((response) => {
                if ('error' in response)
                    throw response;

                this.setState({saving: false});
                this.props.history.push("/");
            })
            .catch((error) => {
                if (!error || typeof error.error !== 'string') {
                    console.log('ops', error);

                    error = { error: SERVER_ERROR };
                }

                this.setState({
                    saving: false,
                    saveError: error,
                });
            });
/*        this.props.saveAccount(this.props.account.data).then((response) => {
            console.log('saved', response);

            if (response)
                this.props.history.push("/");
        });*/
    }

    // TODO: modal "saving..."

    // TODO: i18n

    render() {
        const account = this.props.account;

        if (account.loading) {
            return <p>Loading...</p>;
        }

        if (account.error) {
            return <p>Error... TODO</p>;
        }

        if (!account.data) {
            return <p>....</p>;
        }

        const data = account.data;
        const savingAction = data.savingAction || {};

        let scoutsSummary = '';

        const scouts = data.scouts;

        if (scouts) {
            scoutsSummary = scouts.map(function (s) {return s.name}).join()
        }

        return (
            <div className="Account">
                <div className="container">
                    <h1>Something about you</h1>
                    <form>
                        <div className="form-group">
                            <label htmlFor="email">Email</label>
                            <input type="text"
                                   readOnly
                                   className="form-control-plaintext"
                                   id="email"
                                   value={data.email}
                            />
                        </div>
                        <div className="form-group">
                            <label htmlFor="fullName">Name</label>
                            <input type="text"
                                   className="form-control"
                                   id="fullName"
                                   placeholder="Nome e cognome"
                                   aria-describedby="fullNameHelpBlock"
                                   value={data.name}
                                   onChange={this.handleNameChange}
                            />
                            <small id="fullNameHelpBlock" className="form-text text-muted">
                                Qualcosa che aiuti gli altri ad identificarti: nome e cognome,
                                padre/madre/nonno/zio/angelo custode/tutore di..., anche il
                                nick nel gruppo whatsapp può andare, purché sia pronunciabile ;-)
                            </small>
                        </div>
                        <div className="form-group">
                            <label htmlFor="address">Address</label>
                            <input type="text"
                                   className="form-control"
                                   id="address"
                                   placeholder="Indirizzo approssimativo"
                                   aria-describedby="addressHelpBlock"
                                   value={data.address}
                                   onChange={this.handleAddressChange}
                            />
                            <small id="addressHelpBlock" className="form-text text-muted">
                                È sufficiente una indicazione vaga: comune, frazione, quartiere.
                                Serve a capire come organizzare al meglio i gruppi di trasporto.
                            </small>
                        </div>

                        <div className="form-group">
                            <label htmlFor="children">Scouts</label>
                            <div className="input-group mb-2" onClick={this.handleEditScouts}>
                                <input type="text"
                                       readOnly
                                       className="form-control"
                                       id="children"
                                       placeholder="None, tap to insert someone..."
                                       value={scoutsSummary}
                                       area-describedby="childrenHelpBlock"
                                />
                                <div className="input-group-append">
                                    <div className="input-group-text"><i className="material-icons">edit</i></div>
                                </div>
                            </div>
                            <small id="childrenHelpBlock" className="form-text text-muted">
                                I bimbi che ogni settimana devi scarrozzare a giro per il mondo.
                            </small>
                        </div>

                        <button type="submit"
                                className="btn btn-primary mx-2"
                                onClick={this.handleSave}
                                disabled={savingAction.inProgress}
                        >Update and return to homepage</button>

                    </form>
                </div>
                {/*
                                    overlayClassName="modal fade show"
                    bodyOpenClassName="modal-open"
                    className="modal-dialog modal-dialog-centered"

                */}
                <Modal
                    style={customStyles}
                    isOpen={this.state.saving}
                    contentLabel="Saving..."
                >
                    <h2 >Hello</h2>
                    <div>I am a modal</div>
                    <form>
                        <input />
                        <button>tab navigation</button>
                        <button>stays</button>
                        <button>inside</button>
                        <button>the modal</button>
                    </form>
                    Wait please
                </Modal>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(Account);
