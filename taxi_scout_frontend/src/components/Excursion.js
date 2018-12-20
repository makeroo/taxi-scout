import React, {Component} from 'react';
import {Link} from "react-router-dom";
import PickupSummary from "./PickupSummary";
import ExcursionConfiguration from "./ExcursionConfiguration";
import RideType from "./RideType";

class Excursion extends Component {
    render() {
        return (
            <div className="container">
                <h2>Uscita del 15/12</h2>

                <ExcursionConfiguration/>

                <h3>Andata: ore 16.00</h3>

                <RideType/>

                <PickupSummary/>

                <h3>Ritorno: ore 18:30</h3>

                <RideType/>

                <PickupSummary/>
            </div>
        );
    }
}


export default Excursion;
