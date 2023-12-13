import { useEffect, useState } from 'react';
import { Link, useNavigate, useOutlet, useOutletContext, useParams } from 'react-router-dom';
import Toast from './alerting/Toast';
const Home = () => {
    const { jwtToken } = useOutletContext();
    const { loggedInUserName } = useOutletContext();
    const { apiUrl } = useOutletContext();
    const { loggedInUserId } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const [loanSummary, setLoanSummary] = useState();
    const [summary, setSummary] = useState();

    const navigate = useNavigate();


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

        fetch(`${apiUrl}/users/${loggedInUserId}/summary`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.message)
                } else {
                    setSummary(data);
                }
            })
            .catch(err => {
                console.log(err.message)
            })

    }, []);

    return (
        <>
            <div className="container-fluid">
                <div className="row">
                    {jwtToken !== ""
                        ? <>

                            <h1></h1>
                            <div className="container-fluid">
                                <h1>Dashboard</h1>
                                <div className="d-flex justify-content-between">
                                    <div className="col-md-4 d-flex flex-column">
                                        <div className="p-4 dashboard-item">
                                            <h2>Income</h2>
                                            <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.incomeSummary ? summary.incomeSummary.totalIncome : 0)}</h1>
                                        </div>
                                        <div className="p-4 dashboard-item">
                                            <h2>Expenses</h2>
                                            <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.expenseSummary ? summary.expenseSummary.totalCost : 0)}</h1>
                                        </div>
                                        <div className="p-4 dashboard-item">
                                            <h2>Net Funds</h2>
                                            <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary ? summary.netFunds : 0)}</h1>
                                        </div>
                                        <div className="p-4 dashboard-item">
                                            <h2>Total Debt</h2>
                                            <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.expenseSummary ? summary.expenseSummary.totalBalance : 0)}</h1>
                                        </div>
                                    </div>
                                    {/* Expenses */}
                                    <div className="col-md-8 d-flex flex-column">
                                        <div className="p-4 dashboard-item dashboard-item-tall">
                                            <h2>Expenses</h2>
                                            <table className="table-responsive table table-striped table-hover">
                                                <thead>
                                                    <th className="text-start">Name</th>
                                                    <th className="text-start">Source</th>
                                                    <th className="text-start">Amount</th>
                                                </thead>
                                                <tbody>
                                                    {summary && summary.expenseSummary && summary.expenseSummary.expenses &&
                                                        summary.expenseSummary.expenses.map((e) => (
                                                            <>
                                                                <tr key={e.name}>
                                                                    <td className="text-start">{e.name}</td>
                                                                    <td className="text-start">{e.source}</td>
                                                                    <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(e.amount)}</td>
                                                                </tr>
                                                            </>
                                                        ))
                                                    }
                                                </tbody>
                                            </table>
                                        </div>
                                        {/* Incomes */}
                                        <div className="p-4 dashboard-item dashboard-item-tall">
                                            <h2>Incomes</h2>
                                            <table className="table table-striped table-hover table-tall">
                                                <thead>
                                                    <th className="text-start">Name</th>
                                                    <th className="text-start">Source</th>
                                                    <th className="text-start">Amount</th>
                                                </thead>
                                                <tbody>
                                                    {summary && summary.incomeSummary && summary.incomeSummary.incomes &&
                                                        summary.incomeSummary.incomes.map((i) => (
                                                            <>
                                                                <tr key={i.name}>
                                                                    <td className="text-start">{i.name}</td>
                                                                    <td className="text-start">{i.source}</td>
                                                                    <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(i.amount)}</td>
                                                                </tr>
                                                            </>
                                                        ))
                                                    }
                                                </tbody>
                                            </table>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </>
                        : <div className="col-md-8 offset-md-2 text-center">
                            <h1>Please <Link to="/login">Login</Link> to view personalized dashboard</h1>
                        </div>
                    }
                </div>
            </div>
        </>
    )
}
export default Home;