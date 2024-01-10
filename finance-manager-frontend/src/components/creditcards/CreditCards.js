import { useEffect, useState } from "react";
import { useNavigate, useOutletContext, useParams } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";
import NewCreditCard from "./NewCreditCard";
import ManageCreditCard from "./ManageCreditCard";

const CreditCards = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [creditCard, setCreditCard] = useState([]);
    const [creditCards, setCreditCards] = useState([]);
    const [selectedCreditCardId, setSelectedCreditCardId] = useState("");
    const [search, setSearch] = useState("");

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

        fetch(`${apiUrl}/users/${userId}/credit-cards`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setCreditCards(sortData(data));
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
            setSelectedCreditCardId(id)
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

        fetch(`${apiUrl}/users/${userId}/credit-cards${searchUrl}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setCreditCards(sortData(data));
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })
    }

    function fetchCreditCardById() {
        if (selectedCreditCardId && selectedCreditCardId !== "") {
            const headers = new Headers();
            headers.append("Content-Type", "application/json")
            headers.append("Authorization", `Bearer ${jwtToken}`)
            const requestOptions = {
                method: "GET",
                headers: headers,
            }

            fetch(`${apiUrl}/users/${userId}/credit-cards/${selectedCreditCardId}`, requestOptions)
                .then((response) => response.json())
                .then((data) => {
                    if (data.error) {
                        Toast(data.message, "error")
                    } else {
                        setCreditCard(data);
                    }
                })
                .catch(err => {
                    console.log(err)
                    Toast(err.message, "error")
                })
        }
    };

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
        fetchCreditCardById();

    }, [selectedCreditCardId]);

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }
    }, [creditCard]);

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchData();

    }, []);

    return (
        <div className="container-fluid">
            <h1>Credit Cards</h1>
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
                                    <th className="text-start">Limit</th>
                                    <th className="text-start">Balance</th>
                                    <th className="text-start">APR</th>
                                    <th className="text-start">Min Payment</th>
                                    <th className="text-start">Min Payment %</th>
                                    <th className="text-start">Next Payment</th>
                                </tr>
                            </thead>
                            <tbody>
                                {creditCards.map((c) => (
                                    <>
                                        <tr key={c.id} onClick={updateSelectedId(c.id)}>
                                            <td className="text-start">{c.name}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(c.limit)}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(c.balance)}</td>
                                            <td className="text-start">{Intl.NumberFormat("en-US", interestFormatOptions).format(c.apr)}</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(c.minPayment)}</td>
                                            <td className="text-start">{Intl.NumberFormat("en-US", numberFormatOptions).format(c.minPaymentPercentage)}%</td>
                                            <td className="text-start">${Intl.NumberFormat("en-US", numberFormatOptions).format(c.payment)}</td>
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
                    <NewCreditCard fetchData={fetchData}/>
                </div>
                <div className="p-4 col-md-6 content">
                    <ManageCreditCard 
                        fetchData={fetchData}
                        fetchCreditCardById={fetchCreditCardById}
                        creditCardId={selectedCreditCardId}
                        setSelectedCreditCardId={setSelectedCreditCardId}
                        creditCard={creditCard}
                        setCreditCard={setCreditCard}/>
                </div>
            </div>
        </div>
    )
}
export default CreditCards;