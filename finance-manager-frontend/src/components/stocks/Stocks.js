import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import PortfolioSummary from "./PortfolioSummary";

const Stocks = () => {
    const { jwtToken } = useOutletContext();
    const navigate = useNavigate();


    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
    }, []);

    return (
        <div className="container-fluid">
            <h1>Portfolio</h1>
            <PortfolioSummary />
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content content-xtall">
                </div>
            </div>
        </div>
    )
}
export default Stocks;