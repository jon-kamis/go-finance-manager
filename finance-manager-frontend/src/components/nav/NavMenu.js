import { Link, useOutletContext } from 'react-router-dom';

const NavMenu = (props) => {

    const NavData = [
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
            path: `/users/${props.loggedInUserId}/loans`,
        },
        {
            id: 4,
            text: "Incomes",
            requiresJwt: true,
            requiresAdmin: false,
            path: `/users/${props.loggedInUserId}/incomes`,
        },
        {
            id: 5,
            text: "Bills",
            requiresJwt: true,
            requiresAdmin: false,
            path: `/users/${props.loggedInUserId}/bills`,
        }
    ];

    const hasRole = (key) => {
        return props.roles !== null && props.roles.indexOf(key) !== -1;
    }

    return (
        <div className="navMenu d-flex justify-content-around">
                {NavData.map((n) => {
                    return (!n.requiresJwt || props.jwtToken !== "") && (!n.requiresAdmin || hasRole("admin")) &&
                        <>
                            <div className="flex-col">
                                <Link to={n.path} className="list-group-item list-group-item-action"><p>{n.text}</p></Link>
                            </div>
                        </>
                })}
        </div>
    )
};

export default NavMenu;