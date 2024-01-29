import { forwardRef, useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import { LineChart } from '@mui/x-charts/LineChart'
import { ArrowDropUp, ArrowDropDown } from '@mui/icons-material'
import Toast from "../alerting/Toast";
import { format, parseISO } from "date-fns";

const PortfolioHistory = forwardRef((props, ref) => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const [portfolioHistory, setPortfolioHistory] = useState([]);

    let { userId } = useParams();

    const navigate = useNavigate();

    function fetchPortfolioHistory() {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        console.log("fetching user stock portfolio history")
        fetch(`${apiUrl}/users/${userId}/stock-portfolio-history?histLength=${props.histLength}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setPortfolioHistory(data);
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

        fetchPortfolioHistory();

    }, [props.histLength])

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchPortfolioHistory();

    }, []);

    return (
        <div className="container-fluid content col-md-12">
            {portfolioHistory && portfolioHistory.count > 0 &&
                <>
                    <div className="d-flex justify-content-between">

                        <div className="flex-col">
                            <h5>High</h5>
                            <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioHistory.high)}</h4>
                        </div>
                        <div className="flex-col">
                            <h5>Low</h5>
                            <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioHistory.low)}</h4>
                        </div>
                        <div className="flex-col">
                            <h5>Open</h5>
                            <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioHistory.open)}</h4>
                        </div>
                        <div className="flex-col">
                            <h5>Close</h5>
                            <h4>${Intl.NumberFormat("en-US", numberFormatOptions).format(portfolioHistory.close)}</h4>
                        </div>
                        <div className="flex-col">
                            <h5>Net Change</h5>
                            <h4 className={portfolioHistory.delta > 0 ? "text-success" : "text-failure"}>
                            {portfolioHistory.delta > 0 ? <ArrowDropUp/> : <ArrowDropDown/>}
                            ${Intl.NumberFormat("en-US", numberFormatOptions).format(Math.abs(portfolioHistory.delta))}
                            </h4>
                        </div>
                        <div className="flex-col">
                            <h5>Net Change (%)</h5>
                            <h4 className={portfolioHistory.delta > 0 ? "text-success" : "text-failure"}>
                            {portfolioHistory.delta > 0 ? <ArrowDropUp/> : <ArrowDropDown/>}
                            {Intl.NumberFormat("en-US", numberFormatOptions).format(Math.abs(portfolioHistory.deltaPercentage))}%
                            </h4>
                        </div>
                    </div>
                    <div className="d-flex">

                        <LineChart
                            series={[
                                {
                                    data: portfolioHistory.items.map((v) => (v.open)),
                                    showMark: false,
                                    label: "open",
                                    color: "purple",
                                    id: "openData"
                                },
                                {
                                    data: portfolioHistory.items.map((v) => (v.close)),
                                    showMark: false,
                                    label: "close",
                                    color: portfolioHistory.items[0].close > portfolioHistory.items[portfolioHistory.count - 1].close ? "red" : "green",
                                    id: "closeData"
                                },
                                {
                                    data: portfolioHistory.items.map((v) => (v.high)),
                                    showMark: false,
                                    label: "high",
                                    color: "blue",
                                    id: "highData"
                                },
                                {
                                    data: portfolioHistory.items.map((v) => (v.low)),
                                    showMark: false,
                                    label: "low",
                                    color: "orange",
                                    id: "lowData"
                                },
                            ]}
                            xAxis={[{ scaleType: 'point', data: portfolioHistory.items.map((v) => format(parseISO(v.date), 'MMM do yyyy')) }]}
                            height={500}
                            slotProps={{
                                legend: { hidden: true }
                            }}
                            sx={{
                                '.MuiLineElement-series-openData': {
                                    stroke: 'none',
                                },
                                '.MuiLineElement-series-highData': {
                                    stroke: 'none',
                                },
                                '.MuiLineElement-series-lowData': {
                                    stroke: 'none',
                                }
                            }}
                        />
                    </div>
                </>
            }
        </div>

    )
});

export default PortfolioHistory;