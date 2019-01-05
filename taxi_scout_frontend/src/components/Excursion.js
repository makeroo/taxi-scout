import React, {Component} from 'react';
import {Link} from "react-router-dom";
import PickupSummary from "./PickupSummary";
import ExcursionConfiguration from "./ExcursionConfiguration";
import RideType from "./RideType";
import {connect} from "react-redux";

const mapStateToProps = (state) => {
    return {
        detail: state.excursion.data.detail,
    };
};


// FIXME: rename Activity o ProgramActivity
class Excursion extends Component {
    render() {
        return (
            <div className="container">
                <h2>Uscita del {this.props.detail.date /*FIXME: format date*/} <Link className="float-right" to="/"><i className="material-icons align-middle">home</i></Link></h2>
                <small style={{position: 'relative', top:'-1em'}}>Luogo di ritrovo: {this.props.detail.location}</small>

                <ExcursionConfiguration/>

                <h3>Andata: ore {this.props.detail.from}</h3>

                <RideType direction="out"/>

                <PickupSummary direction="out"/>

                <h3>Ritorno: ore {this.props.detail.to}</h3>

                <RideType direction="return"/>

                <PickupSummary direction="return"/>
            </div>
        );
    }
}


export default connect(mapStateToProps)(Excursion);
