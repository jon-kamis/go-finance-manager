import { useEffect, useState } from "react";
import { Link, Outlet, useNavigate, useOutletContext, useParams } from "react-router-dom";
import { format, parseISO } from "date-fns";
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import NewIncome from "./NewIncome";
import ManageIncome from "./ManageIncome";

const Incomes = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();
    const { loggedInUserId } = useOutletContext();

    const [incomes, setIncomes] = useState([]);
    const [selectedIncomeId, setSelectedIncomeId] = useState()
    const [error, setError] = useState(false);
    const [search, setSearch] = useState("");

    const navigate = useNavigate();

    const numberFormatOptions = { maximumFractionDigits: 2, minimumFractionDigits: 2 }
    const interestFormatOptions = { maximumFractionDigits: 3, minimumFractionDigits: 2 }

    let { userId } = useParams();

    const handleChange = () => (event) => {
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
        {
            value !== ""
                ? searchUrl = `?search=${value}`
                : searchUrl = ``
        }

        fetch(`${apiUrl}/users/${userId}/incomes${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setIncomes(data);
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })
    }

    function fetchData() {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/users/${userId}/incomes`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setIncomes(data);
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
            setSelectedIncomeId(id)
        }
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData()
        setSelectedIncomeId(incomes != null && incomes.length > 0 ? incomes[0].id : null)
    }, []);

    return (
        <div className="container-fluid">
            <c-edit onrecordsaved={fetchData}></c-edit>

            <h1>Income Methods</h1>

            <div className="d-flex">
                <div className="p-4 flex-col col-md-12 content content-xtall">
                    <Input
                        title={"Search"}
                        type={"text"}
                        className={"form-control"}
                        name={"search"}
                        value={search}
                        onChange={handleChange("")}
                    />
                    <div className="content-xtall-tablecontainer">
                        <table className="table table-responsive table-striped table-hover">

                            <thead>
                                <tr>
                                    <th className="text-end">Name</th>
                                    <th className="text-end">Payment Type</th>
                                    <th className="text-end">Rate</th>
                                    <th className="text-end">Hours</th>
                                    <th className="text-end">Est. Gross Pay</th>
                                    <th classname="text-end">Est. Taxes</th>
                                    <th classname="text-end">Est. Net Pay</th>
                                    <th className="text-end">Frequency</th>
                                    <th className="text-end">Tax Percentage</th>
                                    <th className="text-end">Starting Date</th>
                                    <th className="text-end">Est. Next Date</th>
                                </tr>
                            </thead>
                            <tbody>
                                {incomes.map((i) => (
                                    <>
                                        <tr key={i.id} onClick={updateSelectedId(i.id)}>
                                            <td className="text-end">{i.name}</td>
                                            <td className="text-end">{i.type}</td>
                                            <td className="text-end">${Intl.NumberFormat("en-US", interestFormatOptions).format(i.rate)}</td>
                                            <td className="text-end">{Intl.NumberFormat("en-US", interestFormatOptions).format(i.hours)}</td>
                                            <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(i.grossPay)}</td>
                                            <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(i.taxes)}</td>
                                            <td className="text-end">${Intl.NumberFormat("en-US", numberFormatOptions).format(i.netPay)}</td>
                                            <td className="text-end">{i.frequency}</td>
                                            <td className="text-end">{Intl.NumberFormat("en-US", interestFormatOptions).format(i.taxPercentage)}</td>
                                            <td className="text-end">{format(parseISO(i.startDt), 'MMM do yyyy')}</td>
                                            <td className="text-end">{format(parseISO(i.nextDt), 'MMM do yyyy')}</td>
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
                    <NewIncome search={setSearch} data={setIncomes} />
                </div>
                <div className="p-4 col-md-6 content">
                    <ManageIncome fetchData={fetchData} incomeId={selectedIncomeId} setIncomeId = {setSelectedIncomeId}/>
                </div>
            </div>
        </div>
    )
}
export default Incomes;