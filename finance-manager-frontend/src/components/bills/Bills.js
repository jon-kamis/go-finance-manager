import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, parseISO } from "date-fns";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const Bills = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { loggedInUserId } = useOutletContext();

    const [bills, setBills] = useState([]);
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

        fetch(`${apiUrl}/users/${userId}/bills${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setBills(data);
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

        fetch(`${apiUrl}/users/${userId}/bills`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setBills(data);
                }
            })
            .catch(err => {
                console.log(err)
                Toast(err.message, "error")
            })

    }, []);

    return (
        <div className="container-fluid">
            <h2>Bills</h2>
            <hr />
            <Input
                title={"Search"}
                type={"text"}
                className={"form-control"}
                name={"search"}
                value={search}
                onChange={handleChange("")}
            />
            <div className="chartContent">
            <table className="table table-striped table-hover">

                <thead>
                    <tr>
                        <th className="text-start">Name</th>
                        <th className="text-start">Amount</th>
                    </tr>
                </thead>
                <tbody>
                    {bills.map((b) => (
                        <>
                            <tr key={b.id}>
                                <td className="text-start">
                                    <Link to={`/users/${userId}/bills/${b.id}`}>
                                        {b.name}
                                    </Link>
                                </td>
                                <td className="text-start">${Intl.NumberFormat("en-US", interestFormatOptions).format(b.amount)}</td>
                            </tr>
                        </>
                    ))}
                </tbody>
            </table>
            </div>
            <Link to={`/users/${loggedInUserId}/bills/new`}><span className="badge bg-success">Add New</span></Link>
        </div>
    )
}
export default Bills;