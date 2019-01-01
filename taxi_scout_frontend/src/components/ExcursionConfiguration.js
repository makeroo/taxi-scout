import React, {Component} from 'react';
import {connect} from "react-redux";
import { map } from "lodash";
import { setScoutParticipate } from "../actions/set_scout_participate";


function myScouts(scoutIndex, myScoutIds) {
    let scouts = map(myScoutIds, scoutId => scoutIndex[scoutId]);

    function c (s1, s2) {
        return s1.name.localeCompare(s2.name);
    }

    scouts.sort(c);

    // chrome messes up with lambda version:
    // scouts.sort((s1, s2) => s1.name.localCompare(s2.name));
    // even if this works instead:
    // scouts.sort((s1, s2) => c(s1, s2));

    return scouts;
}


const mapStateToProps = (state) => {
    return {
        // TODO: memoize!
        myScouts: myScouts(
            state.excursion.data.scouts,
            state.excursion.out.tutors[state.account.data.id].scouts,
        ),
    };
};


const mapDispatchToProps = (dispatch) => {
    return {
        setScoutParticipate: action => dispatch(action),
    };
};


class ExcursionConfiguration extends Component {
    /*constructor(props) {
        super(props);

        this.togglePartecipate = this.togglePartecipate.bind(this);
    }*/

    toggleParticipate(scout) {
        // why not call directly setScout in onClick expression?
        this.props.setScoutParticipate(
            setScoutParticipate(scout.id, !scout.participate)
        );
    }

    render() {
        return (
            <div className="row">
                <table className="table">
                    <tbody>
                    {map(this.props.myScouts, (scout) => (
                        <tr key={scout.id} onClick={ evt => this.toggleParticipate(scout)}>
                            <td>{scout.name}</td>
                            <td className="text-center">
                                <i className="material-icons">{scout.participate ? 'check' : 'clear'}</i>
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(ExcursionConfiguration);
