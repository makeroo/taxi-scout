import {INVITATION_FETCH_DATA_SUCCESS, INVITATION_HAS_ERRORED, INVITATION_IS_LOADING} from "../constants/action-types";

export function invitation(
    state = {
        error: null,
        loading: false,
        data: null,
    },
    action
) {
    switch (action.type) {
        case INVITATION_IS_LOADING:
            return {
                error: null,
                loading: true,
                data: null,
            };

        case INVITATION_HAS_ERRORED:
            return {
                error: action.error,
                loading: false,
                data: null,
            };

        case INVITATION_FETCH_DATA_SUCCESS:
            return {
                error: null,
                loading: false,
                data: action.invitation,
            };

        default:
            return state;
    }
}
