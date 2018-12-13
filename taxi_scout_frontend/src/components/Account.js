import React, {Component} from 'react';
import {withRouter} from "react-router-dom";
import {connect} from "react-redux";


const mapStateToProps = (state) => {
    return {

    };
};


const mapDispatchToProps = (dispatch) => {
    return {

    };
};


class Account extends Component {
    constructor(props) {
        super(props);

        this.handleEditScouts = this.handleEditScouts.bind(this);
        this.handleSave = this.handleSave.bind(this);
    }

    handleEditScouts(evt) {

    }

    handleSave(evt) {
        evt.preventDefault();

        this.props.history.push("/");
    }

    render() {
        return (
            <div className="Account">
                <div className="container">
                    <h1>Informazioni personali</h1>
                    <form>
                        <div className="form-group">
                            <label htmlFor="email">Email</label>
                            <input type="text"
                                   readOnly
                                   className="form-control-plaintext"
                                   id="email"
                                   value="email@example.com"/>
                                {/*<input type="email"
                                   className="form-control"
                                   id="email"
                                   placeholder="Inserisci l'email"
                            />*/}
                            {/*<button type="submit" className="btn btn-primary mb-2">Confirm identity</button>*/}
                        </div>
                        <div className="form-group">
                            <label htmlFor="fullName">Nome</label>
                            <input type="text"
                                   className="form-control"
                                   id="fullName"
                                   placeholder="Nome e cognome"
                                   aria-describedby="fullNameHelpBlock"
                            />
                            <small id="fullNameHelpBlock" className="form-text text-muted">
                                Qualcosa che aiuti gli altri ad identificarti: nome e cognome,
                                padre/madre/nonno/zio/angelo custode/tutore di..., anche il
                                nick nel gruppo whatsapp può andare, purché sia pronunciabile ;-)
                            </small>
                        </div>
                        <div className="form-group">
                            <label htmlFor="address">Indirizzo</label>
                            <input type="text"
                                   className="form-control"
                                   id="address"
                                   placeholder="Indirizzo approssimativo"
                                   aria-describedby="addressHelpBlock"
                            />
                            <small id="addressHelpBlock" className="form-text text-muted">
                                È sufficiente una indicazione vaga: comune, frazione, quartiere.
                                Serve a capire come organizzare al meglio i gruppi di trasporto.
                            </small>
                        </div>

                        <button type="button"
                                className="btn btn-info mx-2"
                                onClick={this.handleEditScouts}
                        >Scoutisti</button>
                        <button type="submit"
                                className="btn btn-primary mx-2"
                                onClick={this.handleSave}
                        >Aggiorna</button>

                        {/*
                        <div className="form-group">
                            <label htmlFor="password">Password</label>
                            <input type="password"
                                   className="form-control"
                                   id="password"
                                   placeholder="Password"/>
                        </div>
                        */}
                    </form>
                </div>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(Account);
