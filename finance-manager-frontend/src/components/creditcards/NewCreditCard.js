import { forwardRef, useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const NewCreditCard = forwardRef((props, ref) => {
    const { jwtToken } = useOutletContext();

    const [creditCard, setCreditCard] = useState([]);

    const navigate = useNavigate();

    let { userId } = useParams();

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setCreditCard({
            ...creditCard,
            [name]: value,
        })
    }

    const saveCreditCard = (event) => {
        event.preventDefault();
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        creditCard.balance = parseFloat(creditCard.balance)
        creditCard.limit = parseFloat(creditCard.limit)
        creditCard.apr = parseFloat(creditCard.apr)
        creditCard.minPayment = parseFloat(creditCard.minPayment)
        creditCard.minPaymentPercentage = parseFloat(creditCard.minPaymentPercentage)

        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(creditCard, null, 3),
        }

        fetch(`/users/${userId}/credit-cards`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    Toast("Save successful!", "success")
                    props.fetchData()
                    creditCard.name = ""
                    creditCard.balance = 0
                    creditCard.limit = 0
                    creditCard.apr = 0
                    creditCard.minPayment = 0
                    creditCard.minPaymentPercentage = 0
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
            <div className="container-fluid">
                <h2>New Credit Card</h2>
                <div className="d-flex">
                    <div className="p-4 col-md-12">

                        <form onSubmit={saveCreditCard}>
                            <input type="hidden" name="id" value={creditCard.id}></input>
                            <Input
                                title={"Name"}
                                type={"text"}
                                className={"form-control"}
                                name={"name"}
                                value={creditCard.name}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Credit Limit"}
                                type={"number"}
                                className={"form-control"}
                                name={"limit"}
                                value={creditCard.limit}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Current Balance"}
                                type={"number"}
                                className={"form-control"}
                                name={"balance"}
                                value={creditCard.balance}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"APR"}
                                type={"number"}
                                className={"form-control"}
                                name={"apr"}
                                value={creditCard.apr}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Minimum Payment"}
                                type={"number"}
                                className={"form-control"}
                                name={"minPayment"}
                                value={creditCard.minPayment}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Minimum Payment Percentage"}
                                type={"number"}
                                className={"form-control"}
                                name={"minPaymentPercentage"}
                                value={creditCard.minPaymentPercentage}
                                onChange={handleChange("")}
                            />
                        </form>
                    </div>

                </div>
                <div className="d-flex">
                    <div className="p-4 col-md-12">
                        <Input
                            type="submit"
                            className="btn btn-primary"
                            value="Save"
                            onClick={saveCreditCard}
                        />
                    </div>
                </div>
            </div>
        </>
    )
})

export default NewCreditCard;