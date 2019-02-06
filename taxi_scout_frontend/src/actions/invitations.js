import {INVITATION_FETCH_DATA_SUCCESS, INVITATION_HAS_ERRORED, INVITATION_IS_LOADING} from "../constants/action-types";
import {BASE_URL} from "../constants/rest_api";
import {SERVER_ERROR, SERVICE_NOT_AVAILABLE} from "../constants/errors";


export function invitationIsLoading () {
    return {
        type: INVITATION_IS_LOADING
    };
}


export function invitationHasErrored (error) {
    return {
        type: INVITATION_HAS_ERRORED,
        error
    };
}


export function invitationFetchDataSuccess (invitation) {
    return {
        type: INVITATION_FETCH_DATA_SUCCESS,
        invitation
    }
}


export function checkToken(token) {
    return (dispatch) => {
        dispatch(invitationIsLoading());

        return fetch(BASE_URL + '/accounts', {
            method: 'POST',
            credentials: 'same-origin',
            body: JSON.stringify({ invitation: token })
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
            .then((invitation) => {
                if ('error' in invitation)
                    throw invitation;

                //console.log('invitation resp', invitation);

                dispatch(invitationFetchDataSuccess(invitation));

                return invitation;
            })
            .catch((error) => {
                if (!error || typeof error.error !== 'string') {
                    console.log('ops', error);

                    error = {
                        error: SERVER_ERROR
                    };
                }

                dispatch(invitationHasErrored(error));
            });
    };
}
