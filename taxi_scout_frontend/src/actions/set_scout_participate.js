import {SET_SCOUT_PARTICIPATE} from "../constants/action-types";

export function setScoutParticipate(scoutId, participate) {
    return {
        type: SET_SCOUT_PARTICIPATE,
        scoutId,
        participate
    };
}
