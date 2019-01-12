import {coordination} from "./coordination";
import {SET_RIDE_ROLE, SET_RIDES, SET_SCOUT_PARTICIPATE} from "../constants/action-types";
//import {forEach} from "lodash";
import {cloneAndPatch} from "../utils/patch";


export function excursion(
    state = {
        error: null,
        loading: false,
        data: {
            detail: {
                id: 78,
                date: '2018-12-12',
                from: '16:00',
                to: '18:30',
                location: 'Tana',
            },
            scouts: {
                1: {
                    id: 1,
                    name: 'Anita',
                    participate: true,
                },
                4: {
                    id: 4,
                    name: 'Giuliano',
                    participate: true,
                },
                3: {
                    id: 3,
                    name: 'Andrea',
                    participate: true,
                },
                6: {
                    id: 6,
                    name: 'Marta',
                    participate: true,
                },
                7: {
                    id: 7,
                    name: 'Cassandra',
                    participate: true,
                },
                88: {
                    id: 88,
                    name: 'Greta',
                    participate: false,
                },
                8: {
                    id: 8,
                    name: 'Giuseppina',
                    participate: true,
                },
                9: {
                    id: 9,
                    name: 'Morgana',
                    participate: true,
                },
            },
            tutors: {
                4: {
                    id: 4,
                    name: 'Serena',
                    // email,
                    // verified_email
                    address: 'San Cassiano',
                    rides: 1, // enum 0/1/2: 0 neither Out nor Return, 1 either Out or Return, 2 both Out and Return if needed
                },
                6: {
                    id: 6,
                    name: 'Sonia',
                    // email,
                    // verified_email
                    address: 'Ponte',
                    rides: 1,
                },
                2: {
                    id: 2,
                    name: 'Ilenia',
                    // verified_email
                    address: 'Sant\'Alessio',
                    // email,
                    rides: 1,
                },
                5: {
                    id: 5,
                    name: 'Simone',
                    // email
                    // verified_email
                    address: 'San Concordio',
                    rides: 1,
                },
            },
            // out
            // return
        },
    },
    action
) {
    switch (action.type) {
        case SET_SCOUT_PARTICIPATE:
            return cloneAndPatch(
                state,
                ['data', 'scouts', action.scoutId, 'participate'],
                action.participate
            );

        case SET_RIDES:
            return cloneAndPatch(
                state,
                ['data', 'tutors', action.tutorId, 'rides'],
                action.rides
            );

        case SET_RIDE_ROLE: {
            state = { ...state };

            const dir = action.direction;

            state[dir] = coordination(state[dir], action);

            if (state.data.tutors[action.tutorId].rides === 1) {
                const otherDir = dir === 'out' ? 'return' : 'out';
                action.role = action.role === 'R' ? 'F' : 'R';

                state[otherDir] = coordination(state[otherDir], action);
            }

            return state;
        }
        default:
            return {
                ...state,
                out: coordination(state.out, action),
                return: coordination(state.return, action),
            };
    }
}
