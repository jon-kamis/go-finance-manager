import { forwardRef, useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const NewLoan = forwardRef((props, ref) => {
    const { jwtToken } = useOutletContext();

    const [loan, setLoan] = useState([]);
    const [showResult, setShowResult] = useState(false);

    const navigate = useNavigate();

    let { userId } = useParams();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 2 }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setLoan({
            ...loan,
            [name]: value,
        })
    }

    const returnFromResult = (event) => {
        event.preventDefault();
        setShowResult(false)
    }

    const saveLoan = (event) => {
        event.preventDefault();
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        loan.loanTerm = parseFloat(loan.loanTerm)
        loan.total = parseFloat(loan.total)

        let loanToSave = loan
        loanToSave.paymentSchedule = []
        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(loan, null, 3),
        }

        fetch(`/users/${userId}/loans`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    Toast("Save successful!", "success")
                    setShowResult(false);
                    props.fetchData()
                    loan.name=""
                    loan.total=0
                    loan.interestRate=0.0
                    loan.loanTerm=0
                }
            })
            .catch(error => {
                Toast(error.message, "error")
            })
    }

    const handleSubmit = (event) => {
        event.preventDefault();

        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        loan.loanTerm = parseInt(`${loan.loanTerm}`)
        loan.total = parseFloat(`${loan.total}`)
        loan.interestRate = parseFloat(`${loan.interestRate}`)

        let loanToSave = loan
        loanToSave.paymentSchedule = []
        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(loanToSave, null, 3),
        }

        fetch(`/users/${userId}/loans/new/calculate`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    Toast("Success!", "success")
                    setLoan(data)
                    props.setPaymentSchedule(data.paymentSchedule)
                    props.setPaymentScheduleTitle("New Loan")
                    setShowResult(true);
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
                {!showResult
                    ?
                    <>
                        <h2>Calculate New Loan</h2>
                        <div className="d-flex">
                            <div className="p-4 col-md-12">

                                <form onSubmit={handleSubmit}>
                                    <input type="hidden" name="id" value={loan.id}></input>
                                    <Input
                                        title={"Total Amount Borrowed"}
                                        type={"number"}
                                        className={"form-control"}
                                        name={"total"}
                                        value={loan.total}
                                        onChange={handleChange("")}
                                    />
                                    <Input
                                        title={"Interest Rate (Percentage)"}
                                        type={"number"}
                                        className={"form-control"}
                                        name={"interestRate"}
                                        value={loan.interestRate}
                                        onChange={handleChange("")}
                                    />
                                    <Input
                                        title={"Loan Term"}
                                        type={"number"}
                                        className={"form-control"}
                                        name={"loanTerm"}
                                        value={loan.loanTerm}
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
                                    value="Calculate"
                                    onClick={handleSubmit}
                                />
                            </div>
                        </div>
                    </>
                    : <>
                        <h2>New Loan</h2>
                        <div className="d-flex">
                            <div className="col-md-12">
                                <table className="table">
                                    <thead>
                                        <tr>
                                            <th>Principal</th>
                                            <th>Rate</th>
                                            <th>Term</th>
                                            <th>Monthly Payment</th>
                                            <th>Total Interest</th>
                                            <th>Total Cost</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        <tr>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(loan.total)}</td>
                                            <td>{Intl.NumberFormat("en-US", interestFormatOptions).format(loan.interestRate)}</td>
                                            <td>{loan.loanTerm}</td>
                                            <td>{Intl.NumberFormat("en-US", numberFormatOptions).format(loan.monthlyPayment)}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(loan.interest)}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(loan.totalCost)}</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                        <div className="d-flex">
                            <div className="col-md-12">
                                <Input
                                    title={"Enter a name to save"}
                                    type={"text"}
                                    className={"form-control"}
                                    name={"name"}
                                    value={loan.name}
                                    onChange={handleChange("")}
                                />
                            </div>
                        </div>
                        <div className="d-flex justify-content-between">
                            <div className="flex-col">
                                <Input
                                    type="submit"
                                    className="btn btn-primary"
                                    value="Save"
                                    onClick={saveLoan}
                                />
                            </div>
                            <div className="flex-col">
                            <Input
                                type="submit"
                                className="btn btn-primary"
                                value="Go Back"
                                onClick={returnFromResult}
                            />
                            </div>
                        </div>
                    </>
                }
            </div>
        </>
    )
})

export default NewLoan;