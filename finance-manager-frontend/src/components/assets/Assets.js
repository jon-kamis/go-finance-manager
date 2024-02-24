import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Toast from "../alerting/Toast";
import Input from "../form/Input";
import { format, formatRFC3339, parseISO } from "date-fns";
import { utcToZonedTime, zonedTimeToUtc } from "date-fns-tz";
import Select from "../form/Select";

const Assets = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [calcSavingsReq, setCalcSavingsReq] = useState([]);
    const [calcSavingsResp, setCalcSavingsResp] = useState([]);

    const navigate = useNavigate();

    let { userId } = useParams();

    const makeCalcSavingsReq = () => (event) => {
        event.preventDefault();
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        calcSavingsReq.goal = calcSavingsReq && calcSavingsReq.goal ? parseFloat(calcSavingsReq.goal) : 0
        calcSavingsReq.amount = calcSavingsReq && calcSavingsReq.amount ? parseFloat(calcSavingsReq.amount) : 0

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "POST",
            headers: headers,
            body: JSON.stringify(calcSavingsReq, null, 3),
        }

        fetch(`${apiUrl}/calc-savings`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setCalcSavingsResp(data);
                }
            })
            .catch(err => {
                console.log(err)
                Toast(err.message, "error")
            })

    }

    const handleChange = () => (event) => {
        let value = event.target.value;
        let name = event.target.name;
        console.log(`Attempting to update field ${name} to value ${value}`)
        setCalcSavingsReq({
            ...calcSavingsReq,
            [name]: value,
        })
    }

    const handleDateChange = () => (event) => {
        let value = formatRFC3339(zonedTimeToUtc(event.target.value, 'America/New_York'), { fractionDigits: 3 });
        let name = event.target.name;

        console.log(`Attempting to update field ${name} to value ${value}`)
        setCalcSavingsReq({
            ...calcSavingsReq,
            [name]: value,
        })
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
    }, []);

    useEffect(() => {
        setCalcSavingsResp([])
    }, [calcSavingsReq]);

    return (
        <div className="container-fluid">
            <h1>Assets</h1>
            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content">
                    <h2>Savings Calculator</h2>
                    <form onSubmit={makeCalcSavingsReq}>
                        <input type="hidden" name="id" value={calcSavingsReq.id}></input>
                        <div className="d-flex justify-content-around">
                            <div className="col-md-5">
                                <Input
                                    title={"Next Payday"}
                                    type={"date"}
                                    className={"form-control"}
                                    name={"nextPay"}
                                    value={calcSavingsReq && calcSavingsReq.nextPay ? format(utcToZonedTime(calcSavingsReq.nextPay, 'America/New_York'), 'yyyy-MM-dd', { timeZone: 'America/New_York' }) : ""}
                                    onChange={handleDateChange("")}
                                />
                            </div>
                            <div className="col-md-5">
                                <Input
                                    title={"Deadline"}
                                    type={"date"}
                                    className={"form-control"}
                                    name={"deadline"}
                                    value={calcSavingsReq && calcSavingsReq.deadline ? format(utcToZonedTime(calcSavingsReq.deadline, 'America/New_York'), 'yyyy-MM-dd', { timeZone: 'America/New_York' }) : ""}
                                    onChange={handleDateChange("")}
                                />
                            </div>
                        </div>
                        <div className="d-flex justify-content-around">
                            <div className="col-md-5">
                                <Input
                                    title={"Savings Goal"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"goal"}
                                    value={calcSavingsReq.goal}
                                    onChange={handleChange("")}
                                />
                            </div>
                            <div className="col-md-5">
                                <Input
                                    title={"Amount to Save"}
                                    type={"number"}
                                    className={"form-control"}
                                    name={"amount"}
                                    value={calcSavingsReq.amount}
                                    onChange={handleChange("")}
                                />
                            </div>
                        </div>
                        <div className="d-flex justify-content-around">
                            <div className="col-md-11">
                                <Select
                                    title={"Pay Frequency"}
                                    className={"form-control"}
                                    name={"payFrequency"}
                                    value={calcSavingsReq.payFrequency}
                                    onChange={handleChange("")}
                                    options={[{ id: "weekly", value: "weekly" }, { id: "bi-weekly", value: "bi-weekly" }, { id: "monthly", value: "monthly" }]}
                                    placeHolder={"Select"}
                                />
                            </div>
                        </div>
                        <div className="d-flex justify-content-around">
                            <Input
                                type="submit"
                                className="btn btn-primary"
                                value="Calculate"
                                onClick={makeCalcSavingsReq()}
                            />
                        </div>
                    </form>
                    {calcSavingsResp && calcSavingsResp.deadline &&
                        <>
                            <h2>Results</h2>
                            <h3>Deadline: {format(parseISO(calcSavingsResp.deadline), 'MMM do yyyy')}</h3>
                            <h3>Pays Before Deadline: {calcSavingsResp.numPays}</h3>
                            {calcSavingsResp.goal > 0 &&
                                <>
                                    <h3>Goal: ${calcSavingsResp.goal}</h3>
                                    <h3>Required savings per pay: ${calcSavingsResp.perPay}</h3>
                                </>
                            }
                            {calcSavingsResp.manPerPay > 0 &&
                                <>
                                    <h3>Manual Payments: ${calcSavingsResp.manPerPay}</h3>
                                    <h3>Amount saved: ${calcSavingsResp.actual}</h3>
                                </>
                            }

                        </>
                    }
                </div>
            </div>
        </div>
    )
}
export default Assets;