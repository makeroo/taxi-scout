import {INVITATION_FETCH_DATA_SUCCESS, INVITATION_HAS_ERRORED, INVITATION_IS_LOADING} from "../constants/action-types";
import {BASE_URL} from "../constants/rest_api";


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

        return fetch(BASE_URL + '/invitations/' + token)
            .then((response) => {
                if (!response.ok) {
                    // TODO: parse error
                    throw Error(response.statusText);
                }
                return response.json();
            })
            .then((invitation) => {
                dispatch(invitationFetchDataSuccess(invitation));
                return invitation;
            })
            .catch((error) => {
                // TODO: parse error
                dispatch(invitationHasErrored(error));
            });
    };
}
