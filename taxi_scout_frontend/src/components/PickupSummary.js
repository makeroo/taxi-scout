import React, {Component} from 'react';
import {connect} from "react-redux";
import "./PickupSummary.scss";


const mapStateToProps = (state) => {
    return {

    };
};


const mapDispatchToProps = (dispatch) => {
    return {

    };
};


class PickupSummary extends Component {
    constructor(props) {
        super(props);

//        this.handleEditScouts = this.handleEditScouts.bind(this);
    }

/*    handleEditScouts(evt) {
        this.props.history.push("/children/");
    }
*/
    render() {
        return (
            <div className="row">
                <table className="table">
                    <tbody>
                    <tr>
                        <th>Serena</th>
                        <td>
                            <small>taxista</small>
                            <br/>
                            <small>s. donato</small>

                        </td>
                        <td><i className="material-icons">announcement</i></td>
                    </tr>
                    <tr>
                        <td>Anita</td>
                        <td><i className="material-icons">check</i></td>
                        <td>&nbsp;</td>
                    </tr>
                    <tr>
                        <th>Sonia</th>
                        <td>
                            <small>cerca</small>
                            <br/>
                            <small>ponte</small>
                        </td>
                        <td><i className="material-icons">question_answer</i></td> {/* forum, chat_bubble(_outline) */}
                    </tr>
                    <tr>
                        <td>Giuliano</td>
                        <td><i className="material-icons">clear</i></td>
                        <td>&nbsp;</td>
                    </tr>
                    <tr>
                        <th>Ilenia</th>
                        <td>
                            <small>cerca</small>
                            <br/>
                            <small>ponte</small>
                        </td>
                        <td><i className="material-icons">chat_bubble</i></td> {/* forum, chat_bubble(_outline) */}
                    </tr>
                    <tr>
                        <td>Andrea</td>
                        <td><i className="material-icons">check</i></td>
                        <td>
                            <small>sonia</small>
                            <small>esselunga</small>
                            <small>15:45</small>
                        </td>
                    </tr>
                    <tr>
                        <th>io</th>
                        <td colSpan="2">
                            <form className="form-inline">
                                <div className="form-group-sm">
                                    <label htmlFor="children">Posti liberi</label>
                                    <input type="number"
                                           className="form-control free-seats-input"
                                           id="children"
                                           placeholder="Nessuno inserito"
                                           value="3"
                                           area-describedby="childrenHelpBlock"
                                    />
                                </div>
                            </form>
                        </td>
                    </tr>
                    <tr>
                        <td>Cassandra</td>
                        <td><i className="material-icons">clear</i></td>
                        <td>&nbsp;</td>
                    </tr>
                    <tr>
                        <td>Greta</td>
                        <td><i className="material-icons">clear</i></td>
                        <td>&nbsp;</td>
                    </tr>
                    <tr>
                        <td>Giuseppina</td>
                        <td><i className="material-icons">clear</i></td>
                        <td>&nbsp;</td>
                    </tr>
                    </tbody>
                </table>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(PickupSummary);
