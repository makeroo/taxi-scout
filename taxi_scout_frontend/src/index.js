import React from 'react';
import ReactDOM from 'react-dom';
import Modal from 'react-modal';
import './index.scss';
import App from './components/App';
import { Provider } from "react-redux";
import store from "./store/index";
import * as serviceWorker from './serviceWorker';


//console.log(window.location);
Modal.setAppElement('#root');


ReactDOM.render(
    <Provider store={store}>
        <App />
    </Provider>,
    document.getElementById('root')
);


// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: http://bit.ly/CRA-PWA
serviceWorker.unregister();
