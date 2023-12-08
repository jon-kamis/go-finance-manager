import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "./form/Input";
import Toast from "./alerting/Toast";

const Users = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { loggedInUserId } = useOutletContext();

    const [loans, setLoans] = useState([]);
    const [error, setError] = useState(false);
    const [search, setSearch] = useState("");

    const navigate = useNavigate();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 2 }

    let { userId } = useParams();

    const handleChange = () => (event) => {
        let value = event.target.value;
        setSearch(value)

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

        let searchUrl = ""
        {
            value !== ""
            ? searchUrl = `?search=${value}`
            : searchUrl = ``
        }

        fetch(`${apiUrl}/users/${userId}/loans${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setLoans(data);
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

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/users/${userId}/loans`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setLoans(data);
                }
            })
            .catch(err => {
                console.log(err)
                Toast(err.message, "error")
            })

    }, []);

    return (
        <div>
            <h2>Loans</h2>
            <hr />
            <Input
                title={"Search"}
                type={"text"}
                className={"form-control"}
                name={"search"}
                value={search}
                onChange={handleChange("")}
            />
            <table className="table table-striped table-hover">

                <thead>
                    <tr>
                        <th className="text-end">Name</th>
                        <th className="text-end">Principal</th>
                        <th className="text-end">Rate</th>
                        <th className="text-end">Term</th>
                        <th className="text-end">Monthly Payment</th>
                        <th className="text-end">Total Interest</th>
                        <th className="text-end">Total Cost</th>
                    </tr>
                </thead>
                <tbody>
                    {loans.map((l) => (
                        <>
                            <tr key={l.id}>
                                <td className="text-end">
                                    <Link to={`/users/${userId}/loans/${l.id}`}>
                                        {l.name}
                                    </Link>
                                </td>
                                <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.total)}</td>
                                <td className="text-end">{Intl.NumberFormat("en-US", interestFormatOptions).format(l.interestRate)}</td>
                                <td className="text-end">{l.loanTerm}</td>
                                <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.monthlyPayment)}</td>
                                <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.interest)}</td>
                                <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.totalCost)}</td>
                            </tr>
                        </>
                    ))}
                </tbody>
            </table>
            <Link to={`/users/${loggedInUserId}/loans/new`}><span className="badge bg-success">Create New</span></Link>
        </div>
    )
}
export default Users;