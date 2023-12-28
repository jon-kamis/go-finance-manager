import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const Users = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [users, setUsers] = useState([]);
    const [error, setError] = useState(false);
    const [search, setSearch] = useState("");

    const navigate = useNavigate();

    function sortData(data) {
        let sortedData = data

        sortedData.sort((a, b) => a.lastName.toLowerCase() > b.lastName.toLowerCase() ? 1 : -1);

        return sortedData;
    }

    const refreshData = () => (event) => {
        let value = event.target.value;
        setSearch(value)

        fetchData(value)

        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
    }

    function fetchData(value) {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        let searchUrl = ""

        value !== ""
            ? searchUrl = `?search=${value}`
            : searchUrl = ``

        fetch(`${apiUrl}/users/all${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setUsers(sortData(data));
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })
    }


    function deleteUser(id) {
        return () => {
            console.log("attempting to delete user with id: " + id)

            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)
            const requestOptions = {
                method: "DELETE",
                headers: headers,
            }

            fetch(`${apiUrl}/users/${id}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast("An Error Occured", "error")
                    } else {
                        Toast("success!", "success")
                        fetchData("");
                    }

                })
                .catch(err => {
                    console.log(err)
                    Toast(err.message, "error")
                })
        }
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
            })
            .catch(err => {
                console.log(err)
                Toast("Failed to retrieve users", "error")
            })

    }, []);

    return (
        <div className="container-fluid">
            <h1>Users</h1>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content content-xtall">
                    <div className="row">
                        <div className="col-md-12">
                            <Input
                                title={"Search"}
                                type={"text"}
                                className={"form-control"}
                                name={"search"}
                                value={search}
                                onChange={refreshData("")}
                            />
                        </div>
                    </div>
                    <div className="content-xtall-tablecontainer">
                        <table className="table table-striped table-hover">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th>Username</th>
                                    <th>Email</th>
                                    <th></th>
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
                                                <td>
                                                    <Input
                                                        type="submit"
                                                        className="btn btn-danger"
                                                        value="Delete User"
                                                        onClick={deleteUser(u.id)} />
                                                </td>
                                            </tr>
                                        ))}
                                    </>
                                }
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    )
}
export default Users;