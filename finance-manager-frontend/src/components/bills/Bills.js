import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, parseISO } from "date-fns";
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import ManageBill from "./ManageBill";
import NewBill from "./NewBill";

const Bills = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { loggedInUserId } = useOutletContext();

    const [bills, setBills] = useState([]);
    const [search, setSearch] = useState("");
    const [selectedBillId, setSelectedBillId] = useState("");

    const navigate = useNavigate();

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

    function fetchData() {

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

    }

    function updateSelectedId(id) {
        return () => {
            console.log(id)
            setSelectedBillId(id)
        }
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData()
        setSelectedBillId(bills != null && bills.length > 0 ? bills[0].id : null)

    }, []);

    return (
        <div className="container-fluid">
            <h1>Bills</h1>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content content-xtall">
                    <Input
                        title={"Search"}
                        type={"text"}
                        className={"form-control"}
                        name={"search"}
                        value={search}
                        onChange={handleChange("")}
                    />
                    <div className="content-xtall-tablecontainer">

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
                                        <tr key={b.id} onClick={updateSelectedId(b.id)}>
                                            <td className="text-start">{b.name}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", interestFormatOptions).format(b.amount)}</td>
                                        </tr>
                                    </>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <div className="d-flex">
                <div className="p-4 col-md-6 content">
                    <NewBill search={setSearch} data={setBills} />
                </div>
                <div className="p-4 col-md-6 content">
                    <ManageBill fetchData={fetchData} billId={selectedBillId} setBillId={setSelectedBillId} />
                </div>
            </div>
        </div>
    )
}
export default Bills;