import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import ManageLoan from "./ManageLoan";
import NewLoan from "./NewLoan";
import PaymentSchedule from "./PaymentSchedule";
import LoanComparisonTable from "./LoanComparisonTable";
import PaymentScheduleComparison from "./PaymentScheduleComparison";

const Loans = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [loan, setLoan] = useState([]);
    const [updatedLoan, setUpdatedLoan] = useState([]);
    const [loans, setLoans] = useState([]);
    const [selectedLoanId, setSelectedLoanId] = useState("");
    const [search, setSearch] = useState("");
    const [paymentSchedule, setPaymentSchedule] = useState();
    const [paymentScheduleTitle, setPaymentScheduleTitle] = useState();
    const [compare, setCompare] = useState(false);

    const navigate = useNavigate();

    const numberFormatOptions = useOutletContext();
    const interestFormatOptions = useOutletContext();

    let { userId } = useParams();

    function sortData(data) {
        let sortedData = data

        sortedData.sort((a, b) => a.name.toLowerCase() > b.name.toLowerCase() ? 1 : -1);

        return sortedData;
    }

    function fetchData() {

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/users/${userId}/loans`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setLoans(sortData(data));
                }
            })
            .catch(err => {
                console.log(err)
                Toast(err.message, "error")
            })

    }

    function updateSelectedId(id) {
        return () => {
            console.log(id)
            setSelectedLoanId(id)
            setCompare(false)
        }
    }

    const refreshData = () => (event) => {
        let value = event.target.value;
        setSearch(value)

        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        let searchUrl = ""

        value !== ""
            ? searchUrl = `?search=${value}`
            : searchUrl = ``

        fetch(`${apiUrl}/users/${userId}/loans${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setLoans(sortData(data));
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })
    }

    function fetchLoanById() {
        if (selectedLoanId && selectedLoanId !== "") {
            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)
            const requestOptions = {
                method: "GET",
                headers: headers,
            }

            fetch(`${apiUrl}/users/${userId}/loans/${selectedLoanId}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast(data.message, "error")
                    } else {
                        setLoan(data);
                    }
                })
                .catch(err => {
                    console.log(err)
                    Toast(err.message, "error")
                })
        }
    };

    function calcLoanPaymentSchedule() {
        if (loan && loan.id) {
            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)

            const requestOptions = {
                method: "POST",
                headers: headers,
                credentials: "include",
                body: JSON.stringify(loan, null, 3),
            }

            fetch(`/users/${userId}/loans/${selectedLoanId}/calculate`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        console.log(data.message)
                        Toast("An error occured during calculation", "error")
                    } else {
                        Toast("Success!", "success")
                        setPaymentSchedule(data.paymentSchedule)
                        setPaymentScheduleTitle(data.name)
                    }
                })
                .catch(error => {
                    console.log(error.message)
                    Toast("Unexpected error occured during calculation", "error")
                })
        }
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
        fetchLoanById();

    }, [selectedLoanId]);

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
        calcLoanPaymentSchedule();
    }, [loan]);

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData();

    }, []);

    return (
        <div className="container-fluid">
            <h1>Loans</h1>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content content-xtall">
                    <div className="row">
                        <div className="col-md-12">
                            <Input
                                title={"Search"}
                                type={"text"}
                                className={"form-control"}
                                name={"search"}
                                value={search}
                                onChange={refreshData("")}
                            />
                        </div>
                    </div>
                    <div className="content-xtall-tablecontainer">
                        <table className="table table-striped table-hover">

                            <thead>
                                <tr>
                                    <th className="text-start">Name</th>
                                    <th className="text-start">Principal</th>
                                    <th className="text-start">Rate</th>
                                    <th className="text-start">Term</th>
                                    <th className="text-start">Monthly Payment</th>
                                    <th className="text-start">Total Interest</th>
                                    <th className="text-start">Total Cost</th>
                                </tr>
                            </thead>
                            <tbody>
                                {loans.map((l) => (
                                    <>
                                        <tr key={l.id} onClick={updateSelectedId(l.id)}>
                                            <td className="text-start">{l.name}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.total)}</td>
                                            <td className="text-start">{Intl.NumberFormat("en-US", interestFormatOptions).format(l.interestRate)}</td>
                                            <td className="text-start">{l.loanTerm}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.monthlyPayment)}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.interest)}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(l.totalCost)}</td>
                                        </tr>
                                    </>
                                ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
            <div className="d-flex">
                <div className="p-4 col-md-6 content">
                    <NewLoan fetchData={fetchData} setPaymentSchedule={setPaymentSchedule} setPaymentScheduleTitle={setPaymentScheduleTitle} />
                </div>
                <div className="p-4 col-md-6 content">
                    <ManageLoan fetchData={fetchData} loanId={selectedLoanId} setLoanId={setSelectedLoanId} loan={loan} setLoan={setLoan} setUpdatedLoan={setUpdatedLoan} fetchLoanById={fetchLoanById} setCompare={setCompare} />
                </div>
            </div>
            {compare &&
                <>
                    <LoanComparisonTable
                        loan={loan}
                        updatedLoan={updatedLoan}
                        />
                    <PaymentScheduleComparison
                    loan={loan}
                    updatedLoan={updatedLoan}
                    userId={userId}
                    loanId={selectedLoanId} />
                </>
            }
            {!compare && paymentSchedule && paymentSchedule !== "" &&
                <PaymentSchedule 
                    schedule={paymentSchedule} 
                    title={paymentScheduleTitle} 
                />
            }

        </div>
    )
}
export default Loans;