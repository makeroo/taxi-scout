import {
    ACCOUNT_FETCH_DATA_SUCCESS,
    ACCOUNT_HAS_ERRORED,
    ACCOUNT_INFO_FETCH_DATA_SUCCESS,
    ACCOUNT_IS_LOADING, ACCOUNT_UPDATE_ADDRESS, ACCOUNT_UPDATE_NAME
} from "../constants/action-types";

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

        case ACCOUNT_INFO_FETCH_DATA_SUCCESS:
            return {
                ...state,
                data: {
                    ...state.data,
                    groups: action.groups,
                    scouts: action.scouts,
                }
            };

        case ACCOUNT_UPDATE_NAME:
            return {
                ...state,
                data: {
                    ...state.data,
                    name: action.name
                }
            };

        case ACCOUNT_UPDATE_ADDRESS:
            return {
                ...state,
                data: {
                    ...state.data,
                    address: action.address
                }
            };

        default:
            return state;
    }
}
