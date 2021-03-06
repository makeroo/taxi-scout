import React, {Component} from 'react';
import {BrowserRouter as Router, Route} from "react-router-dom";
import ReduxToastr from 'react-redux-toastr';
import Home from "./Home";
import Invitation from "./Invitation";
import Account from "./Account";
import Excursion from "./Excursion";
import Children from "./Children";
import Login from "./Login";
import ForgotPassword from "./ForgotPassword";
import ChangePassword from "./ChangePassword";
import './App.scss';


class App extends Component {
    render() {
        return (
            <div>
                <Router>
                    <div>
                        {/* <div className="App">
                        //    <header className="App-header">
                        //    </header>
                        </div>
                        <nav>
                            <ul>
                                <li>
                                    <Link to="/">Home</Link>
                                </li>
                                <li>
                                    <Link to="/about/">About</Link>
                                </li>
                                <li>
                                    <Link to="/users/">Users</Link>
                                </li>
                            </ul>
                        </nav>
                        */}

                        <Route path="/" exact component={Home}/>
                        <Route path="/login" component={Login}/>
                        <Route path="/forgot-password" component={ForgotPassword}/>
                        <Route path="/invitation/:token" component={Invitation}/>
                        <Route path="/account/" component={Account}/>
                        <Route path="/excursion/" component={Excursion}/>
                        <Route path="/children/" component={Children}/>
                        <Route path="/change-password/" component={ChangePassword}/>
                    </div>
                </Router>
                <ReduxToastr
                    timeOut={4000}
                    newestOnTop={false}
                    preventDuplicates
                    position="bottom-center"
                    transitionIn="fadeIn"
                    transitionOut="fadeOut"
                    progressBar
                    closeOnToastrClick/>
            </div>
        );
    }
}

export default App;
