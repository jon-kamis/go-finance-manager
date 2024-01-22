import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import { LineChart } from '@mui/x-charts/LineChart'
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import { format, parseISO } from "date-fns";
import PortfolioHistory from "./PortfolioHistory";

const PortfolioSummary = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const [portfolioSummary, setPortfolioSummary] = useState([]);
    const [posHist, setPosHist] = useState([]);

    let { userId } = useParams();

    const navigate = useNavigate();


    function fetchPortfolioSummary() {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/users/${userId}/stocks`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setPortfolioSummary(data);
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        if (portfolioSummary.positions) {
            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)
            const requestOptions = {
                method: "GET",
                headers: headers,
            }

            let stockList = []
            portfolioSummary.positions.map(s => stockList.push(s.ticker))
            let stockStr = stockList.join(",")

            fetch(`${apiUrl}/stocks?tickers=${stockStr}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast(data.message, "error")
                    } else {
                        setPosHist(data);
                    }
                })
                .catch(err => {
                    Toast(err.message, "error")
                    console.log(err)
                })
        }

    }, [portfolioSummary]);

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchPortfolioSummary();

    }, []);

    return (
        <>
            <div className="container-fluid content">
                <div className="col-md-12 d-flex p-4 justify-content-between">
                    <div className="flex-col">
                        <h5>Current Value</h5>
                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioSummary ? portfolioSummary.currentValue : 0)}</h4>
                    </div>
                    <div className="flex-col">
                        <h5>Daily High</h5>
                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioSummary ? portfolioSummary.currentHigh : 0)}</h4>
                    </div>
                    <div className="flex-col">
                        <h5>Daily Low</h5>
                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioSummary ? portfolioSummary.currentLow : 0)}</h4>
                    </div>
                    <div className="flex-col">
                        <h5>Daily Open</h5>
                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioSummary ? portfolioSummary.currentOpen : 0)}</h4>
                    </div>
                    <div className="flex-col">
                        <h5>Daily Close</h5>
                        <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioSummary ? portfolioSummary.currentClose : 0)}</h4>
                    </div>
                    <div className="flex-col">
                        <h5>As of Date</h5>
                        <h4>{portfolioSummary && portfolioSummary.asOf ? format(parseISO(portfolioSummary.asOf), 'MMM do yyyy') : "-"}</h4>
                    </div>
                </div>
                <div className="d-flex">
                    <div className="flex-col p-4 col-md-3">
                        <div className="flex-row">
                            <h2>Positions</h2>
                        </div>
                        <div className="content-xtall-tablecontainer">
                            {posHist && posHist.length > 0 && posHist.map((p) => (

                                <div className="flex-row">

                                    <h4>{p.ticker}</h4>
                                    <LineChart
                                        series={[{
                                            data: p.values.map((v) => (v.close)),
                                            showMark: false,
                                            color: p.values[0].close > p.values[p.count - 1].close ? "red" : "green"
                                        }]}
                                        xAxis={[{ scaleType: 'point', data: p.values.map((v) => format(parseISO(v.date), 'MMM do yyyy')) }]}
                                        height={200}
                                    />

                                </div>

                            ))}
                        </div>
                    </div>
                    <div className="flex-col p-4 col-md-9">
                        <div className="flex-row">
                            <h2>Portfolio Balance</h2>
                        </div>
                        <PortfolioHistory />
                    </div>
                </div>
            </div>
        </>
    )
}
export default PortfolioSummary;