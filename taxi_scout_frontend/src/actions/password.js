import {
    PASSWORD_SAVE,
} from "../constants/action-types";
import {jsonFetch, parseError} from "../utils/json_fetch";

// account related page loading

// account editing

export function savePassword(account, newPassword) {
    return (dispatch) => {
        return jsonFetch(`/account/${account.id}/password`, { 'p': newPassword });
    };
}
