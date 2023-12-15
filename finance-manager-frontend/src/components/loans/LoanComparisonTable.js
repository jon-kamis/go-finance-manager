import { forwardRef, useEffect, useState } from "react";
import { GetDeltaMonthText, GetDeltaText, GetTextCompareClass } from "./LoanComparisonTools";
import { useOutletContext } from "react-router-dom";

const LoanComparisonTable = forwardRef((props, ref) => {

    const [loan, setLoan] = useState([]);
    const [updatedLoan, setUpdatedLoan] = useState([]);

    const numberFormatOptions = useOutletContext();
    const interestFormatOptions = useOutletContext();

    useEffect(() => {
        setLoan(props.loan)
    }, [props.loan])

    useEffect(() => {
        setUpdatedLoan(props.updatedLoan)
    }, [props.updatedLoan])

    return (
        <div className="d-flex">
            <div className="p-4 col-md-12 content">
            <h2>Comparing changes for Loan{loan.name && loan.name != "" ? ` ${loan.name}` : ""}</h2>
                <div className="content-tall-tablecontainer">
                    {loan && loan.name && updatedLoan && updatedLoan.name &&
                        <table className="table">
                            <thead>
                                <tr>
                                    <th>Name</th>
                                    <th className="text-start">Balance</th>
                                    <th className="text-start">Total Cost</th>
                                    <th className="text-start">Total Interest</th>
                                    <th className="text-start">Monthly Payment</th>
                                    <th className="text-start">Interest Rate</th>
                                    <th className="text-start">Loan Term</th>
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
                                <tr>
                                    <td>New Value</td>
                                    <td className={GetTextCompareClass(loan.total, updatedLoan.total)}>
                                        ${updatedLoan.total ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.total) : null}
                                    </td>
                                    <td className={GetTextCompareClass(loan.totalCost, updatedLoan.totalCost)}>
                                        ${updatedLoan.totalCost ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.totalCost) : null}
                                    </td>
                                    <td className={GetTextCompareClass(loan.interest, updatedLoan.interest)}>
                                        ${updatedLoan.interest ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.interest) : null}
                                    </td>
                                    <td className={GetTextCompareClass(loan.monthlyPayment, updatedLoan.monthlyPayment)}>
                                        ${updatedLoan.monthlyPayment ? Intl.NumberFormat("en-US", numberFormatOptions).format(updatedLoan.monthlyPayment) : null}
                                    </td>
                                    <td className={GetTextCompareClass(loan.interestRate, updatedLoan.interestRate)}>
                                        {updatedLoan.interestRate ? Intl.NumberFormat("en-US", interestFormatOptions).format(updatedLoan.interestRate) : null}
                                    </td>
                                    <td className={GetTextCompareClass(loan.loanTerm, updatedLoan.loanTerm)}>
                                        {updatedLoan.loanTerm ? updatedLoan.loanTerm : null}
                                    </td>
                                </tr>
                                <tr>
                                    <td>Delta</td>
                                    <td className={GetTextCompareClass(loan.total, updatedLoan.total)}>
                                        {GetDeltaText(loan.total, updatedLoan.total, numberFormatOptions, true)}
                                    </td>
                                    <td className={GetTextCompareClass(loan.totalCost, updatedLoan.totalCost)}>
                                        {GetDeltaText(loan.totalCost, updatedLoan.totalCost, numberFormatOptions, true)}
                                    </td>
                                    <td className={GetTextCompareClass(loan.interest, updatedLoan.interest)}>
                                        {GetDeltaText(loan.interest, updatedLoan.interest, numberFormatOptions, true)}
                                    </td>
                                    <td className={GetTextCompareClass(loan.monthlyPayment, updatedLoan.monthlyPayment)}>
                                        {GetDeltaText(loan.monthlyPayment, updatedLoan.monthlyPayment, numberFormatOptions, true)}
                                    </td>
                                    <td className={GetTextCompareClass(loan.interestRate, updatedLoan.interestRate)}>
                                        {GetDeltaText(loan.interestRate, updatedLoan.interestRate, interestFormatOptions, false)}
                                    </td>
                                    <td className={GetTextCompareClass(loan.loanTerm, updatedLoan.loanTerm, numberFormatOptions)}>
                                        {GetDeltaMonthText(loan.loanTerm, updatedLoan.loanTerm, numberFormatOptions, false)}
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    }
                </div>
            </div>
        </div>
    )
})

export default LoanComparisonTable;