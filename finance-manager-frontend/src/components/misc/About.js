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
                            <h1>Learning Project</h1>
                            <br />
                            <h2>
                                Finance Manager is a side project developed as a way to learn new code.<br />
                                This was my first project using React as well as my first with GO. As such, it is not a refined application and does not accurately represent my typical quality of code.
                            </h2>
                        </div>
                    </div>
                </div>
                <div className="d-flex">
                    <div className="p-4 flex-col col-md-12 content">
                        <h1>Intention</h1>
                        <br/>
                        <h2>This application is <b><u>not</u></b> intended for public use</h2>
                        <br/>
                        <h2>The sole intentions of this project are the following:</h2>
                        <h3 className="text-start">1. Provide me, the developer, an opportunity to learn new technologies</h3>
                        <h3 className="text-start">2. Provide an example of my work for future interviews or job opportunities</h3>
                        <h3 className="text-start">3. Provide a functioning finance management app for personal use</h3>
                    </div>
                </div>
                <div className="d-flex">
                    <div className="p-4 flex-col col-md-12 content">
                        <h1>Disclaimer</h1>
                        <br/>
                        <h2>I am primarily a backend engineer and am trying to learn frontend development. I am absolutely not a designer and did my best here</h2>
                    </div>
                </div>
            </div>
        </>
    )
}
export default About;