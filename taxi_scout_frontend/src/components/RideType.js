import React, {Component} from 'react';
import {connect} from "react-redux";
import classNames from "classnames";
import {setRideRole} from "../actions/set_ride_role";
import {setRides} from "../actions/set_rides";
import "./RideType.scss";

const mapStateToProps = (state, props) => {
    return {
        tutor: state.excursion[props.direction].tutors[state.account.data.id],
    };
};


const mapDispatchToProps = (dispatch) => {
    return {
        setRideRole: (tutorId, role) => dispatch(
            setRideRole(tutorId, role, ['out', 'return'])
        ),
        setRides: (tutorId, rides) => dispatch(
            setRides(tutorId, rides)
        )
    };
};

// FIXME: rename -> RideRole
class RideType extends Component {
    constructor(props) {
        super(props);

        this.handleRideMutex = this.handleRideMutex.bind(this);
    }

    selectRole (role) {
        this.props.setRideRole(this.props.tutor.id, role);
    }

    handleRideMutex (evt) {
        this.props.setRides(this.props.tutor.id, evt.target.value === 'on' ? 1 : 2);
    }

    render() {
        const taxiClass = classNames('btn', 'ridetype-selector', {
            'btn-primary': this.props.tutor.role === 'F',
            'btn-outline-primary': this.props.tutor.role !== 'F',
        });
        const riderClass = classNames('btn', 'ridetype-selector', {
            'btn-primary': this.props.tutor.role === 'R',
            'btn-outline-primary': this.props.tutor.role !== 'R',
        });

        return (
            <div className="row mb-1">
                <div className="col">
                    <button type="button"
                            className={taxiClass}
                            onClick={ evt => this.selectRole('F') }
                    >Ho posti liberi</button>
                </div>
                <div className="col">
                    <button type="button"
                            className={riderClass}
                            onClick={ evt => this.selectRole('R') }
                    >Cerco taxista</button>
                </div>
                <div className="col">
                    <div className="form-group form-check">
                        <input type="checkbox" className="form-check-input"
                               id={'riderMutex' + this.props.direction}
                               checked={this.props.tutor.rides === 1}
                               onChange={this.handleRideMutex}
                        />
                        <label className="form-check-label"
                               htmlFor={'riderMutex' + this.props.direction}
                        >Vado o torno</label>
                    </div>
                </div>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(RideType);
