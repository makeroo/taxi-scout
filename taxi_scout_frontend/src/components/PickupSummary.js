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
                            <br/>
                            <small>esselunga</small>
                            <small className="mx-2">15:45</small>
                        </td>
                    </tr>
                    <tr>
                        <th>io</th>
                        <td colSpan="2">
                            <label htmlFor="children">Posti liberi</label>
                            <div className="input-group input-group-sm" style={{width:'65%'}}>
                                <div className="input-group-prepend">
                                    <span className="input-group-text" id="inputGroupPrepend"><i className="material-icons">add</i></span>
                                </div>
                                <input type="text" className="form-control"
                                       id="children"
                                       value="2"
                                       aria-describedby="inputGroupPrepend"
                                       required/>
                                <div className="input-group-append">
                                    <span className="input-group-text" id="inputGroupPrepend"><i className="material-icons">remove</i></span>
                                </div>
                            </div>
                        </td>
                    </tr>
                    <tr>
                        <td colSpan="3">
                            <table className="w-100">
                                <caption style={{'caption-side':'top'}}>Riassunto incontri</caption>
                                <tbody>
                                <tr>
                                    <td>sonia</td>
                                    <td>esselunga</td>
                                    <td>15:45</td>
                                </tr>
                                </tbody>
                            </table>
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
