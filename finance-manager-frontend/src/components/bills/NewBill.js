import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import { format, formatRFC3339, parse, parseISO } from "date-fns";
import { utcToZonedTime, zonedTimeToUtc } from "date-fns-tz";
import Select from "../form/Select";

const NewBill = () => {
    const { jwtToken } = useOutletContext();

    const [bill, setBill] = useState([]);

    const navigate = useNavigate();

    let { userId } = useParams();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 2 }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setBill({
            ...bill,
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

        bill.amount = parseFloat(bill.amount)

        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(bill, null, 3),
        }

        fetch(`/users/${userId}/bills`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast("An error occured while saving", "error")
                } else {
                    Toast("Save successful!", "success")
                    navigate(`/users/${userId}/bills`)
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

                    <h2>Create New Bill</h2>
                    <hr />
                    <form onSubmit={saveChanges}>
                        <input type="hidden" name="id" value="new"></input>
                        <Input
                            title={"Name"}
                            type={"text"}
                            className={"form-control"}
                            name={"name"}
                            value={bill.name}
                            onChange={handleChange("")}
                        />
                        <Input
                            title={"Amount (per month)"}
                            type={"float"}
                            className={"form-control"}
                            name={"amount"}
                            value={bill.amount}
                            onChange={handleChange("")}
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
export default NewBill;