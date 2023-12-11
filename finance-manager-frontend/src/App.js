import { useCallback, useState } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import { ToastContainer, toast } from 'react-toastify';
import background from './images/main_background.jpg';
import Logo from './images/logo_v1.png';
import { NavData } from "./components/nav/NavData";
import HomeIcon from '@mui/icons-material/Home';
import People from '@mui/icons-material/People';
import Info from '@mui/icons-material/Info';
import BarChart from '@mui/icons-material/BarChart';
import 'react-toastify/dist/ReactToastify.min.css';
import './App.css'

function App() {
  const [jwtToken, setJwtToken] = useState("");
  const [apiUrl, setAPIUrl] = useState("http://localhost:8080")
  const [tickInterval, setTickInterval] = useState();
  const [loggedInUserId, setLoggedInUserId] = useState();
  const [loggedInUserName, setLoggedInUserName] = useState("");
  const [roles, setRoles] = useState([]);
  const navigate = useNavigate();

  const [numberFormatOptions, setNumberFormatOptions] = useState({ maximumFractionDigits: 2, minimumFractionDigits: 2 })
  const [interestFormatOptions, setInterestFormatOptions] = useState({ maximumFractionDigits: 3, minimumFractionDigits: 2 })

  const logOut = () => {
    setJwtToken("");
    setRoles("");
    setLoggedInUserId("");
    setLoggedInUserName("");
    navigate("/login")
  }

  const hasRole = (key) => {
    return roles.indexOf(key) !== -1;
  }

  const toggleRefresh = useCallback((status) => {
    console.log("clicked");

    if (status) {
      console.log("turning on ticking")
      let i = setInterval(() => {
        const requestOptions = {
          method: "GET",
          credentials: "include"
        }

        fetch(`/refresh`, requestOptions)
          .then((response) => response.json())
          .then((data) => {
            if (data.access_token) {
              setJwtToken(data.access_token)
            }
          })
          .catch(error => {
            console.log("user is not logged in", error)
          })

      }, 600000);
      setTickInterval(i)
      console.log("setting tick interval to", i);
    } else {
      console.log("turning off ticking");
      console.log("turning off tickInterval", tickInterval);
      setTickInterval(null);
      clearInterval(tickInterval);
    }
  }, [tickInterval])

  return (
    <div style={{ backgroundImage: `url(${background})` }}>
      <div className="container-fluid">
        <div className="row">
          <div className="col-md-2">
            <div className="appLogo">
              <img src={Logo} height="200%" width="150%" object-fit="contain" alt="logo"></img>
            </div>
          </div>
          <div className="col-md-8 offset-md-1 navMenu">
            <div className="d-flex justify-content-around">
              {NavData.map((n) => {
                {console.log(n)}
                return (!n.requiresJwt || jwtToken !== "") && (!n.requiresAdmin || hasRole("admin")) &&
                <>
                <div className="flex-col">
                  <Link to={n.path} className="list-group-item list-group-item-action"><p>{n.text}</p></Link>
                </div>
                </>
              })}
            </div>
          </div>
          <div className="col-md-1 text-end">
            {jwtToken === ""
              ? <Link to="/login"><span className="badge bg-success">Login</span></Link>
              : <a href="#!" onClick={logOut}><span className="badge bg-danger">Logout</span></a>
            }
          </div>
        </div>
        <div className="row min-vh-100">
          <div className="col-md-10 offset-md-1">
            <Outlet context={{
              jwtToken,
              apiUrl,
              loggedInUserName,
              roles,
              loggedInUserId,
              numberFormatOptions,
              interestFormatOptions,
              setRoles,
              setJwtToken,
              setLoggedInUserId,
              setLoggedInUserName,
              toggleRefresh,
              hasRole,
            }} />
            <ToastContainer
              position="bottom-center"
              autoClose={5000}
              hideProgressBar={false}
              newestOnTop={false}
              closeOnClick
              rtl={false}
              pauseOnFocusLoss
              draggable={false}
              pauseOnHover
              theme="light" />
          </div>
        </div>
      </div>
    </div>
  )
}

export default App;
