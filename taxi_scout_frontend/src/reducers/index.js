import { combineReducers } from 'redux';
import { excursion } from "./excursion";
import { account } from "./account";
import { scouts } from "./scouts";

export default combineReducers({
    account,
    scouts,
    excursion,
});
