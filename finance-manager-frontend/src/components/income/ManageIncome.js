import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, utcToZonedTime, zonedTimeToUtc } from 'date-fns-tz'

import Input from "../form/Input";
import Toast from "../alerting/Toast";
import Select from "../form/Select";
import { formatRFC3339 } from "date-fns";

const ManageIncome = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [income, setIncome] = useState();
    const [updatedIncome, setUpdatedIncome] = useState("");
    const [showUpdateForm, setShowUpdateForm] = useState("");

    const navigate = useNavigate();

    let { userId } = useParams();
    let { incomeId } = useParams();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 3 }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        console.log(`Attempting to update field ${name} to value ${value}`)
        setUpdatedIncome({
            ...updatedIncome,
            [name]: value,
        })
    }

    const handleDateChange = () => (event) => {
        let value = formatRFC3339(zonedTimeToUtc(event.target.value, 'America/New_York'), {fractionDigits: 3});
        let name = event.target.name;

        console.log(`Attempting to update field ${name} to value ${value}`)
        setUpdatedIncome({
            ...updatedIncome,
            [name]: value,
        })
    }

    const deleteIncome = (event) => {
        event.preventDefault();
        console.log("Entered deleteIncome Method")
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

        fetch(`/users/${userId}/incomes/${income.id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    Toast("Error Deleting Income", "error")
                } else {
                    Toast("Delete successful!", "success")
                    navigate(`/users/${userId}/incomes`);
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

        updatedIncome.rate = parseFloat(updatedIncome.downPayment)
        updatedIncome.loanTerm = parseFloat(updatedIncome.loanTerm)
        updatedIncome.total = parseFloat(updatedIncome.total)

        const requestOptions = {
            method: "PUT",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(updatedIncome, null, 3),
        }

        fetch(`/users/${userId}/incomes/${incomeId}`, requestOptions)
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

        fetch(`${apiUrl}/users/${userId}/incomes/${incomeId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setIncome(data);
                    setUpdatedIncome(data);
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
                    <h2>Manage Loan</h2>
                </div>
                <div className="col-md-1">
                    <Link to={`/users/${userId}/incomes`}><span className="badge bg-danger">Go back</span></Link>
                </div>
            </div>
            <div className="row">
                <div className="col-md-10 offset-md-1">
                    <hr />
                    {income &&
                        <div className="chartContent">
                            <table className="table">
                                <thead>
                                    <tr>
                                        <th className="text-end">Name</th>
                                        <th className="text-end">Payment Type</th>
                                        <th className="text-end">Rate</th>
                                        <th className="text-end">Hours</th>
                                        <th className="text-end">Est. Gross Pay</th>
                                        <th classname="text-end">Est. Taxes</th>
                                        <th classname="text-end">Est. Net Pay</th>
                                        <th className="text-end">Frequency</th>
                                        <th className="text-end">Tax Percentage</th>
                                        <th className="text-end">Starting Date</th>
                                        <th className="text-end">Est. Next Date</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr key={income.id}>
                                        <td className="text-end">{income.name}</td>
                                        <td className="text-end">{income.type}</td>
                                        <td className="text-end">${Intl.NumberFormat("en-US", interestFormatOptions).format(income.rate)}</td>
                                        <td className="text-end">{Intl.NumberFormat("en-US", interestFormatOptions).format(income.hours)}</td>
                                        <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(income.grossPay)}</td>
                                        <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(income.taxes)}</td>
                                        <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(income.netPay)}</td>
                                        <td className="text-end">{income.frequency}</td>
                                        <td className="text-end">{Intl.NumberFormat("en-US", interestFormatOptions).format(income.taxPercentage)}</td>
                                        <td className="text-end">{format(utcToZonedTime(income.startDt, 'America/New_York'), 'MMM do yyyy', { timeZone: 'America/New_York' })}</td>
                                        <td className="text-end">{format(utcToZonedTime(income.nextDt, 'America/New_York'), 'MMM do yyyy', { timeZone: 'America/New_York' })}</td>
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
                        onClick={deleteIncome}
                    />
                </div>
            </div>
            {showUpdateForm &&
                <div className="row">
                    <div className="col-md-10 offset-md-1">
                        <>
                            <h2>Edit Income Details</h2>
                            <hr />
                            <form onSubmit={saveChanges}>
                                <input type="hidden" name="id" value={updatedIncome.id}></input>
                                <Input
                                    title={"Name"}
                                    type={"text"}
                                    className={"form-control"}
                                    name={"name"}
                                    value={updatedIncome.name}
                                    onChange={handleChange("")}
                                />
                                <Select
                                    title={"Type"}
                                    className={"form-control"}
                                    name={"type"}
                                    value={updatedIncome.type}
                                    onChange={handleChange("")}
                                    options={[{ id: "hourly", value: "hourly" }, { id: "salary", value: "salary" }]}
                                    placeHolder={"Select"}
                                />
                                <Input
                                    title={"Pay Rate"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"rate"}
                                    value={updatedIncome.rate}
                                    onChange={handleChange("")}
                                />
                                <Input
                                    title={"Hours per Pay (Set to 0 to calculate automatically based on a 40Hr work Week)"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"hours"}
                                    value={updatedIncome.hours}
                                    onChange={handleChange("")}
                                />
                                <Select
                                    title={"Pay Frequency"}
                                    className={"form-control"}
                                    name={"frequency"}
                                    value={updatedIncome.frequency}
                                    onChange={handleChange("")}
                                    options={[{ id: "weekly", value: "weekly" }, { id: "bi-weekly", value: "bi-weekly" }, { id: "monthly", value: "monthly" }]}
                                    placeHolder={"Select"}
                                />
                                <Input
                                    title={"Estimated Tax percentage"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"taxPercentage"}
                                    value={updatedIncome.taxPercentage}
                                    onChange={handleChange("")}
                                />
                                <Input
                                    title={"Starting Date"}
                                    type={"date"}
                                    className={"form-control"}
                                    name={"startDt"}
                                    value={format(utcToZonedTime(updatedIncome.startDt, 'America/New_York'), 'yyyy-MM-dd', { timeZone: 'America/New_York' })}
                                    onChange={handleDateChange("")}
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
export default ManageIncome;