import { Link, useOutletContext } from 'react-router-dom';
import Ticket from './../images/movie_tickets.jpg';
const Home = () => {
    const { jwtToken } = useOutletContext()

    return (
        <>
            <div className="text-center">
                <h2>Home</h2>
                <hr />
                <h3>
                    Welcome to Finance Manager!
                </h3>
                {jwtToken !== ""
                   ? <>
                        <h3>Dashboard</h3>
                    </>
                : <>
                    <p>Please <Link to="/login">Login</Link> to view personalized dashboard</p>
                </>
                }
            </div>
        </>
    )
}
export default Home;