import { forwardRef, useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, utcToZonedTime, zonedTimeToUtc } from 'date-fns-tz'

import Input from "../form/Input";
import Toast from "../alerting/Toast";
import { formatRFC3339 } from "date-fns";

const ManageBill = forwardRef((props, ref) => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [bill, setBill] = useState([]);

    const navigate = useNavigate();

    let { userId } = useParams();


    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        console.log(`Attempting to update field ${name} to value ${value}`)
        setBill({
            ...bill,
            [name]: value,
        })
    }

    const handleDateChange = () => (event) => {
        let value = formatRFC3339(zonedTimeToUtc(event.target.value, 'America/New_York'), { fractionDigits: 3 });
        let name = event.target.name;

        console.log(`Attempting to update field ${name} to value ${value}`)
        setBill({
            ...bill,
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
                    props.fetchData()
                    props.setBillId(null)
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

        bill.amount = parseFloat(bill.amount)

        const requestOptions = {
            method: "PUT",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(bill, null, 3),
        }

        fetch(`/users/${userId}/bills/${props.billId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast("An error occured while saving", "error")
                } else {
                    Toast("Save successful!", "success")
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
        if (props.billId) {
            fetch(`${apiUrl}/users/${userId}/bills/${props.billId}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast(data.message, "error")
                    } else {
                        setBill(data);
                        setBill(data);
                    }
                })
                .catch(err => {
                    console.log(err)
                    Toast(err.message, "error")
                })
        } else {
            bill.name=""
            bill.amount=0
            bill.id=""
        }
    };

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData()

    }, [props.billId])

    return (
        <>
            <div className="container-fluid">
                <h2>Edit Bill Details</h2>

                <div className="d-flex">
                    <div className="p-4 col-md-12">

                        <form onSubmit={saveChanges}>
                            <input type="hidden" name="id" value={bill.id}></input>
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
                            onClick={deleteBill}
                        />
                    </div>
                </div>
            </div>
        </>
    )
})
export default ManageBill;