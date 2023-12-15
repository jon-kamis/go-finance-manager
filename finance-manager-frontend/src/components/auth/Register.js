import { useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const Register = () => {
    const { setJwtToken } = useOutletContext();
    const { setAlertClassName } = useOutletContext();
    const { setAlertMessage } = useOutletContext();
    const { toggleRefresh } = useOutletContext();

    const navigate = useNavigate();
    const [errors, setErrors] = useState([]);

    const hasError = (key) => {
        return errors.indexOf(key) !== -1;
    }

    const [user, setUser] = useState({
        username: "",
        email: "",
        firstName: "",
        lastName: "",
        password: ""
    })

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setUser({
            ...user,
            [name]: value,
        })
    }

    const handleSubmit = (event) => {
        event.preventDefault();

        const requestOptions = {
            method: "POST",
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: "include",
            body: JSON.stringify(user, null, 3),
        }

        fetch(`/register`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setJwtToken("");
                    Toast("Account Created!", "success")
                    navigate("/login");
                }
            })
            .catch(error => {
                Toast(error.message, "error")
            })

    }

    return (
        <>
            <div className="col-md-6 offset-md-3">
                <h2>Register</h2>
                <hr />
                <form onSubmit={handleSubmit}>
                    <input type="hidden" name="id" value={user.id}></input>
                    <Input
                        title={"Username"}
                        type={"text"}
                        className={"form-control"}
                        name={"username"}
                        value={user.username}
                        onChange={handleChange("")}
                        errorDiv={hasError("username") ? "text-danger" : "d-none"}
                        errorMsg={"Please enter a username"}
                    />
                    <Input
                        title={"First Name"}
                        type={"text"}
                        className={"form-control"}
                        name={"firstName"}
                        value={user.firstName}
                        onChange={handleChange("")}
                        errorDiv={hasError("firstName") ? "text-danger" : "d-none"}
                        errorMsg={"Please enter a first name"}
                    />
                    <Input
                        title={"Last Name"}
                        type={"text"}
                        className={"form-control"}
                        name={"lastName"}
                        value={user.lastName}
                        onChange={handleChange("")}
                        errorDiv={hasError("lastName") ? "text-danger" : "d-none"}
                        errorMsg={"Please enter a last name"}
                    />
                    <Input
                        title={"Email"}
                        type={"email"}
                        className={"form-control"}
                        name={"email"}
                        value={user.email}
                        onChange={handleChange("")}
                        errorDiv={hasError("email") ? "text-danger" : "d-none"}
                        errorMsg={"Please enter a email"}
                    />
                    <Input
                        title={"Password"}
                        type={"password"}
                        className={"form-control"}
                        name={"password"}
                        value={user.password}
                        onChange={handleChange("")}
                        errorDiv={hasError("password") ? "text-danger" : "d-none"}
                        errorMsg={"Please enter a password"}
                    />
                    <hr />
                    <Input
                        type="submit"
                        className="btn btn-primary"
                        value="Register"
                    />
                </form>
                <hr />
                <p>Already have an account? <Link to="/login">Login</Link></p>
            </div>
        </>
    )
}
export default Register;