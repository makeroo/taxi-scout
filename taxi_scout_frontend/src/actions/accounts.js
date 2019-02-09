import {
    ACCOUNT_FETCH_DATA_SUCCESS,
    ACCOUNT_HAS_ERRORED,
    ACCOUNT_INFO_FETCH_DATA_SUCCESS,
    ACCOUNT_IS_LOADING,
    ACCOUNT_IS_SAVING,
    ACCOUNT_SAVE_HAS_ERRORED,
    ACCOUNT_SAVE_SUCCEDED,
    ACCOUNT_UPDATE_ADDRESS,
    ACCOUNT_UPDATE_NAME
} from "../constants/action-types";
import {jsonFetch, parseError} from "../utils/json_fetch";


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
    };
}

export function accountInfoFetchDataSuccess(groups, scouts) {
    return {
        type: ACCOUNT_INFO_FETCH_DATA_SUCCESS,
        groups,
        scouts
    };
}

export function accountUpdateName(name) {
    return {
        type: ACCOUNT_UPDATE_NAME,
        name
    };
}

export function accountUpdateAddress(address) {
    return {
        type: ACCOUNT_UPDATE_ADDRESS,
        address
    };
}

export function accountIsSaving() {
    return {
        type: ACCOUNT_IS_SAVING
    };
}

export function saveAccountHasErrored(error) {
    return {
        type: ACCOUNT_SAVE_HAS_ERRORED,
        error
    };
}

export function saveAccountSucceded() {
    return {
        type: ACCOUNT_SAVE_SUCCEDED
    }
}

export function saveAccount(account) {
    return (dispatch) => {
        dispatch(accountIsSaving());

        return jsonFetch(`/account/${account.id}`, account)
            .then(() => {
                dispatch(saveAccountSucceded());

                // this value will be returned up to the event handler
                // but state will already be updated so there is no need for it actually
                return true;
            })
            .catch((error) => {
                dispatch(saveAccountHasErrored(parseError(error)));

                return false;
            });
    };
}

export function fetchMyAccount() {
    return (dispatch) => {
        dispatch(accountIsLoading());

        return jsonFetch('/account/me')
            .then((account) => {
                dispatch(accountFetchDataSuccess(account));

                return Promise.all([
                    jsonFetch(`/account/${account.id}/groups`),
                    jsonFetch(`/account/${account.id}/scouts`)
                ]);
            })
            .then((responses) => {
                let groups = responses[0];
                let scouts = responses[1];

                dispatch(accountInfoFetchDataSuccess(groups, scouts));
            })
            .catch((error) => {
                dispatch(accountHasErrored(parseError(error)));
            });
    }
}