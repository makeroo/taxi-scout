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
import {BASE_URL} from "../constants/rest_api";
import {SERVER_ERROR, SERVICE_NOT_AVAILABLE} from "../constants/errors";


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

        return fetch(`${BASE_URL}/account/${account.id}`, {
            credentials: 'same-origin',
            method: 'POST',
            body: JSON.stringify(account)
        })
            .then((response) => {
                switch (response.status) {
                    case 502:
                        throw { error: SERVICE_NOT_AVAILABLE };

                    default:
                        return response.json();
                }
            })
            .then((response) => {
                if ('error' in response)
                    throw response;

                dispatch(saveAccountSucceded());

                // this value will be returned up to the event handler
                // but state will already be updated so there is no need for it actually
                return true
            })
            .catch((error) => {
                if (!error || typeof error.error !== 'string') {
                    console.log('ops', error);

                    error = {
                        error: SERVER_ERROR
                    };
                }

                dispatch(saveAccountHasErrored(error));

                return false
            });
    };
}

export function fetchMyAccount() {
    return (dispatch) => {
        dispatch(accountIsLoading());

        return fetch(BASE_URL + '/account/me', {
            credentials: 'same-origin'
        })
            .then((response) => {
                switch (response.status) {
                    case 502:
                        const error = { error: SERVICE_NOT_AVAILABLE };

                        throw error;

                    default:
                        return response.json();
                }
            })
            .then((account) => {
                if ('error' in account)
                    throw account;

                dispatch(accountFetchDataSuccess(account));

                return Promise.all([
                    fetch(BASE_URL + `/account/${account.id}/groups`, {
                        credentials: 'same-origin'
                    }),
                    fetch(BASE_URL + `/account/${account.id}/scouts`, {
                        credentials: 'same-origin'
                    })
                ]);
            })
            .then((responses) => {
                if (responses[0].status === 502 || responses[1].status === 502) {
                    throw { error : SERVICE_NOT_AVAILABLE };
                }

                return Promise.all([
                    responses[0].json(),
                    responses[1].json()
                ])
            })
            .then((responses) => {
                let groups = responses[0];
                let scouts = responses[1];

                console.log('ricevuto 2', responses);

                if ('error' in groups)
                    throw groups;
                if ('error' in scouts)
                    throw scouts;

                dispatch(accountInfoFetchDataSuccess(groups, scouts));
            })
            .catch((error) => {
                if (!error || typeof error.error !== 'string') {
                    console.log('ops', error);

                    error = {
                        error: SERVER_ERROR
                    };
                }

                dispatch(accountHasErrored(error));
            });
    }
}