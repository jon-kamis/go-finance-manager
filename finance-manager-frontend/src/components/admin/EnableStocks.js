import { useEffect, useState } from "react";
import { Link, useNavigate, useOutletContext } from "react-router-dom";
import Input from "../form/Input";
import Toast from "../alerting/Toast";

const EnableStocks = () => {
    const { apiUrl } = useOutletContext();
    const { jwtToken } = useOutletContext();

    const [stocksEnabled, setStocksEnabled] = useState(false);
    const [apiKey, setApiKey] = useState("");

    const navigate = useNavigate();

    const handleChange = () => (event) => {
        let value = event.target.value;
        setApiKey(value)
    };

    const updateApiKey = () => (event) => {

        event.preventDefault()

        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        
        let reqBody = {}
        reqBody.key = apiKey
        
        const requestOptions = {
            method: "POST",
            headers: headers,
            body: JSON.stringify(reqBody, null, 3),
        }

        fetch(`${apiUrl}/modules/stocks/key`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setApiKey("")
                    Toast(data.message, "error")
                } else {
                    setApiKey("")
                    fetchIsStocksEnabled();
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })

    };

    function fetchIsStocksEnabled() {
        const headers = new Headers();
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", `Bearer ${jwtToken}`)
        const requestOptions = {
            method: "GET",
            headers: headers,
        }

        fetch(`${apiUrl}/modules/stocks`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    Toast(data.message, "error")
                } else {
                    setStocksEnabled(data.enabled);
                }
            })
            .catch(err => {
                Toast(err.message, "error")
                console.log(err)
            })
    }

    useEffect(() => {
        if (jwtToken === null || jwtToken === "") {
            navigate("/")
        }

        fetchIsStocksEnabled();

    }, []);

    return (
        <>
            <h2>Stocks Module</h2>
            <div className="row">
                <h3>{stocksEnabled ? <h3 className="text-success"><b>Enabled</b></h3> : <h3 className="text-failure"><b>Disabled</b></h3>}</h3>
                <div className="col-md-10">
                    <Input
                        title={"Api Key"}
                        type={"password"}
                        className={"form-control"}
                        name={"apiKey"}
                        value={apiKey}
                        onChange={handleChange("")}
                    />
                </div>
                <div className="col-md-2">
                    <Input
                        type="submit"
                        className="btn btn-primary"
                        value="Update"
                        onClick={updateApiKey()}
                    />
                </div>
            </div>
        </>
    )
}
export default EnableStocks;