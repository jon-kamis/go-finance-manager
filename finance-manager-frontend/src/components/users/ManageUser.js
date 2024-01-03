import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Select from "../form/Select";

const Users = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [user, setUser] = useState([]);
    const [userRoles, setUserRoles] = useState([]);
    const [allRoles, setAllRoles] = useState([]);

    const [availableRoles, setAvailableRoles] = useState([]);

    const navigate = useNavigate();

    let { id } = useParams();

    function fetchRoles() {
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

        fetch(`${apiUrl}/users/${id}/roles`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                } else {
                    setUserRoles(data);
                }
            })
            .catch(err => {
                console.log(err)
                setUserRoles([]);
            })
    }

    function fetchAvailableRoles() {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/roles`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    setAllRoles([]);
                } else {
                    setAvailableRoles(trimRoles(data));
                }
            })
            .catch(err => {
                console.log(err)
                setAllRoles([]);
            })
    }

    function removeRole(roleId) {
        return () => {
            if (jwtToken === null || jwtToken === "") {
                navigate("/")
            }

            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)
            const requestOptions = {
                method: "DELETE",
                headers: headers,
            }

            fetch(`${apiUrl}/users/${id}/roles/${roleId}`, requestOptions)
                .then((response) => response.json())
                .then(
                    fetchRoles()
                )
                .catch(err => {
                    console.log(err)
                })

        }
    }

    function addRole(roleId) {
        return () => {
            if (jwtToken === null || jwtToken === "") {
                navigate("/")
            }

            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)
            const requestOptions = {
                method: "POST",
                headers: headers,
            }

            fetch(`${apiUrl}/users/${id}/roles/${roleId}`, requestOptions)
                .then((response) => response.json())
                .then(
                    fetchRoles()
                )
                .catch(err => {
                    console.log(err)
                })

        }
    }

    /* Trims out Roles the user already has */
    function trimRoles(roles) {
        let takenIds = userRoles.map((r) => {
            return r.roleId;
        });

        return roles.filter((r) => !takenIds.includes(r.id));

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

        fetch(`${apiUrl}/users/${id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                setUser(data);
            })
            .catch(err => {
                console.log(err)
                setUser([]);
            })

        fetchRoles()
        fetchAvailableRoles()

    }, []);

    useEffect(() => {

        fetchAvailableRoles()

    }, [userRoles])

    return (
        <>
            <h1>Manage User</h1>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content">
                    <div className="row">
                        <div className="col-md-12">
                            <h2>User Info</h2>
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
                        </div>
                    </div>
                </div>
            </div>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-6 content">
                    <div className="row">
                        <div className="col-md-12">
                            <h2>Current Roles</h2>
                            <table className="table">
                                <thead>
                                    <tr>
                                        <th>Name</th>
                                        <th></th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {userRoles && userRoles.map((r) => (
                                        <>
                                            <tr key={r.id}>
                                                <td className="text-start">{r.code}</td>
                                                <td className="text-end">
                                                    <Input
                                                        type="submit"
                                                        className="btn btn-danger"
                                                        value="Remove"
                                                        onClick={removeRole(r.id)} />
                                                </td>
                                            </tr>
                                        </>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                <div className="p-4 flex-col col-md-6 content">
                    <div className="row">
                        <div className="col-md-12">
                            <h2>Add Roles</h2>
                            <table className="table">
                                <thead>
                                    <tr>
                                        <th>Name</th>
                                        <th></th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {availableRoles && availableRoles.map((r) => (
                                        <>
                                            <tr key={r.id}>
                                                <td className="text-start">{r.code}</td>
                                                <td className="text-end">
                                                    <Input
                                                        type="submit"
                                                        className="btn btn-success"
                                                        value="Add"
                                                        onClick={addRole(r.id)} />
                                                </td>
                                            </tr>
                                        </>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>

        </>
    )
}
export default Users;