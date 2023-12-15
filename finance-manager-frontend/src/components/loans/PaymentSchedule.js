import { forwardRef } from "react";
import { useOutletContext } from "react-router-dom";

const PaymentSchedule = forwardRef((props, ref) => {

    const { numberFormatOptions } = useOutletContext();

return (
    <div className="d-flex">
                    <div className="p-4 col-md-12 content">
                        <h2>Payment Schedule{props.title && props.title != "" ? ` for ${props.title}` : ""}</h2>
                        <div className="content-xtall-tablecontainer">

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
                                    {props.schedule !== null &&
                                        props.schedule.map((p) => (
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
                </div>
)
});

export default PaymentSchedule;