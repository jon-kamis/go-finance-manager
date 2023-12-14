import { useEffect, useState } from 'react';
import { Link, useNavigate, useOutlet, useOutletContext, useParams } from 'react-router-dom';
import { BarPlot } from '@mui/x-charts/BarChart';
import { ChartContainer } from '@mui/x-charts/ChartContainer';
import { ChartsXAxis } from '@mui/x-charts/ChartsXAxis';
import { ChartsYAxis } from '@mui/x-charts/ChartsYAxis';

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
                                            <div className="dashboard-item-tall-tablecontainer">
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
                                        </div>
                                        {/* Incomes */}
                                        <div className="p-4 dashboard-item dashboard-item-tall">
                                            <h2>Incomes</h2>
                                            <div className="dashboard-item-tall-tablecontainer">
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
                            </div>
                            <div className="container-fluid">
                                <div className="col-md-12 d-flex flex-column">
                                    {summary && summary.expenseSummary &&
                                        <div className="p-4 dashboard-item">
                                            <h1>Expenses by Category</h1>
                                            <div className="dashboard-center">
                                                <ChartContainer

                                                    xAxis={[
                                                        {
                                                            id: 'categories',
                                                            data: ['Taxes', 'Loans', 'Bills'],
                                                            scaleType: 'band',
                                                        },
                                                    ]}
                                                    series={[
                                                        {
                                                            type: 'bar',
                                                            stack: '',
                                                            yAxisKey: 'cost',
                                                            data: [summary.expenseSummary.taxes, summary.expenseSummary.loanCost, summary.expenseSummary.bills],
                                                        },
                                                    ]}
                                                    yAxis={[
                                                        {
                                                            id: 'cost',
                                                            scaleType: 'linear',
                                                        },
                                                    ]}
                                                    width={1000}
                                                    height={300}

                                                >
                                                    <BarPlot
                                                        margin={{
                                                            left: 20,
                                                            right: 10,
                                                            top: 0,
                                                            bottom: 0,
                                                        }}
                                                    />
                                                    <ChartsXAxis label="Category" position="bottom" axisId="categories" />
                                                    <ChartsYAxis label="Cost ($)" position="left" axisId="cost" />
                                                </ChartContainer>
                                            </div>
                                        </div>
                                    }
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