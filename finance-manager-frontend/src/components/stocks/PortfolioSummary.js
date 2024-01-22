import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import { LineChart } from '@mui/x-charts/LineChart'
import Toast from "../alerting/Toast";
import { format, parseISO } from "date-fns";
import PortfolioHistory from "./PortfolioHistory";
import PositionDetail from "./PositionDetail"
import Select from "../form/Select";

const PortfolioSummary = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const [portfolioSummary, setPortfolioSummary] = useState([]);
    const [posHist, setPosHist] = useState([]);
    const [histLength, setHistLength] = useState("month")

    let { userId } = useParams();

    const navigate = useNavigate();

    const handleChange = () => (event) => {
        let value = event.target.value;
        setHistLength(value);
    }

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

            fetch(`${apiUrl}/stocks?tickers=${stockStr}&histLength=${histLength}`, requestOptions)
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

    }, [portfolioSummary, histLength]);

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchPortfolioSummary();

    }, []);

    return (
        <>
            <div className="container-fluid content p-4">
                <h2>Current Statistics</h2>
                <div className="col-md-12 d-flex justify-content-between">

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
                <br />
                <div className="d-flex">

                    <div className="flex-col col-md-12">
                        <div className="d-flex col-md-12 justify-content-between">

                            <div className="flex-col col-md-10">
                                <h2>Balance</h2>
                            </div>
                            <div className="flex-col col-md-2">
                                <Select
                                    title={"History Length"}
                                    className={"form-control"}
                                    name={"histLength"}
                                    value={histLength}
                                    onChange={handleChange()}
                                    options={[{ id: "week", value: "week" }, { id: "month", value: "month" }, { id: "year", value: "year" }]}
                                    placeHolder={"Select"}
                                />
                            </div>
                        </div>
                        <PortfolioHistory histLength={histLength} />
                    </div>
                </div>
            </div>
            <div className="content container-fluid p-4">
                <div className="content-xtall-tablecontainer ">
                    <h2>Current Positions</h2>

                    {posHist && posHist.length > 0 && posHist.map((p) => (

                        <PositionDetail position={p} portfolioSummary={portfolioSummary.positions[portfolioSummary.positions.findIndex(i => i.ticker == p.ticker)]}/>

                    ))}
                    <div className="container-fluid">
                        <div className="flex-col p-4 col-md-12">
                            <div className="flex-row">
                            </div>
                            <div className="">

                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </>
    )
}
export default PortfolioSummary;