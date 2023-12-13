import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import { format, formatRFC3339, parse, parseISO } from "date-fns";
import { utcToZonedTime, zonedTimeToUtc } from "date-fns-tz";
import Select from "../form/Select";

const NewIncome = () => {
    const { jwtToken } = useOutletContext();

    const [income, setIncome] = useState([]);

    const navigate = useNavigate();

    let { userId } = useParams();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 2 }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setIncome({
            ...income,
            [name]: value,
        })
    }

    const handleDateChange = () => (event) => {
        let value = parse(event.target.value, 'yyyy-MM-dd', new Date(), {}).toISOString();
        let name = event.target.name;

        console.log(`Attempting to update field ${name} to value ${value}`)

        setIncome({
            ...income,
            [name]: value,
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

        income.taxPercentage = parseFloat(income.taxPercentage)
        income.hours = income.hours ? parseFloat(income.hours) : 0
        income.rate = parseFloat(income.rate)

        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(income, null, 3),
        }

        fetch(`/users/${userId}/incomes`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast("An error occured while saving", "error")
                } else {
                    Toast("Save successful!", "success")
                    navigate(`/users/${userId}/incomes`)
                }
            })
            .catch(error => {
                Toast(error.message, "error")
            })
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
        
    }, [])

    return (
        <>
            <div className="col-md-10 offset-md-1">
                <div className="row">

                    <h2>Create New Income</h2>
                    <hr />
                    <form onSubmit={saveChanges}>
                        <input type="hidden" name="id" value="new"></input>
                        <Input
                            title={"Name"}
                            type={"text"}
                            className={"form-control"}
                            name={"name"}
                            value={income.name}
                            onChange={handleChange("")}
                        />
                        <Select
                            title={"Type"}
                            className={"form-control"}
                            name={"type"}
                            value={income.type}
                            onChange={handleChange("")}
                            options={[{ id: "hourly", value: "hourly" }, { id: "salary", value: "salary" }]}
                            placeHolder={"Select"}
                        />
                        <Input
                            title={"Pay Rate"}
                            type={"float"}
                            className={"form-control"}
                            name={"rate"}
                            value={income.rate}
                            onChange={handleChange("")}
                        />
                        <Input
                            title={"Hours per Pay (Leave blank to calculate full-time)"}
                            type={"number"}
                            className={"form-control"}
                            name={"hours"}
                            value={income.hours}
                            onChange={handleChange("")}
                        />
                        <Select
                            title={"Pay Frequency"}
                            className={"form-control"}
                            name={"frequency"}
                            value={income.frequency}
                            onChange={handleChange("")}
                            options={[{ id: "weekly", value: "weekly" }, { id: "bi-weekly", value: "bi-weekly" }, { id: "monthly", value: "monthly" }]}
                            placeHolder={"Select"}
                        />
                        <Input
                            title={"Estimated Tax percentage"}
                            type={"float"}
                            className={"form-control"}
                            name={"taxPercentage"}
                            value={income.taxPercentage}
                            onChange={handleChange("")}
                        />
                        <Input
                            title={"Starting Date"}
                            type={"date"}
                            className={"form-control"}
                            name={"startDt"}
                            value={income.startDt ? format(parseISO(income.startDt), 'yyyy-MM-dd') : ""}
                            onChange={handleDateChange("")}
                        />
                        <Input
                            type="submit"
                            className="btn btn-primary"
                            value="Submit"
                        />
                    </form>
                </div>
            </div>
        </>
    )
}
export default NewIncome;