import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import Input from "./form/Input";

const Users = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [users, setUsers] = useState([]);
    const [error, setError] = useState(false);
    const [search, setSearch] = useState("");

    const navigate = useNavigate();

    const handleChange = () => (event) => {
        let value = event.target.value;
        setSearch(value)
        
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

        let searchUrl = ""
        {value !== "" 
        ? searchUrl = `?search=${value}`
        : searchUrl = ``}

        fetch(`${apiUrl}/users/all${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUsers(data);
                setError(false);
            })
            .catch(err => {
                console.log(err)
                setUsers([]);
                setError(true);
            })
    }

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

        fetch(`${apiUrl}/users/all`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUsers(data);
                setError(false);
            })
            .catch(err => {
                console.log(err)
                setUsers([]);
                setError(true);
            })

    }, []);

    return (
        <div>
            <h2>Users</h2>
            <hr />
                <Input
                    title={"Search"}
                    type={"text"}
                    className={"form-control"}
                    name={"search"}
                    value={search}
                    onChange={handleChange("")}
                />
            <table className="table table-striped table-hover">

                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Username</th>
                        <th>Email</th>
                    </tr>
                </thead>
                <tbody>
                    {!error &&
                        <>
                            {users.map((u) => (
                                <tr key={u.id}>
                                    <td>
                                        <Link to={`/admin/users/${u.id}`}>
                                            {u.lastName}, {u.firstName}
                                        </Link>
                                    </td>
                                    <td>{u.username}</td>
                                    <td>{u.email}</td>
                                </tr>
                            ))}
                        </>
                    }
                </tbody>
            </table>
        </div>
    )
}
export default Users;