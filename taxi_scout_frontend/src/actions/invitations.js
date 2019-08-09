import {INVITATION_FETCH_DATA_SUCCESS, INVITATION_HAS_ERRORED, INVITATION_IS_LOADING} from "../constants/action-types";
import {jsonFetch, parseError} from "../utils/json_fetch";


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

        return jsonFetch('/accounts', { invitation: token })
            .then((invitation) => {
                //console.log('invitation resp', invitation);

                dispatch(invitationFetchDataSuccess(invitation));

                return invitation;
            })
            .catch((error) => {
                dispatch(invitationHasErrored(parseError(error)));
            });
    };
}


export function sendInvitation(email) {
    return (dispatch) => {
        dispatch(invitationIsLoading());

        return jsonFetch('/invitations', { email }, 'POST')
            .then((res) => {
                dispatch(invitationFetchDataSuccess(res))
            })
            .catch((error) => {
                dispatch(invitationHasErrored(parseError(error)));
                return Promise.reject(error);
            });
    }
}
