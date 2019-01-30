import {ACCOUNT_FETCH_DATA_SUCCESS, ACCOUNT_HAS_ERRORED, ACCOUNT_IS_LOADING} from "../constants/action-types";

export function account(
    state = {
        error: null,
        loading: false,
        data: null,
    },
    action
) {
    switch (action.type) {
        case ACCOUNT_IS_LOADING:
            return {
                error: null,
                loading: true,
                data: null,
            };

        case ACCOUNT_HAS_ERRORED:
            return {
                error: action.error,
                loading: false,
                data: null,
            };

        case ACCOUNT_FETCH_DATA_SUCCESS:
            return {
                error: null,
                loading: false,
                data: action.account,
            };

        default:
            return state;
    }
}
