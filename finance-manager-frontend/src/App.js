import { useCallback, useState } from "react";
import { Link, Outlet, useNavigate } from "react-router-dom";
import { ToastContainer, toast } from 'react-toastify';
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
  const [roles, setRoles] = useState([]);
  const navigate = useNavigate();

  const [numberFormatOptions, setNumberFormatOptions] = useState({ maximumFractionDigits: 2, minimumFractionDigits: 2 })
  const [interestFormatOptions, setInterestFormatOptions] = useState({ maximumFractionDigits: 3, minimumFractionDigits: 2 })

  const logOut = () => {
    setJwtToken("");
    setRoles("");
    setLoggedInUserId("");
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
    <div className="container-fluid">
      <div className="row min-vh-100">
        <div className="col-md-1 navMenu">
          <div className="navMenu">
            <Link to="/" className="list-group-item list-group-item-action"><p><HomeIcon/> Home</p></Link>
            <Link to="/about" className="list-group-item list-group-item-action"><p><Info/> About</p></Link>
            {jwtToken !== "" &&
              <>
                <Link to={`/users/${loggedInUserId}/loans`} className="list-group-item list-group-item-action"><p><BarChart /> Loans</p></Link>
                {roles.length > 0 && hasRole("admin") &&
                  <>
                    <Link to="/users" className="list-group-item list-group-item-action"><p><People/> Users</p></Link>
                  </>
                }
              </>
            }
          </div>
        </div>
        <div className="col-md-11">
          <div className="row">
          <div className="col text-end">
              {jwtToken === ""
                ? <Link to="/login"><span className="badge bg-success">Login</span></Link>
                : <a href="#!" onClick={logOut}><span className="badge bg-danger">Logout</span></a>
              }
            </div>
          </div>
          <div className="row">
            <div className="col-md-12">
              <h1 className="mt-3 text-center">Finance Manager</h1>
            </div>
            
            <hr className="mb-3" />
          </div>
          <div className="row">
            <div className="col-md-10">
            </div>
          </div>
          <Outlet context={{
            jwtToken,
            apiUrl,
            roles,
            loggedInUserId,
            numberFormatOptions,
            interestFormatOptions,
            setRoles,
            setJwtToken,
            setLoggedInUserId,
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
  )
}

export default App;
