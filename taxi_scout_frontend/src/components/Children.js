import React, {Component} from 'react';


class Children extends Component {
    constructor(props) {
        super(props);

        this.handleBack = this.handleBack.bind(this);
    }

    handleBack() {
        this.props.history.push("/account/");
    }

    render() {
        return (
            <div className="container">
                <h2>I tuoi scout</h2>
                <div className="row mt-4">
                    <table className="table">
                        <tr>
                            <td>Cassandra</td>
                            <td className="text-center">
                                <i className="material-icons mx-2">edit</i>
                                <i className="material-icons">remove</i>
                            </td>
                        </tr>
                        <tr>
                            <td>Greta</td>
                            <td className="text-center">
                                <i className="material-icons mx-2">edit</i>
                                <i className="material-icons">remove</i>
                            </td>
                        </tr>
                        <tr>
                            <td>Morgana</td>
                            <td className="text-center">
                                <i className="material-icons mx-2">edit</i>
                                <i className="material-icons">remove</i>
                            </td>
                        </tr>
                        <tr>
                            <td>Giuseppina</td>
                            <td className="text-center">
                                <i className="material-icons mx-2">edit</i>
                                <i className="material-icons">remove</i>
                            </td>
                        </tr>
                    </table>
                </div>
                <div className="row justify-content-end">
                    <div className="col-2">
                        <i className="material-icons">add</i>
                    </div>
                </div>
                <div className="row">
                    <div className="col">
                        <button type="button"
                                className="btn btn-primary"
                                onClick={this.handleBack}
                        >Indietro</button>
                    </div>
                </div>
            </div>
        );
    }
}

export default Children;
