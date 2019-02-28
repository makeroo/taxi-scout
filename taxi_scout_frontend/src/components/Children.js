import React, {Component} from 'react';
import {fetchMyAccount, editScout, scoutUpdateName} from "../actions/accounts";
import {connect} from "react-redux";


const mapStateToProps = (state) => {
    return {
        account: state.account,
    }
};


const mapDispatchToProps = (dispatch) => {
    return {
        fetchMyAccount: () => {
            dispatch(fetchMyAccount())
        },

        editScout: (account, index) => {
            dispatch(editScout(account, index))
        },

        editScoutName: (newName) => {
            dispatch(scoutUpdateName(newName))
        },
    };
};


class Children extends Component {
    constructor(props) {
        super(props);

/*        this.state = {
            editingLineIdx: -1,
        };
*/
        this.handleScoutNameChange = this.handleScoutNameChange.bind(this);
        this.handleAddScout = this.handleAddScout.bind(this);
        this.handleBack = this.handleBack.bind(this);
    }

    componentDidMount() {
        // TODO: check if account is already available or reload it anyway?
        this.props.fetchMyAccount();
    }

    handleBack() {
        // TODO: save if needed
        this.props.history.push("/account/");
    }

    handleScoutNameChange(evt) {
        this.props.editScoutName(evt.target.value);
    }

    handleEditScout(index) {
        this.props.editScout(this.props.account, index);
    }

    handleRemoveScout(index) {
        // TODO: modal confirm
        console.log('remove scout', index);
    }

    handleAddScout() {
        this.props.editScout(this.props.account, -1);
/*        this.saveIfNeeded()
            .then(() => {
                const scouts = this.props.account.scouts;

                scouts.push({
                    name: '',
                });

                this.setState({
                    ...this.state,
                    editingLineIdx: scouts.length - 1
                })
            })
            .catch((error) => {
                console.log('todo', error);
            });*/
    }
/*
    saveIfNeeded() {
        if (this.state.editingLineIdx < 0)
            return Promise.resolve();

        const scouts = this.props.account.scouts;
        const scout = scouts[this.state.editingLineIdx];

        if (scout.name === this.state.name)
            return Promise.resolve();

        return scout.id === undefined ?
            jsonFetch('/scouts', scout)
            :
            jsonFetch(`/scout/${scout.id}`, scout, 'PUT');
    }
*/
    render() {
        const account = this.props.account;

        if (account.loading) {
            return <p>Loading...</p>;
        }

        if (account.error) {
            return <p>Error... TODO</p>;
        }

        if (!account.data) {
            return <p>...</p>;
        }

        let scouts = account.scouts || [];
        let groups = account.groups || [];
        let scoutEditing = account.scoutEditing || {};

        return (
            <div className="container">
                <h2>Scouts you take care of</h2>
                <div className="row mt-4">
                    <table className="table">
                        <tbody>
                        {scouts.map((scout, index) => (
                            index === scoutEditing.index ? (
                                <tr key={scout.id || -1}>
                                    <td><input type="text"
                                               className="form-control form-control-sm"
                                               value={scout.name}
                                               onChange={evt => this.handleScoutNameChange(evt, scout)}
                                    /></td>
                                    <td className="text-center">
                                        <i className="material-icons mx-2">edit</i>
                                        <i className="material-icons">remove</i>
                                    </td>
                                </tr>
                            ) : (
                                <tr key={scout.id}>
                                    <td>{scout.name}</td>
                                    <td className="text-center">
                                        <i className="material-icons mx-2" onClick={evt => this.handleEditScout(index)}>edit</i>
                                        <i className="material-icons" onClick={evt => this.handleRemoveScout(index)}>remove</i>
                                    </td>
                                </tr>
                            )
                        ))}
                        </tbody>
                    </table>
                </div>
                { groups.length &&
                <div className="row justify-content-end">
                    <div className="col-2">
                        <i className="material-icons" onClick={this.handleAddScout}>add</i>
                    </div>
                </div>
                }
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


export default connect(mapStateToProps, mapDispatchToProps)(Children);
