import {SET_RIDE_ROLE} from "../constants/action-types";
import {cloneAndPatch} from "../utils/patch";


export function coordination(
    state = {
        tutors: {
            4: {
                id: 4,
                role: 'R',
                free_seats: 2,
                scouts: [1],
            },
            6: {
                id: 6,
                role: 'R',
                free_seats: 0,
                scouts: [4],
            },
            2: {
                id: 2,
                role: 'R',
                free_seats: 0,
                scouts: [3, 6],
            },
            5: {
                id: 5,
                role: 'F',
                free_seets: 2,
                scouts: [7, 88, 8, 9],
            }
        },
        meetings: [
            {
                id: 45,
                taxi: 6,
                riders: [6],
                point: 'esselunga',
                time: '15:45',
            },
        ]
    },
    action
) {
    switch (action.type) {
        case SET_RIDE_ROLE: {
            return cloneAndPatch(
                state,
                ['tutors', action.tutorId, 'role'],
                action.role
            );
        }
        default:
            return state;
    }
}
