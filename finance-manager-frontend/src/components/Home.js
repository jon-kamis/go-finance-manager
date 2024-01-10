import { useEffect, useState } from 'react';
import { Link, useNavigate, useOutlet, useOutletContext, useParams } from 'react-router-dom';
import { PieChart } from '@mui/x-charts/PieChart';

const Home = () => {
    const { jwtToken } = useOutletContext();
    const { apiUrl } = useOutletContext();
    const { loggedInUserId } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const [summary, setSummary] = useState();

    const navigate = useNavigate();

    function getSource(source) {
        switch (source) {
            case "bill":
                return "Bill"
            case "credit-card":
                return "Credit Card"
            case "taxes":
                return "Taxes"
            case "loan":
                return "Loan"
        }
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
                {jwtToken !== ""
                    ? <>

                        <h1></h1>
                        <div className="container-fluid">
                            <h1>Dashboard</h1>
                            <div className="col-md-12 d-flex flex-column">
                                {summary && summary.expenseSummary &&
                                    <div className="p-4 content">
                                        <h2>Monthly Breakdown</h2>
                                        <div className="d-flex justify-content-around">
                                            <div className="flex-col">
                                                <PieChart
                                                    series={[
                                                        {
                                                            data: [
                                                                { id: 0, value: summary.expenseSummary.bills, label: 'Bills' },
                                                                { id: 1, value: summary.expenseSummary.creditCards, label: 'Credit Cards' },
                                                                { id: 2, value: summary.expenseSummary.loanCost, label: 'Loans' },
                                                                { id: 3, value: summary.expenseSummary.taxes, label: 'Taxes' },
                                                                { id: 4, value: summary.netFunds, label: 'Net Funds' }
                                                            ]
                                                        }
                                                    ]}
                                                    width={600}
                                                    height={400}
                                                />
                                            </div>
                                        </div>
                                    </div>
                                }
                            </div>
                            <div className="d-flex justify-content-between">
                                <div className="col-md-4 d-flex flex-column">
                                    <div className="p-4 content">
                                        <h2>Income</h2>
                                        <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.incomeSummary ? summary.incomeSummary.totalIncome : 0)}</h1>
                                    </div>
                                    <div className="p-4 content">
                                        <h2>Expenses</h2>
                                        <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.expenseSummary ? summary.expenseSummary.totalCost : 0)}</h1>
                                    </div>
                                    <div className="p-4 content">
                                        <h2>Net Funds</h2>
                                        <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary ? summary.netFunds : 0)}</h1>
                                    </div>
                                    <div className="p-4 content">
                                        <h2>Total Debt</h2>
                                        <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.expenseSummary ? summary.expenseSummary.totalBalance : 0)}</h1>
                                    </div>
                                </div>
                                <div className="col-md-8 d-flex flex-column">
                                    {/* Incomes */}
                                    <div className="p-4 content content-tall">
                                        <h2>Incomes</h2>
                                        <div className="content-tall-tablecontainer">
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
                                    {/* Expenses */}
                                    <div className="p-4 content content-tall">
                                        <h2>Expenses</h2>
                                        <div className="content-tall-tablecontainer">
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
                                                                    <td className="text-start">{getSource(e.source)}</td>
                                                                    <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(e.amount)}</td>
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
                        </div>
                        <div className="container-fluid d-flex justify-content-between">
                            <div className="col-md-4 d-flex flex-column">
                                <div className="p-4 content">
                                    <h2>Total Credit</h2>
                                    <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.creditSummary ? summary.creditSummary.total : 0)}</h1>
                                </div>
                            </div>
                            <div className="col-md-4 d-flex flex-column">
                                <div className="p-4 content">
                                    <h2>Available Credit</h2>
                                    <h1>${Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.creditSummary ? summary.creditSummary.available : 0)}</h1>
                                </div>
                            </div>
                            <div className="col-md-4 d-flex flex-column">
                                <div className="p-4 content">
                                    <h2>Credit Utilization</h2>
                                    <h1>{Intl.NumberFormat("en-US", numberFormatOptions).format(summary && summary.creditSummary ? summary.creditSummary.utilization : 0)}%</h1>
                                </div>
                            </div>
                        </div>
                    </>
                    :
                    <>
                        <div className="d-flex">
                            <div className="p-4 col-md-12 content text-center">
                                <h1>Welcome to Finance Manager!</h1>
                                <br />
                                <h3>This application is designed to help manage personal finances. This is a personal project not intended for public use. Values are not gauranteed to be accurate</h3>
                                <br />
                                <br />
                                <h2>Please <Link to="/login">Login</Link> to access the application</h2>
                            </div>
                        </div>
                    </>
                }
            </div>
        </>
    )
}
export default Home;