import React, {Component} from 'react';
import {connect} from "react-redux";
import "./RideType.scss";

const mapStateToProps = (state) => {
    return {

    };
};


const mapDispatchToProps = (dispatch) => {
    return {

    };
};

// FIXME: rename -> RideRole
class RideType extends Component {
    /*constructor(props) {
        super(props);

//        this.handleEditScouts = this.handleEditScouts.bind(this);
    }*/

    /*    handleEditScouts(evt) {
            this.props.history.push("/children/");
        }
    */
    render() {
        return (
            <div className="row mb-1">
                <div className="col">
                    <button type="button"
                            className="btn btn-primary ridetype-selector"
                    >Ho posti liberi</button>
                </div>
                <div className="col">
                    <button type="button"
                            className="btn btn-outline-primary ridetype-selector"
                    >Cerco taxista</button>
                </div>
                <div className="col">
                    <div className="form-group form-check">
                        <input type="checkbox" className="form-check-input" id="exampleCheck1"/>
                        <label className="form-check-label" htmlFor="exampleCheck1">Uno o l'altro</label>
                    </div>
                </div>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(RideType);
