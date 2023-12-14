import { Link, useOutletContext } from 'react-router-dom';
const About = () => {
    const { jwtToken } = useOutletContext()

    return (
        <>
            <div className="container-fluid">
                <div className="d-flex">
                    <div className="flex-col col-md-12">
                        <h1>About</h1>
                        <div className="p-4 content">


                            <h2>
                                Learning Project
                            </h2>
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