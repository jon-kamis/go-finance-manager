import { forwardRef, useEffect, useState } from "react";
import Toast from "../alerting/Toast";
import { useNavigate, useOutletContext } from "react-router-dom";
import { GetDeltaTextForSingleValue, GetTextCompareClass } from "./LoanComparisonTools";

const PaymentScheduleComparison = forwardRef((props, ref) => {

    const [paymentSchedComparison, setPaymentSchedComparison] = useState([])
    const [loan, setLoan] = useState([])
    const [updatedLoan, setUpdatedLoan] = useState([])

    const { jwtToken } = useOutletContext();
    const { numberFormatOptions } = useOutletContext();

    const navigate = useNavigate();

    const fetchPaymentSchedComparison = () => {

        if (loan && loan.name && updatedLoan && updatedLoan.name) {
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

            fetch(`/users/${props.userId}/loans/${props.loanId}/compare-payments`, requestOptions)
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
        }
    };

    useEffect(() => {
        setLoan(props.loan)
    }, [props.loan])

    useEffect(() => {
        setUpdatedLoan(props.updatedLoan)
    }, [props.updatedLoan])

    useEffect(() => {
        fetchPaymentSchedComparison();
    }, [loan, updatedLoan])

    return (
        <div className="d-flex">
            <div className="p-4 col-md-12 content content-xtall">
                <h2>Payment Calendar for Loan{loan.name && loan.name != "" ? ` ${loan.name}` : ""}</h2>
                {paymentSchedComparison &&
                    <>
                        <hr />
                        <h2>Payment Schedule</h2>
                        <div className="content-xtall-tablecontainer">
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
                                        loan.paymentSchedule != null && paymentSchedComparison.length === 0 &&
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
                                                    <p className={GetTextCompareClass(p.principal, p.principalNew)}>
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalNew)}
                                                    </p>
                                                    <p className={GetTextCompareClass(p.principal, p.principalNew)}>
                                                        {GetDeltaTextForSingleValue(p.principalDelta, numberFormatOptions, true)}
                                                    </p>
                                                </td>
                                                <td>
                                                    <p className="text-end">
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interest)}
                                                    </p>
                                                    <p className={GetTextCompareClass(p.interest, p.interestNew)}>
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestNew)}
                                                    </p>
                                                    <p className={GetTextCompareClass(loan.interest, updatedLoan.interestNew)}>
                                                        {GetDeltaTextForSingleValue(p.interestDelta, numberFormatOptions, true)}
                                                    </p>
                                                </td>
                                                <td>
                                                    <p className="text-end">
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalToDate)}
                                                    </p>
                                                    <p className={GetTextCompareClass(p.principalToDate, p.principalToDateNew)}>
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.principalToDateNew)}
                                                    </p>
                                                    <p className={GetTextCompareClass(loan.principalToDate, updatedLoan.principalToDateNew)}>
                                                        {GetDeltaTextForSingleValue(p.principalToDateDelta, numberFormatOptions, true)}
                                                    </p>
                                                </td>
                                                <td>
                                                    <p className="text-end">
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestToDate)}
                                                    </p>
                                                    <p className={GetTextCompareClass(p.interestToDate, p.interestToDateNew)}>
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.interestToDateNew)}
                                                    </p>
                                                    <p className={GetTextCompareClass(loan.interestToDate, updatedLoan.interestToDateNew)}>
                                                        {GetDeltaTextForSingleValue(p.interestToDateDelta, numberFormatOptions, true)}
                                                    </p>
                                                </td>
                                                <td>
                                                    <p className="text-end">
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.remainingBalance)}
                                                    </p>
                                                    <p className={GetTextCompareClass(p.remainingBalance, p.remainingBalanceNew)}>
                                                        ${Intl.NumberFormat("en-US", numberFormatOptions).format(p.remainingBalanceNew)}
                                                    </p>
                                                    <p className={GetTextCompareClass(loan.remainingBalance, updatedLoan.remainingBalanceNew)}>
                                                        {GetDeltaTextForSingleValue(p.remainingBalanceDelta, numberFormatOptions, true)}
                                                    </p>
                                                </td>
                                            </tr>
                                        ))
                                    }
                                </tbody>
                            </table>
                        </div>

                    </>

                }
            </div>
        </div>
    )
});

export default PaymentScheduleComparison;