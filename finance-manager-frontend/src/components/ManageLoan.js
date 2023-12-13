import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "./form/Input";
import Toast from "./alerting/Toast";

const ManageLoan = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [loan, setLoan] = useState();
    const [updatedLoan, setUpdatedLoan] = useState("");
    const [paymentSchedComparison, setPaymentSchedComparison] = useState([])
    const [showUpdateForm, setShowUpdateForm] = useState(false);
    const [showCompareUpdate, setShowCompareUpdate] = useState(false);

    const navigate = useNavigate();

    let { userId } = useParams();
    let { loanId } = useParams();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 3 }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setShowCompareUpdate(false)
        setUpdatedLoan({
            ...updatedLoan,
            [name]: value,
        })
    }

    const toggleShowUpdateForm = () => {
        setShowUpdateForm(!showUpdateForm)
        setShowCompareUpdate(false)
    }

    const calcLoan = (event) => {
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

        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(loan, null, 3),
        }

        fetch(`/users/${userId}/loans/${loanId}/calculate`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.message)
                    Toast("An error occured during calculation", "error")
                } else {
                    Toast("Success!", "success")
                    setLoan(data)
                }
            })
            .catch(error => {
                console.log(error.message)
                Toast("Unexpected error occured during calculation", "error")
            })
    }

    const calcUpdateLoan = (event) => {
        event.preventDefault();

        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        updatedLoan.downPayment = parseFloat(`${updatedLoan.downPayment}`)
        updatedLoan.loanTerm = parseInt(`${updatedLoan.loanTerm}`)
        updatedLoan.total = parseFloat(`${updatedLoan.total}`)
        updatedLoan.interestRate = parseFloat(`${updatedLoan.interestRate}`)
        updatedLoan.monthlyPayment = 0

        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(updatedLoan, null, 3),
        }

        fetch(`/users/${userId}/loans/${loanId}/calculate`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.message)
                    Toast("An error occured during calculation", "error")
                } else {
                    Toast("Success!", "success")
                    setShowCompareUpdate(true)
                    setUpdatedLoan(data)
                    fetchPaymentSchedComparison()
                }
            })
            .catch(error => {
                console.log(error.message)
                Toast("Unexpected error occured during calculation", "error")
            })

    }

    const fetchPaymentSchedComparison = () => {

        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)

        const requestOptions = {
            method: "POST",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(updatedLoan, null, 3),
        }

        fetch(`/users/${userId}/loans/${loanId}/compare-payments`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.message)
                    Toast("An error occured during calculation", "error")
                } else {
                    Toast("Success!", "success")
                    setPaymentSchedComparison(data)
                }
            })
            .catch(error => {
                console.log(error.message)
                Toast("Unexpected error occured during calculation", "error")
            })
    };

    const deleteLoan = (event) => {
        event.preventDefault();
        console.log("Entered deleteLoan Method")
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

        fetch(`/users/${userId}/loans/${loan.id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error)
                    Toast("Error Deleting Loan", "error")
                } else {
                    Toast("Delete successful!", "success")
                    navigate(`/users/${userId}/loans`);
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

        updatedLoan.downPayment = parseFloat(updatedLoan.downPayment)
        updatedLoan.loanTerm = parseFloat(updatedLoan.loanTerm)
        updatedLoan.total = parseFloat(updatedLoan.total)

        let loanToSave = updatedLoan
        loanToSave.paymentSchedule = []
        const requestOptions = {
            method: "PUT",
            headers: headers,
            credentials: "include",
            body: JSON.stringify(loanToSave, null, 3),
        }

        fetch(`/users/${userId}/loans/${loanId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast("An error occured while saving", "error")
                } else {
                    Toast("Save successful!", "success")
                    setShowCompareUpdate(false);
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

        fetch(`${apiUrl}/users/${userId}/loans/${loanId}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setLoan(data);
                    setUpdatedLoan(data);
                }
            })
            .catch(err => {
                console.log(err)
                Toast(err.message, "error")
            })
    };

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData()

    }, [])

    const getTextCompareClass = (val1, val2) => {
        if (val1 > val2) {
            return "text-end text-success"
        } else if (val1 === val2) {
            return "text-end"
        } else {
            return "text-end text-danger"
        }
    };

    const getDeltaText = (val1, val2, formatOptions, addDollarSign) => {
        if (val1 > val2) {
            return `- ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val1 - val2))}`
        } else if (val1 === val2) {
            return "-"
        } else {
            return `+ ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val1 - val2))}`
        }
    }

    const getDeltaTextForSingleValue = (val, formatOptions, addDollarSign) => {
        if (val > 0) {
            return `+ ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val))}`
        } else if (val < 0) {
            return `- ${addDollarSign ? "$" : ""}${Intl.NumberFormat("en-US", formatOptions).format(Math.abs(val))}`
        } else {
            return "-"
        }
    }

    const getDeltaMonthText = (val1, val2) => {
        if (val1 > val2) {
            return `- ${val1 - val2}`
        } else if (val1 === val2) {
            return "-"
        } else {
            return `+ ${val1 - val2}`
        }
    }

    return (
        <>
            <div className="row">
                <div className="col-md-9 offset-md-1">
                    <h2>Manage Loan</h2>
                </div>
                <div className="col-md-1">
                    <Link to={`/users/${userId}/loans`}><span className="badge bg-danger">Go back</span></Link>
                </div>
            </div>
            <div className="row">
                <div className="col-md-10 offset-md-1">
                    <hr />
                    {loan &&
                        <div className="chartContent">
                            <table className="table">
                                <thead>
                                    <tr>
                                        <th>Name</th>
                                        <th className="text-end">Balance</th>
                                        <th className="text-end">Total Cost</th>
                                        <th className="text-end">Total Interest</th>
                                        <th className="text-end">Monthly Payment</th>
                                        <th className="text-end">Interest Rate</th>
                                        <th className="text-end">Loan Term</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr>
                                        <td>{loan.name ? loan.name : null}</td>
                                        <td className="text-end">${loan.total ? Intl.NumberFormat("en-US", numberFormatOptions).format(loan.total) : null}</td>
                                        <td className="text-end">${loan.totalCost ? Intl.NumberFormat("en-US", numberFormatOptions).format(loan.totalCost) : null}</td>
                                        <td className="text-end">${loan.interest ? Intl.NumberFormat("en-US", numberFormatOptions).format(loan.interest) : null}</td>
                                        <td className="text-end">${loan.monthlyPayment ? Intl.NumberFormat("en-US", numberFormatOptions).format(loan.monthlyPayment) : null}</td>
                                        <td className="text-end">{loan.interestRate ? Intl.NumberFormat("en-US", interestFormatOptions).format(loan.interestRate) : null}</td>
                                        <td className="text-end">{loan.loanTerm ? loan.loanTerm : null}</td>
                                    </tr>
                                    {showCompareUpdate &&
                                        <>
                                            <tr>
                                                <td>New Value</td>
                                                <td className={getTextCompareClass(loan.total, updatedLoan.total)}>
                                                    ${updatedLoan.total ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.total) : null}
                                                </td>
                                                <td className={getTextCompareClass(loan.totalCost, updatedLoan.totalCost)}>
                                                    ${updatedLoan.totalCost ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.totalCost) : null}
                                                </td>
                                                <td className={getTextCompareClass(loan.interest, updatedLoan.interest)}>
                                                    ${updatedLoan.interest ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.interest) : null}
                                                </td>
                                                <td className={getTextCompareClass(loan.monthlyPayment, updatedLoan.monthlyPayment)}>
                                                    ${updatedLoan.monthlyPayment ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.monthlyPayment) : null}
                                                </td>
                                                <td className={getTextCompareClass(loan.interestRate, updatedLoan.interestRate)}>
                                                    {updatedLoan.interestRate ? Intl.NumberFormat("en-US", interestFormatOptions).format(updatedLoan.interestRate) : null}
                                                </td>
                                                <td className={getTextCompareClass(loan.loanTerm, updatedLoan.loanTerm)}>
                                                    {updatedLoan.loanTerm ? updatedLoan.loanTerm : null}
                                                </td>
                                            </tr>
                                            <tr>
                                                <td>Delta</td>
                                                <td className={getTextCompareClass(loan.total, updatedLoan.total)}>
                                                    {getDeltaText(loan.total, updatedLoan.total, numberFormatOptions, true)}
                                                </td>
                                                <td className={getTextCompareClass(loan.totalCost, updatedLoan.totalCost)}>
                                                    {getDeltaText(loan.totalCost, updatedLoan.totalCost, numberFormatOptions, true)}
                                                </td>
                                                <td className={getTextCompareClass(loan.interest, updatedLoan.interest)}>
                                                    {getDeltaText(loan.interest, updatedLoan.interest, numberFormatOptions, true)}
                                                </td>
                                                <td className={getTextCompareClass(loan.monthlyPayment, updatedLoan.monthlyPayment)}>
                                                    {getDeltaText(loan.monthlyPayment, updatedLoan.monthlyPayment, numberFormatOptions, true)}
                                                </td>
                                                <td className={getTextCompareClass(loan.interestRate, updatedLoan.interestRate)}>
                                                    {getDeltaText(loan.interestRate, updatedLoan.interestRate, interestFormatOptions, false)}
                                                </td>
                                                <td className={getTextCompareClass(loan.loanTerm, updatedLoan.loanTerm, numberFormatOptions)}>
                                                    {getDeltaMonthText(loan.loanTerm, updatedLoan.loanTerm, numberFormatOptions, false)}
                                                </td>
                                            </tr>
                                        </>
                                    }
                                </tbody>
                            </table>
                        </div>
                    }
                </div>
            </div>
            <div className="row">
                <div className="col-md-1 offset-md-1">
                    <Input
                        type="submit"
                        className="btn btn-success"
                        value={showUpdateForm ? "Cancel Edit" : "Edit"}
                        onClick={toggleShowUpdateForm}
                    />
                </div>
                <div className="col-md-6">
                    {showCompareUpdate &&
                        <Input
                            type="submit"
                            className="btn btn-success"
                            value="Save Changes"
                            onClick={saveChanges}
                        />
                    }
                </div>
                <div className="col-md-2">
                    <Input
                        type="submit"
                        className="btn btn-primary"
                        value="View Payment Schedule"
                        onClick={calcLoan}
                    />
                </div>
                <div className="col-md-1">
                    <Input
                        type="submit"
                        className="btn btn-danger"
                        value="Delete"
                        onClick={deleteLoan}
                    />
                </div>
            </div>
            {showUpdateForm &&
                <div className="row">
                    <div className="col-md-10 offset-md-1">
                        <>
                            <h2>Edit Loan Details</h2>
                            <hr />
                            <form onSubmit={calcUpdateLoan}>
                                <input type="hidden" name="id" value={updatedLoan.id}></input>
                                <Input
                                    title={"Name"}
                                    type={"text"}
                                    className={"form-control"}
                                    name={"name"}
                                    value={updatedLoan.name}
                                    onChange={handleChange("")}
                                />
                                <Input
                                    title={"Total Amount Borrowed"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"total"}
                                    value={updatedLoan.total}
                                    onChange={handleChange("")}
                                />
                                <Input
                                    title={"Interest Rate (Percentage)"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"interestRate"}
                                    value={updatedLoan.interestRate}
                                    onChange={handleChange("")}
                                />
                                <Input
                                    title={"Loan Term"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"loanTerm"}
                                    value={updatedLoan.loanTerm}
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
                    </div>
                </div>
            }
            <div className="row">
                <div className="col-md-10 offset-md-1">
                    {loan && loan.paymentSchedule &&
                        <>
                            <div className="row">
                                <hr />
                                <h2>Payment Schedule</h2>
                                <div className="chartContent">
                                    <table className="table table-striped">
                                        <thead>
                                            <th>Month</th>
                                            <th className="text-end">Principal</th>
                                            <th className="text-end">Interest</th>
                                            <th className="text-end">PrincipalToDate</th>
                                            <th className="text-end">InterestToDate</th>
                                            <th className="text-end">Remaining Balance</th>
                                        </thead>
                                        <tbody>
                                            {
                                                loan.paymentSchedule !== null && paymentSchedComparison.length === 0 &&
                                                loan.paymentSchedule.map((p) => (
                                                    <tr key={p.id}>
                                                        <td>{p.month}</td>
                                                        <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principal)}</td>
                                                        <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interest)}</td>
                                                        <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalToDate)}</td>
                                                        <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestToDate)}</td>
                                                        <td>${Intl.NumberFormat("en-US", numberFormatOptions).format(p.remainingBalance)}</td>
                                                    </tr>
                                                ))
                                            }
                                            {
                                                paymentSchedComparison.length > 0 &&
                                                paymentSchedComparison.map((p) => (
                                                    <tr key={p.id}>
                                                        <td>{p.month}</td>
                                                        <td>
                                                            <p className="text-end">
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principal)}<br />
                                                            </p>
                                                            <p className={getTextCompareClass(p.principal, p.principalNew)}>
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalNew)}
                                                            </p>
                                                            <p className={getTextCompareClass(p.principal, p.principalNew)}>
                                                                {getDeltaTextForSingleValue(p.principalDelta, numberFormatOptions, true)}
                                                            </p>
                                                        </td>
                                                        <td>
                                                            <p className="text-end">
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interest)}
                                                            </p>
                                                            <p className={getTextCompareClass(p.interest, p.interestNew)}>
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestNew)}
                                                            </p>
                                                            <p className={getTextCompareClass(loan.interest, updatedLoan.interestNew)}>
                                                                {getDeltaTextForSingleValue(p.interestDelta, numberFormatOptions, true)}
                                                            </p>
                                                        </td>
                                                        <td>
                                                            <p className="text-end">
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalToDate)}
                                                            </p>
                                                            <p className={getTextCompareClass(p.principalToDate, p.principalToDateNew)}>
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalToDateNew)}
                                                            </p>
                                                            <p className={getTextCompareClass(loan.principalToDate, updatedLoan.principalToDateNew)}>
                                                                {getDeltaTextForSingleValue(p.principalToDateDelta, numberFormatOptions, true)}
                                                            </p>
                                                        </td>
                                                        <td>
                                                            <p className="text-end">
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestToDate)}
                                                            </p>
                                                            <p className={getTextCompareClass(p.interestToDate, p.interestToDateNew)}>
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestToDateNew)}
                                                            </p>
                                                            <p className={getTextCompareClass(loan.interestToDate, updatedLoan.interestToDateNew)}>
                                                                {getDeltaTextForSingleValue(p.interestToDateDelta, numberFormatOptions, true)}
                                                            </p>
                                                        </td>
                                                        <td>
                                                            <p className="text-end">
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.remainingBalance)}
                                                            </p>
                                                            <p className={getTextCompareClass(p.remainingBalance, p.remainingBalanceNew)}>
                                                                ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.remainingBalanceNew)}
                                                            </p>
                                                            <p className={getTextCompareClass(loan.remainingBalance, updatedLoan.remainingBalanceNew)}>
                                                                {getDeltaTextForSingleValue(p.remainingBalanceDelta, numberFormatOptions, true)}
                                                            </p>
                                                        </td>
                                                    </tr>
                                                ))
                                            }
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </>
                    }
                </div>
            </div>
        </>
    )
}
export default ManageLoan;