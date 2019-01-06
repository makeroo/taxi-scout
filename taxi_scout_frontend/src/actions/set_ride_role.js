import {SET_RIDE_ROLE} from "../constants/action-types";

export function setRideRole(tutorId, role, direction) {
    return {
        type: SET_RIDE_ROLE,
        tutorId,
        role,
        direction
    };
}
