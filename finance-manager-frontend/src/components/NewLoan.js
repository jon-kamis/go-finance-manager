import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "./form/Input";
import Toast from "./alerting/Toast";

const NewLoan = () => {
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

        loan.downPayment = parseFloat(loan.downPayment)
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
                    setShowResult(true);
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

        loan.downPayment = parseFloat(`${loan.downPayment}`)
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
            <div className="col-md-9 offset-md-1">
                <div className="row">
                    {!showResult
                        ? <>
                            <h2>Create New Loan</h2>
                            <hr />
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
                                <hr />
                                <Input
                                    type="submit"
                                    className="btn btn-primary"
                                    value="Calculate"
                                />
                            </form>
                        </>
                        : <>
                            <h3>{"New Loan"}</h3>
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
                            <hr />
                            <div className="col-md-6">
                                <Input
                                    title={"Enter a name to save"}
                                    type={"text"}
                                    className={"form-control"}
                                    name={"name"}
                                    value={loan.name}
                                    onChange={handleChange("")}
                                />
                            </div>
                            <div className="col-md-1">
                                <Input
                                    type="submit"
                                    className="btn btn-primary"
                                    value="Save"
                                    onClick={saveLoan}
                                />
                            </div>
                            <div className="col-md-1">
                                <Input
                                    type="submit"
                                    className="btn btn-primary"
                                    value="Return"
                                    onClick={returnFromResult}
                                />
                            </div>
                        </>
                    }
                </div>
            </div>
            {showResult &&
                <div className="col-md-9 offset-md-1">
                    <div className="row">
                        <hr />
                        <h3>Payment Schedule</h3>
                        <table className="table table-striped">
                            <thead>
                                <th>Month</th>
                                <th>Principal</th>
                                <th>Interest</th>
                                <th>PrincipalToDate</th>
                                <th>InterestToDate</th>
                                <th>Remaining Balance</th>
                            </thead>
                            <tbody>
                                {loan.paymentSchedule !== null &&
                                    loan.paymentSchedule.map((p) => (
                                        <tr key={p.id}>
                                            <td>{p.month}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principal)}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interest)}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalToDate)}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestToDate)}</td>
                                            <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.remainingBalance)}</td>
                                        </tr>
                                    ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            }
        </>
    )
}
export default NewLoan;