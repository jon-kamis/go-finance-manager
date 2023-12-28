import { forwardRef, useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const ManageLoan = forwardRef((props, ref) => {
    const { jwtToken } = useOutletContext();

    const [loan, setLoan] = useState([]);
    const [updatedLoan, setUpdatedLoan] = useState([]);
    const navigate = useNavigate();

    let { userId } = useParams();

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        setUpdatedLoan({
            ...updatedLoan,
            [name]: value,
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

        fetch(`/users/${userId}/loans/${loan.id}/calculate`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.message)
                    Toast("An error occured during calculation", "error")
                } else {
                    Toast("Success!", "success")
                    setUpdatedLoan(data)
                    props.setUpdatedLoan(data)
                    props.setCompare(true)
                }
            })
            .catch(error => {
                console.log(error.message)
                Toast("Unexpected error occured during calculation", "error")
            })

    }

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
                    props.fetchData();
                    props.setSelectedLoanId("");
                    props.loan.id="";
                    props.loan.name="";
                    props.loan.total=0;
                    props.loan.interestRate=0;
                    props.loan.loanTerm=0;
                    props.setUpdatedLoan([]);
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

        if (updatedLoan && updatedLoan.id) {

            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)

            updatedLoan.loanTerm = parseFloat(updatedLoan.loanTerm)
            updatedLoan.total = parseFloat(updatedLoan.total)
            updatedLoan.interestRate = parseFloat(updatedLoan.interestRate)

            let loanToSave = updatedLoan
            loanToSave.paymentSchedule = []
            const requestOptions = {
                method: "PUT",
                headers: headers,
                credentials: "include",
                body: JSON.stringify(loanToSave, null, 3),
            }

            fetch(`/users/${userId}/loans/${loan.id}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast("An error occured while saving", "error")
                    } else {
                        Toast("Save successful!", "success")
                        props.setCompare(false);
                        props.fetchLoanById();
                        props.fetchData();
                    }
                })
                .catch(error => {
                    Toast(error.message, "error")
                })
        }
    }

    useEffect(() => {
        setLoan(props.loan)
        setUpdatedLoan(props.loan)
    }, [props.loan]);

    return (
        <div className="container-fluid">
            <h2>Manage Loan</h2>
            <div className="d-flex">
                <div className="col-md-12">
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
                    </form>
                </div>
            </div>
            <div className="d-flex justify-content-between">
                <div className="flex-col">
                    <Input
                        type="submit"
                        className="btn btn-primary"
                        value="Compare"
                        onClick={calcUpdateLoan}
                    />
                </div>

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
                        onClick={deleteLoan}
                    />
                </div>
            </div>
        </div>
    )
});

export default ManageLoan;