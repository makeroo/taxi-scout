/*
copiato dal tutorial redux...
const initialState = {
    // TODO: my properties
};

const rootReducer = (state = initialState, action) => state;

export default rootReducer;
*/



import { combineReducers } from 'redux';
import { excursion } from "./excursion";
import { account } from "./account";
import { scouts } from "./scouts";

export default combineReducers({
    account,
    scouts,
    excursion,
});
