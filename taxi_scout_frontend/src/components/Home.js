import React, {Component} from 'react';
import {Link} from "react-router-dom";


class Home extends Component {
    render() {
        return (
            <div className="Home">
                <div className="container">
                    <div className="row">
                        <Link className="btn btn-outline-primary mx-auto" to="/account/">Account</Link>
                    </div>
                    <div className="row">
                        <Link className="btn btn-outline-primary mx-auto" to="/excursion/">Prossima uscita</Link>
                    </div>
                </div>
            </div>
        );
    }
}


export default Home;
