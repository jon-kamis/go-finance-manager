export const NavData = [
    {
        id: 0,
        text: "Home",
        requiresJwt: false,
        requiresAdmin: false,
        path: "/",
    },
    {
        id: 1,
        text: "About",
        requiresJwt: false,
        requiresAdmin: false,
        path: "/about",
    },
    {
        id: 2,
        text: "Users",
        requiresJwt: true,
        requiresAdmin: true,
        path: "/users",
    },
    {
        id: 3,
        text: "Loans",
        requiresJwt: true,
        requiresAdmin: false,
        path: "/users/:userId/loans",
    },
    {
        id: 4,
        text: "Incomes",
        requiresJwt: true,
        requiresAdmin: false,
        path: "/users/:userId/incomes",
    }
]