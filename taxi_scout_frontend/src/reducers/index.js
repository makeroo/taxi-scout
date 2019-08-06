import { combineReducers } from 'redux';
import { excursion } from "./excursion";
import { account } from "./account";
import { scouts } from "./scouts";
import { invitation } from "./invitation";
import {reducer as toastrReducer} from 'react-redux-toastr';

export default combineReducers({
    account,
    scouts,
    excursion,
    invitation,
    toastr: toastrReducer,
});
