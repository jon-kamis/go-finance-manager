import { Link, useOutletContext } from 'react-router-dom';
const About = () => {
    const { jwtToken } = useOutletContext()

    return (
        <>
            <div className="container-fluid">
                <div className="row">
                    <div className="col-md-8 offset-md-2 text-center">
                        <div className="content">
                            <h2>About</h2>
                            <hr />

                            <h3>
                                Learning Project
                            </h3>
                            <p>
                                Finance Manager is a side project developed as a way to learn new code.<br />
                                This was my first project using React as well as my first with GO. As such, it is not a refined application and does not accurately represent my typical quality of code.
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}
export default About;