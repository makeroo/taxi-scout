import { combineReducers } from 'redux';
import { excursion } from "./excursion";
import { account } from "./account";
import { scouts } from "./scouts";
import { invitation } from "./invitation";

export default combineReducers({
    account,
    scouts,
    excursion,
    invitation,
});
