import React, {Component} from 'react';
import {connect} from "react-redux";
import classNames from "classnames";
import {setRideRole} from "../actions/set_ride_role";
import {setRides} from "../actions/set_rides";
import "./RideType.scss";

const mapStateToProps = (state, props) => {
    const myId = state.account.data.id;

    return {
        myId,
        tutor_desc: state.excursion.data.tutors[myId],
        tutor_dir: state.excursion[props.direction].tutors[myId],
    };
};


const mapDispatchToProps = (dispatch) => {
    return {
        setRideRole: (tutorId, role, direction) => dispatch(
            setRideRole(tutorId, role, direction)
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
        this.props.setRideRole(this.props.myId, role, this.props.direction);
    }

    handleRideMutex (evt) {
        this.props.setRides(this.props.myId, evt.target.checked ? 1 : 2);
    }

    render() {
        const taxiClass = classNames('btn', 'ridetype-selector', {
            'btn-primary': this.props.tutor_dir.role === 'F',
            'btn-outline-primary': this.props.tutor_dir.role !== 'F',
        });
        const riderClass = classNames('btn', 'ridetype-selector', {
            'btn-primary': this.props.tutor_dir.role === 'R',
            'btn-outline-primary': this.props.tutor_dir.role !== 'R',
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
                               checked={this.props.tutor_desc.rides === 1}
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
