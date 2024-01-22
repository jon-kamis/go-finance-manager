import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import { LineChart } from '@mui/x-charts/LineChart'
import Toast from "../alerting/Toast";
import { format, parseISO } from "date-fns";

const PortfolioHistory = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

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
        fetch(`${apiUrl}/users/${userId}/stock-portfolio-history`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setPortfolioHistory(data);
                    console.log(data)
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

    }, []);

    return (
        <div className="content">
            { portfolioHistory && portfolioHistory.count > 0 &&
                <LineChart
                    series={[{
                        data: portfolioHistory.items.map((v) => (v.balance)),
                        showMark: false,
                        color: portfolioHistory.items[0].balance > portfolioHistory.items[portfolioHistory.count - 1].balance ? "red" : "green"
                    }]}
                    xAxis={[{ scaleType: 'point', data: portfolioHistory.items.map((v) => format(parseISO(v.date), 'MMM do yyyy')) }]}
                    height={500}
                />
            }
        </div>

    )
}
export default PortfolioHistory;