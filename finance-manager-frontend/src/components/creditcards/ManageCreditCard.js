import { forwardRef, useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const ManageCreditCard = forwardRef((props, ref) => {
    const { jwtToken } = useOutletContext();

    const [creditCard, setCreditCard] = useState([]);
    const [updatedCreditCard, setUpdatedCreditCard] = useState([]);
    const navigate = useNavigate();

    let { userId } = useParams();

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setUpdatedCreditCard({
            ...updatedCreditCard,
            [name]: value,
        })
    }

    const deleteCreditCard = (event) => {
        event.preventDefault();
        console.log("Entered deleteCreditCard Method")
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

        fetch(`/users/${userId}/credit-cards/${creditCard.id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    Toast("Error Deleting CreditCard", "error")
                } else {
                    Toast("Delete successful!", "success")
                    props.fetchData();
                    props.setSelectedCreditCardId("");
                    props.creditCard.name = ""
                    props.creditCard.balance = 0
                    props.creditCard.limit = 0
                    props.creditCard.apr = 0
                    props.creditCard.minPayment = 0
                    props.creditCard.minPaymentPercentage = 0
                }
            })
            .catch(error => {
                console.log(error.message)
                Toast("Error Deleting CreditCard", "error")
            })
    }

    const saveChanges = (event) => {
        event.preventDefault();
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        if (updatedCreditCard && updatedCreditCard.id) {

            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)

            updatedCreditCard.balance = parseFloat(updatedCreditCard.balance)
            updatedCreditCard.limit = parseFloat(updatedCreditCard.limit)
            updatedCreditCard.apr = parseFloat(updatedCreditCard.apr)
            updatedCreditCard.minPayment = parseFloat(updatedCreditCard.minPayment)
            updatedCreditCard.minPaymentPercentage = parseFloat(updatedCreditCard.minPaymentPercentage)

            let creditCardToSave = updatedCreditCard
            const requestOptions = {
                method: "PUT",
                headers: headers,
                credentials: "include",
                body: JSON.stringify(creditCardToSave, null, 3),
            }

            fetch(`/users/${userId}/credit-cards/${creditCard.id}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast("An error occured while saving", "error")
                    } else {
                        Toast("Save successful!", "success")
                        props.fetchCreditCardById();
                        props.fetchData();
                    }
                })
                .catch(error => {
                    Toast(error.message, "error")
                })
        }
    }

    useEffect(() => {
        setCreditCard(props.creditCard)
        setUpdatedCreditCard(props.creditCard)
    }, [props.creditCard]);

    return (
        <div className="container-fluid">
            <h2>Manage CreditCard</h2>
            <div className="d-flex">
                <div className="col-md-12">
                    <form onSubmit={saveChanges}>
                        <input type="hidden" name="id" value={updatedCreditCard.id}></input>
                        <Input
                                title={"Name"}
                                type={"text"}
                                className={"form-control"}
                                name={"name"}
                                value={updatedCreditCard.name}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Credit Limit"}
                                type={"number"}
                                className={"form-control"}
                                name={"limit"}
                                value={updatedCreditCard.limit}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Current Balance"}
                                type={"number"}
                                className={"form-control"}
                                name={"balance"}
                                value={updatedCreditCard.balance}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"APR"}
                                type={"number"}
                                className={"form-control"}
                                name={"apr"}
                                value={updatedCreditCard.apr}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Minimum Payment"}
                                type={"number"}
                                className={"form-control"}
                                name={"minPayment"}
                                value={updatedCreditCard.minPayment}
                                onChange={handleChange("")}
                            />
                            <Input
                                title={"Minimum Payment Percentage"}
                                type={"number"}
                                className={"form-control"}
                                name={"minPaymentPercentage"}
                                value={updatedCreditCard.minPaymentPercentage}
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
                        value="Save"
                        onClick={saveChanges}
                    />
                </div>
                <div className="flex-col">
                    <Input
                        type="submit"
                        className="btn btn-danger"
                        value="Delete"
                        onClick={deleteCreditCard}
                    />
                </div>
            </div>
        </div>
    )
});

export default ManageCreditCard;