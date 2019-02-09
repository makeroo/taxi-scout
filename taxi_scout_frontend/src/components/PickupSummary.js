import React, {Component} from 'react';
import {connect} from "react-redux";
import "./PickupSummary.scss";
import {map, isArray} from "lodash";


function sortedValues(entityIndex, keySubset, sortingProperty) {
    let selectedEntities = map(
        keySubset || Object.keys(entityIndex),
        entityId => entityIndex[entityId]
    );

    function comparator (e1, e2) {
        return e1[sortingProperty].localeCompare(e2[sortingProperty]);
    }

    selectedEntities.sort(comparator);

    // chrome messes up with lambda version:
    // selectedEntities.sort((s1, s2) => s1[sortingProperty].localCompare(s2[sortingProperty]));
    // even if this works instead:
    // selectedEntities.sort((s1, s2) => c(s1, s2));

    return selectedEntities;
}

function matchScout (scoutIndex, scoutIds) {
    if (!isArray(scoutIds))
        return false;

    for (var i = scoutIds.length; --i >= 0; )
        if (scoutIds[i] in scoutIndex)
            return true;

        return false;
}

const mapStateToProps = (state, props) => {
    // TODO: use reselect
    const myId = state.account.data.id;
    const coord = state.excursion[props.direction];
    const tutorsIndex = state.excursion.data.tutors;
    const scoutsIndex = state.excursion.data.scouts;
    let myCoord = null;
    const myScouts = state.scouts.data;
    let myMeetings = [];

    const tutors = sortedValues(tutorsIndex, Object.keys(tutorsIndex), 'name');
    let meetingsByScout = {};

    map(coord.meetings, meeting => {
        map(meeting.riders, scoutId => {
            let meetings = meetingsByScout[scoutId];

            if (meetings === undefined) {
                meetings = [];
                meetingsByScout[scoutId] = meetings;
            }

            meetings.push(meeting);
        });

        if (meeting.taxi === myId || matchScout(myScouts, meeting.scouts))
            myMeetings.push(meeting);
    });

    // TODO: sort scout meetings by time

    const c = function (s1, s2) {
        return s1.name.localeCompare(s2.name);
    };

    let rows = [];

    map(tutors, tutor => {
        const tutorCoord = coord.tutors[tutor.id];

        if (tutor.id === myId)
            myCoord = tutorCoord;

        rows.push({
            type: tutor.id === myId ? 'me' : 'tutor',
            tutor,
            tutorCoord,
        });

        let scouts = [];

        map(tutorCoord.scouts, scoutId => {
            scouts.push(scoutsIndex[scoutId]);
        });

        if (scouts.length) {
            scouts.sort(c);

            map(scouts, scout => {
                rows.push({
                    type: 'scout',
                    scout,
                    meetings: meetingsByScout[scout.id]
                });
            });
        }
    });

    return {
        rows,
        tutorsIndex,
        myCoord,
        myMeetings,
        scoutsIndex,
    };
};


const mapDispatchToProps = (dispatch) => {
    return {
        setFreeSeats: (seats) => {
            console.log('dispatch set free seats', seats, dispatch); // TODO
        },
    };
};


class PickupSummary extends Component {
    constructor(props) {
        super(props);

        this.handleFreeSeatsChange = this.handleFreeSeatsChange.bind(this);
        this.incFreeSeats = this.incFreeSeats.bind(this);
        this.decFreeSeats = this.decFreeSeats.bind(this);
    }

    handleFreeSeatsChange(evt) {
        const x = + evt.target.value;

        if (isNaN(x)) {
            console.log('illegal value:', x, typeof x, evt.target.value)
        } else {
            this.props.setFreeSeats(x);
        }
    }

    incFreeSeats(evt) {
        const x = this.props.myCoord.free_seats + 1;

        this.props.setFreeSeats(x);
    }

    decFreeSeats(evt) {
        const x = this.props.myCoord.free_seats - 1;

        if (x >= 0)
            this.props.setFreeSeats(x);
    }

    render() {
        return (
            <div className="row">
                <table className="table">
                    <tbody>
                    {this.props.rows.map(row => (
                        row.type === 'tutor' ? (
                        <tr key={'t' + row.tutor.id}>
                            <th>{row.tutor.name}</th>
                            <td>
                                <small>{row.tutorCoord.role === 'R' ? 'cerca' : 'taxista'}</small>
                                <br/>
                                <small>{row.tutor.address}</small>
                            </td>
                            <td><i className="material-icons">announcement</i>TODO</td>
                        </tr>
                        ) : row.type === 'me' ? (
                        <tr key={'t' + row.tutor.id}>
                            <th>io</th>
                            <td colSpan="2">
                                <label htmlFor="children">Posti liberi</label>
                                <div className="input-group input-group-sm" style={{width:'65%'}}>
                                    <div className="input-group-prepend"
                                         onClick={this.incFreeSeats}
                                    >
                                        <span className="input-group-text" id="inputGroupPrepend"><i className="material-icons">add</i></span>
                                    </div>
                                    <input type="number" className="form-control"
                                           id="children"
                                           value={row.tutorCoord.free_seats}
                                           onChange={this.handleFreeSeatsChange}
                                           aria-describedby="inputGroupPrepend"
                                           required/>
                                    <div className="input-group-append"
                                         onClick={this.decFreeSeats}
                                    >
                                        <span className="input-group-text" id="inputGroupPrepend"><i className="material-icons">remove</i></span>
                                    </div>
                                </div>
                                <table className="w-100">
                                    <caption style={{captionSide:'top'}}>Riassunto incontri</caption>
                                    <tbody>
                                    {this.props.myMeetings.length === 0 ? (
                                        <tr>
                                            <td colSpan={3}>nessuno</td>
                                        </tr>
                                    ) : (map(this.props.myMeetings, meeting => (
                                            <tr key={meeting.id}>
                                                <td>{map(meeting.riders, (scoutId, idx) => (
                                                    (idx > 0 ? ', ' : '') + (this.props.scoutsIndex[scoutId].name)
                                                ))}</td>
                                                <td>{meeting.point}</td>
                                                <td>{meeting.time}</td>
                                            </tr>
                                        ))
                                    )}
                                    </tbody>
                                </table>
                            </td>
                        </tr>
                        ) : row.type === 'scout' ? (
                        <tr key={'s' + row.scout.id}>
                            <td>{row.scout.name}</td>
                            <td><i className="material-icons">{row.meetings ? 'check' : 'clear'}</i></td>
                            <td>
                                {map(row.meetings, meeting => (
                                    <div key={meeting.id}>
                                        <small>{this.props.tutorsIndex[meeting.taxi].name}</small>
                                        <br/>
                                        <small>{meeting.point}</small>
                                        <small className="mx-2">{meeting.time}</small>
                                    </div>
                                ))}
                            </td>
                        </tr>
                        ) : (
                        <tr>
                        </tr>
                        )
                    ))}

                    {/*<tr>
                        <th>Sonia</th>
                        <td>
                            <small>cerca</small>
                            <br/>
                            <small>ponte</small>
                        </td>
                        <td><i className="material-icons">question_answer</i></td>
                    </tr>
                    <tr>
                        <th>Ilenia</th>
                        <td>
                            <small>cerca</small>
                            <br/>
                            <small>ponte</small>
                        </td>
                        <td><i className="material-icons">chat_bubble</i>{/ * forum, chat_bubble(_outline) * /}</td>
                    </tr>*/}
                    </tbody>
                </table>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(PickupSummary);
