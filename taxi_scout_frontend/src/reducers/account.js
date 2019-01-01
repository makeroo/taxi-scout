export function account(
    state = {
        error: null,
        loading: false,
        data: {
            id: 5,
            name: 'Simone',
            email: 'makeroo@gmail.com',
            verified_email: 1,
        },
    },
    action
) {
    switch (action.type) {
        default:
            return state;
    }
}
