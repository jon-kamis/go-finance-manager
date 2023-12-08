import { useEffect, useState } from 'react';
import { Link, useNavigate, useOutlet, useOutletContext, useParams } from 'react-router-dom';
import Toast from './alerting/Toast';
const Home = () => {
    const { jwtToken } = useOutletContext();
    const { apiUrl } = useOutletContext();
    const { loggedInUserId } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const [loanSummary, setLoanSummary] = useState();
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

        fetch(`${apiUrl}/users/${loggedInUserId}/loan-summary`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.message)
                    Toast("An unexpected error occured when loading loan summary", "error")
                } else {
                    setLoanSummary(data);
                }
            })
            .catch(err => {
                console.log(err.message)
                Toast("An unexpected error occured when loading loan summary", "error")
            })

    }, []);

    return (
        <>
            <div className="text-start">
                <h2>
                    Welcome to Finance Manager!
                    <br/><br/>
                    Dashboard
                    <br/><br/>
                </h2>
                {jwtToken !== ""
                    ? <>
                        <div className="row">
                        <h3>Loans</h3>
                            <div className="col-md-12 text-start">
                                <table className="table table-striped table-hover">

                                    <thead>
                                        <tr>
                                            <th className="text-start">Total Number</th>
                                            <th className="text-start">Total Debt</th>
                                            <th className="text-start">Monthly Cost</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {loanSummary &&
                                            <>
                                                <td className="text-start">{loanSummary.count}</td>
                                                <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(loanSummary.totalBalance)}</td>
                                                <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(loanSummary.monthlyCost)}</td>
                                            </>
                                        }
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </>
                    : <>
                        <p>Please <Link to="/login">Login</Link> to view personalized dashboard</p>
                    </>
                }
            </div>
        </>
    )
}
export default Home;