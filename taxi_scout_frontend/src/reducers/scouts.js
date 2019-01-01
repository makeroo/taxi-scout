export function scouts(
    state = {
        error: null,
        loading: false,
        items: {
            7: {
                id: 7,
                name: 'Cassandra',
            },
            88: {
                id: 88,
                name: 'Greta',
            },
            8: {
                id: 8,
                name: 'Giuseppina',
            },
            9: {
                id: 9,
                name: 'Morgana',
            },
        },
    },
    action
) {
    switch (action.type) {
        default:
            return state;
    }
}
