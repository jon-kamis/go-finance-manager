import { useState } from "react";
import { jwtDecode } from "jwt-decode";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import Input from "./form/Input";
import Toast from "./alerting/Toast";

const Login = () => {

    const { setJwtToken } = useOutletContext();
    const { setRoles } = useOutletContext();
    const { setLoggedInUserId } = useOutletContext();
    const { setLoggedInUserName } = useOutletContext();
    const { toggleRefresh } = useOutletContext();

    const navigate = useNavigate();

    const [loginRequest, setLoginRequest] = useState({
        username: "",
        password: "",
    })

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setLoginRequest({
            ...loginRequest,
            [name]: value,
        })
    }

    const setJwtRoles = ((token) => {
        setRoles(jwtDecode(token).roles.split(",").map((r) => r.trimStart()))
    });

    const setUserInfo = ((token) => {
        setLoggedInUserId(jwtDecode(token).sub)
        setLoggedInUserName(jwtDecode(token).name)
    });

    const handleSubmit = (event) => {
        event.preventDefault();

        const requestOptions = {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: "include",
            body: JSON.stringify(loginRequest),
        }

        fetch(`/authenticate`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error");
                } else {
                    setJwtToken(data.access_token);
                    setJwtRoles(data.access_token);
                    setUserInfo(data.access_token);
                    toggleRefresh(true);
                    Toast("Login successful!", "success");
                    navigate("/");
                }
            })
            .catch(error => {
                Toast(error.message, "error");
            })
    }

    return (
        <div className="container-fluid">

            <div className="d-flex justify-content-around">
                <div className="p-4 content flex-column col-md-6">
                    <h1>Login</h1>
                    <form onSubmit={handleSubmit}>
                        <Input
                            title={"Username"}
                            type={"text"}
                            className={"form-control"}
                            name={"username"}
                            value={loginRequest.username}
                            onChange={handleChange("")}
                        />
                        <Input
                            title={"Password"}
                            type={"password"}
                            className={"form-control"}
                            name={"password"}
                            value={loginRequest.password}
                            onChange={handleChange("")}
                        />
                        <Input
                            type="submit"
                            className="btn btn-primary"
                            value="Login"
                        />
                    </form>
                    <hr />
                    <p>Don't have an account? <Link to="/register">Register</Link></p>
                </div>
            </div>
        </div>
    )
}
export default Login;