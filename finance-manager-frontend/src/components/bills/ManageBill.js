import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, utcToZonedTime, zonedTimeToUtc } from 'date-fns-tz'

import Input from "../form/Input";
import Toast from "../alerting/Toast";
import Select from "../form/Select";
import { formatRFC3339 } from "date-fns";

const ManageBill = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [bill, setBill] = useState();
    const [updatedBill, setUpdatedBill] = useState("");
    const [showUpdateForm, setShowUpdateForm] = useState("");

    const navigate = useNavigate();

    let { userId } = useParams();
    let { billId } = useParams();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 3 }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        console.log(`Attempting to update field ${name} to value ${value}`)
        setUpdatedBill({
            ...updatedBill,
            [name]: value,
        })
    }

    const handleDateChange = () => (event) => {
        let value = formatRFC3339(zonedTimeToUtc(event.target.value, 'America/New_York'), {fractionDigits: 3});
        let name = event.target.name;

        console.log(`Attempting to update field ${name} to value ${value}`)
        setUpdatedBill({
            ...updatedBill,
            [name]: value,
        })
    }

    const deleteBill = (event) => {
        event.preventDefault();
        console.log("Entered deleteBill Method")
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        const requestOptions = {
            method: "DELETE",
            headers: headers,
            credentials: "include",
        }

        fetch(`/users/${userId}/bills/${bill.id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    Toast("Error Deleting Bill", "error")
                } else {
                    Toast("Delete successful!", "success")
                    navigate(`/users/${userId}/bills`);
                }
            })
            .catch(error => {
                console.log(error.message)
                Toast("Error Deleting Loan", "error")
            })
    }

    const saveChanges = (event) => {
        event.preventDefault();
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        updatedBill.amount = parseFloat(updatedBill.amount)

        const requestOptions = {
            method: "PUT",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(updatedBill, null, 3),
        }

        fetch(`/users/${userId}/bills/${billId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast("An error occured while saving", "error")
                } else {
                    Toast("Save successful!", "success")
                    setShowUpdateForm(false);
                    fetchData();
                }
            })
            .catch(error => {
                Toast(error.message, "error")
            })
    }

    const fetchData = () => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/users/${userId}/bills/${billId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setBill(data);
                    setUpdatedBill(data);
                }
            })
            .catch(err => {
                console.log(err)
                Toast(err.message, "error")
            })
    };

    const toggleShowUpdateForm = () => {
        setShowUpdateForm(!showUpdateForm)
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData()

    }, [])

    return (
        <>
            <div className="row">
                <div className="col-md-9 offset-md-1">
                    <h2>Manage Bill</h2>
                </div>
                <div className="col-md-1">
                    <Link to={`/users/${userId}/bills`}><span className="badge bg-danger">Go back</span></Link>
                </div>
            </div>
            <div className="row">
                <div className="col-md-10 offset-md-1">
                    <hr />
                    {bill &&
                        <div className="chartContent">
                            <table className="table">
                                <thead>
                                    <tr>
                                        <th className="text-start">Name</th>
                                        <th className="text-start">Amount</th>
                                        
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr key={bill.id}>
                                        <td className="text-start">{bill.name}</td>
                                        <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(bill.amount)}</td>
                                    </tr>

                                </tbody>
                            </table>
                        </div>
                    }
                </div>
            </div>
            <div className="row">
                <div className="col-md-9 offset-md-1">
                    <Input
                        type="submit"
                        className="btn btn-success"
                        value={showUpdateForm ? "Cancel Edit" : "Edit"}
                        onClick={toggleShowUpdateForm}
                    />
                </div>

                <div className="col-md-1">
                    <Input
                        type="submit"
                        className="btn btn-danger"
                        value="Delete"
                        onClick={deleteBill}
                    />
                </div>
            </div>
            {showUpdateForm &&
                <div className="row">
                    <div className="col-md-10 offset-md-1">
                        <>
                            <h2>Edit Bill Details</h2>
                            <hr />
                            <form onSubmit={saveChanges}>
                                <input type="hidden" name="id" value={updatedBill.id}></input>
                                <Input
                                    title={"Name"}
                                    type={"text"}
                                    className={"form-control"}
                                    name={"name"}
                                    value={updatedBill.name}
                                    onChange={handleChange("")}
                                />
                                <Input
                                    title={"Amount (per month)"}
                                    type={"float"}
                                    className={"form-control"}
                                    name={"amount"}
                                    value={updatedBill.amount}
                                    onChange={handleChange("")}
                                />
                            </form>
                            <div className="col-md-8">
                                {showUpdateForm &&
                                    <Input
                                        type="submit"
                                        className="btn btn-success"
                                        value="Save Changes"
                                        onClick={saveChanges}
                                    />
                                }
                            </div>
                        </>
                    </div>
                </div>
            }
        </>
    )
}
export default ManageBill;