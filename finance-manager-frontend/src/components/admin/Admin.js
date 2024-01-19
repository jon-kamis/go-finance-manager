import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import Users from "../users/Users"
import EnableStocks from "./EnableStocks";

const Admin = () => {
    const { jwtToken } = useOutletContext();
    const navigate = useNavigate();


    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
    }, []);

    return (
        <div className="container-fluid">
            <h1>Administration</h1>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content">
                    <EnableStocks />
                </div>
            </div>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content content-xtall">
                    <Users />
                </div>
            </div>
        </div>
    )
}
export default Admin;