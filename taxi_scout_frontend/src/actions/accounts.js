import {ACCOUNT_FETCH_DATA_SUCCESS, ACCOUNT_HAS_ERRORED, ACCOUNT_IS_LOADING} from "../constants/action-types";


export function accountIsLoading() {
    return {
        type: ACCOUNT_IS_LOADING
    };
}


export function accountHasErrored (error) {
    return {
        type: ACCOUNT_HAS_ERRORED,
        error
    };
}

export function accountFetchDataSuccess (account) {
    return {
        type: ACCOUNT_FETCH_DATA_SUCCESS,
        account
    }
}
