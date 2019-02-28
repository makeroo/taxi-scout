import {
    ACCOUNT_FETCH_DATA_SUCCESS,
    ACCOUNT_HAS_ERRORED,
    ACCOUNT_INFO_FETCH_DATA_SUCCESS,
    ACCOUNT_IS_LOADING,
    ACCOUNT_IS_SAVING,
    ACCOUNT_SAVE_HAS_ERRORED,
    ACCOUNT_SAVE_SUCCEDED,
    ACCOUNT_UPDATE_ADDRESS,
    ACCOUNT_UPDATE_NAME,
    SCOUT_SAVE_FAILED,
    SCOUT_SAVING,
    SCOUT_SELECT,
    SCOUT_UPDATE_NAME
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
                ...state,
                error: null,
                loading: true,
                data: null,
            };

        case ACCOUNT_HAS_ERRORED:
            return {
                ...state,
                error: action.error,
                loading: false,
                data: null,
            };

        case ACCOUNT_FETCH_DATA_SUCCESS:
            return {
                ...state,
                error: null,
                loading: false,
                data: action.account,
            };

        case ACCOUNT_INFO_FETCH_DATA_SUCCESS:
            return {
                ...state,
                groups: action.groups,
                scouts: action.scouts,
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

        case ACCOUNT_IS_SAVING:
            return {
                ...state,
                savingAction: {
                    inProgress: true,
                }
            };

        case ACCOUNT_SAVE_HAS_ERRORED:
            return {
                ...state,
                savingAction: {
                    inProgress: false,
                    error: action.error
                }
            };

        case ACCOUNT_SAVE_SUCCEDED:
            return {
                ...state,
                savingAction: {
                    inProgress: false,
                }
            };

        case SCOUT_UPDATE_NAME: {
            const scouts = state.scouts.slice(0);
            const scoutEditing = state.scoutEditing;
            const scout = scouts[scoutEditing.index];
            const origName = scoutEditing.origName || scout.name;

            scouts[state.scoutEditing.index] = {
                ...scout,
                name: action.name,
            };

            return {
                ...state,
                scouts,
                scoutEditing: {
                    ...scoutEditing,
                    origName
                }
            };
        }

        case SCOUT_SELECT:
            if (action.index === -1) {
                let scouts = state.scouts || [];

                if (scouts.length === 0) {
                    scouts = scouts.slice(0);
                } else if (scouts[scouts.length - 1].id >= 0) {
                    scouts = scouts.slice(0);
                } else {
                    return state;
                }

                scouts.push({
                    id: -1,
                    name: '',
                    group: state.groups[0].id
                });

                return {
                    ...state,
                    scouts,
                    scoutEditing: {
                        index: scouts.length - 1,
                    }
                };
            }

            return {
                ...state,
                scoutEditing: {
                    index: action.index
                }
            };

        case SCOUT_SAVING:
            return {
                ...state,
                scoutEditing: {
                    ...state.scoutEditing,
                    saving: true,
                }
            };

        case SCOUT_SAVE_FAILED:
            return {
                ...state,
                scoutEditing: {
                    ...state.scoutEditing,
                    saving: false,
                    error: action.error,
                }
            };

        default:
            return state;
    }
}
