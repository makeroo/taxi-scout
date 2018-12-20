import React, {Component} from 'react';
import {Link} from "react-router-dom";


class Home extends Component {
    render() {
        return (
            <div className="Home">
                <div className="container">
                    <h1>Taxi Scout!</h1>
                    <small>Storia di un branco di auto muniti che tenta di raccapezzarsi, ogni settimana.</small>
                    <div className="row mt-3">
                        <div className="col">
                            <Link className="btn btn-outline-primary" to="/account/"><i className="material-icons align-middle">account_circle</i> Account</Link>
                        </div>
                        <div className="col-5">
                            <small>Soltanto poche informazioni su di te.</small>
                        </div>
                    </div>
                    <div className="row">
                        <div className="col">
                            Ecco quello che ci serve sapere:
                            <ul className="list-unstyled">
                                <li><i className="material-icons">done</i> come ti chiami: per poter interagire con gli altri;</li>
                                <li><i className="material-icons">done</i> dove vivi: per sapere con chi è più opportuno ritrovarsi;</li>
                                <li><i className="material-icons">clear</i> quanti sono i bimbi che porti a giro: per calcolare le auto.</li>
                            </ul>
                        </div>
                    </div>
                    <div className="row mt-3">
                        <div className="col">
                            <Link className="btn btn-outline-primary" to="/excursion/"><i className="material-icons align-middle">assignment</i> Prossima uscita</Link>
                        </div>
                        <div className="col-5">
                            <small>Al lavoro, c'è da organizzarsi.</small>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}


export default Home;
