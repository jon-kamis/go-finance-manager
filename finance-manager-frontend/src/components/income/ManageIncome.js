import { forwardRef, useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, utcToZonedTime, zonedTimeToUtc } from 'date-fns-tz'

import Input from "../form/Input";
import Toast from "../alerting/Toast";
import Select from "../form/Select";
import { formatRFC3339 } from "date-fns";

const ManageIncome = forwardRef((props, ref) => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [income, setIncome] = useState([]);
    const [updatedIncome, setUpdatedIncome] = useState([]);
    const [showUpdateForm, setShowUpdateForm] = useState("");

    const navigate = useNavigate();

    let { userId } = useParams();

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
        let value = formatRFC3339(zonedTimeToUtc(event.target.value, 'America/New_York'), { fractionDigits: 3 });
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
                    props.fetchData()
                    props.setIncomeId(null)
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

        updatedIncome.rate = parseFloat(updatedIncome.rate)

        const requestOptions = {
            method: "PUT",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(updatedIncome, null, 3),
        }

        fetch(`/users/${userId}/incomes/${props.incomeId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast("An error occured while saving", "error")
                } else {
                    Toast("Save successful!", "success")
                    setShowUpdateForm(false);
                    props.fetchData()
                    fetchData()
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
        if (props.incomeId) {
            fetch(`${apiUrl}/users/${userId}/incomes/${props.incomeId}`, requestOptions)
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
        } else {
            income.name = ""
            income.type = ""
            income.hours = 0
            income.rate = 0
            income.frequency = ""
            income.taxPercentage = 0
            income.startDt = null
            setUpdatedIncome(income)
        }
    };

    const toggleShowUpdateForm = () => {
        setShowUpdateForm(!showUpdateForm)
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData()

    }, [props.incomeId]);

    return (
        <>
            <div className="container-fluid">
                <h2>Edit Income Details</h2>
                <div className="d-flex">
                    <div className="p-4 col-md-12">
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
                                title={"Hours per Pay"}
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
                                value={updatedIncome && updatedIncome.startDt ? format(utcToZonedTime(updatedIncome.startDt, 'America/New_York'), 'yyyy-MM-dd', { timeZone: 'America/New_York' }) : ""}
                                onChange={handleDateChange("")}
                            />
                        </form>
                    </div>
                </div>
                <div className="d-flex justify-content-between">
                    <div className="flex-col">
                        <Input
                            type="submit"
                            className="btn btn-primary"
                            value="Save Changes"
                            onClick={saveChanges}
                        />
                    </div>
                    <div className="flex-col">
                        <Input
                            type="submit"
                            className="btn btn-danger"
                            value="Delete"
                            onClick={deleteIncome}
                        />
                    </div>
                </div>
            </div >
        </>
    )
})

export default ManageIncome;