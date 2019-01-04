import {SET_RIDES} from "../constants/action-types";

export function setRides(tutorId, rides) {
    return {
        type: SET_RIDES,
        tutorId,
        rides
    };
}
