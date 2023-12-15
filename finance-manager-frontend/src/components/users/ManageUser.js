import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";

const Users = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [user, setUser] = useState([]);
    const [error, setError] = useState(false);

    const navigate = useNavigate();

    let { id } = useParams();

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/users/${id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUser(data);
                setError(false);
            })
            .catch(err => {
                console.log(err)
                setUser([]);
                setError(true);
            })

    }, []);

    return (
        <div>
            <h2>Manage User</h2>
            <hr />
            {!error &&
                <>
                    {!error &&
                        <>
                            <table className="table">
                                <thead>
                                    <tr>
                                        <th>ID</th>
                                        <th>First Name</th>
                                        <th>Last Name</th>
                                        <th>Username</th>
                                        <th>Email</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>{user.id}</td>
                                        <td>{user.firstName}</td>
                                        <td>{user.lastName}</td>
                                        <td>{user.username}</td>
                                        <td>{user.email}</td>
                                    </tr>
                                </tbody>
                            </table>
                        </>
                    }
                </>
            }
        </div>
    )
}
export default Users;