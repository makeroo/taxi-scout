import React, {Component} from 'react';
import {BrowserRouter as Router, Route} from "react-router-dom";
import Home from "./Home";
import Invitation from "./Invitation";
import Account from "./Account";
import Excursion from "./Excursion";
import Children from "./Children";
import Login from "./Login";
import ForgotPassword from "./ForgotPassword";
import './App.scss';


class App extends Component {
    render() {
        return (
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
                </div>
            </Router>
        );
    }
}

export default App;
