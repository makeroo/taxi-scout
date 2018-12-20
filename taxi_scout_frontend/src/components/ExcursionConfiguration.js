import React, {Component} from 'react';
import {connect} from "react-redux";


const mapStateToProps = (state) => {
    return {

    };
};


const mapDispatchToProps = (dispatch) => {
    return {

    };
};


class ExcursionConfiguration extends Component {
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
                    <tr>
                        <td>Cassandra</td>
                        <td className="text-center">
                            <i className="material-icons">check</i>
                        </td>
                    </tr>
                    <tr>
                        <td>Greta</td>
                        <td className="text-center">
                            <i className="material-icons">clear</i>
                        </td>
                    </tr>
                    <tr>
                        <td>Morgana</td>
                        <td className="text-center">
                            <i className="material-icons">check</i>
                        </td>
                    </tr>
                    <tr>
                        <td>Giuseppina</td>
                        <td className="text-center">
                            <i className="material-icons">check</i>
                        </td>
                    </tr>
                </table>
            </div>
        );
    }
}


export default connect(mapStateToProps, mapDispatchToProps)(ExcursionConfiguration);
