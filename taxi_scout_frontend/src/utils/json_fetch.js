import {BASE_URL} from "../constants/rest_api";
import {SERVER_ERROR, SERVICE_NOT_AVAILABLE} from "../constants/errors";

const ERROR_SERVICE_NOT_AVAILABLE = {
    error: SERVICE_NOT_AVAILABLE
};

const ERROR_SERVER_ERROR = {
    error: SERVER_ERROR
};

export function jsonFetch(path, payload) {
    let fetch_config = {
        credentials: 'same-origin',
    };

    if (payload) {
        fetch_config.method = 'POST';
        fetch_config.body = JSON.stringify(payload);
    }

    return fetch(`${BASE_URL}${path}`, fetch_config).then((response) => {
        switch (response.status) {
            case 502:
                throw ERROR_SERVICE_NOT_AVAILABLE;

            default:
                return response.json();
        }
    }).then((response) => {
        if (response && 'error' in response)
            throw response;

        return response;
    });
}

export function parseError(error) {
    if (error && typeof error.error === 'string')
        return error;

    console.log('unexpected error', error);

    return ERROR_SERVER_ERROR;
}
