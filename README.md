# Welcome to Finance Manager

## About

Finance Manager is a side project developed as a way to learn new code.<br />
This was my first project using React as well as my first with GO. As such, it is not a refined application and does not accurately represent my typical quality of code.

This application is <b><u>not</u></b> intended for public use
<br/>
The sole intentions of this project are the following:
1. Provide me, the developer, an opportunity to learn new technologies
2. Provide an example of my work for future interviews or job opportunities
3. Provide a functioning finance management app for personal use

I am primarily a backend engineer and am trying to learn frontend development. I am absolutely not a designer and while I am proud of the website I am certainly not trying to show off my design skills

## Architecture
Finance Manager is comprised of three components

1. A backend appliation written in go
2. A frontend application written in React JS
3. A Postgres Database

All data is stored on the local file system by postgres

## Application Requirements
To run the application you must install docker
You can download it from the offical website <a href="https://docs.docker.com/engine/install/">here</a>

You will also need a small amount of space for postgres to store required files

## Running the Application

### Requirements
Finance Manager requires the following software
1. GO
2. NodeJS
3. Docker

### Running the full application
Once installed, the app can be started with the following command from the directory this file is located in:

`docker compose up --build`

The application can then be accessed at http://localhost:3000

### Running the application in individual layers
Alternatively, if you would like to build the application by individual layers, you can do the following

1. Start the Database layer by navigating into the finance-manager-backend directory and running
`docker compose -f docker-compose-db-only.yml up --build`
2. Start the backend GO layer by navigating into the finance-manager-backend directory and running
`go run .\cmd\finance-mngr\`
3. Start the frontend React application by navigating into the finance-manager-frontend directory and running
`npm start`

## Accounts Included by default
| Account | Role | Username | Password|
| --- | --- | --- | --- |
| Admin | Administrator | admin | admin|

## Testing the Application
### Backend Testing
The backend GO application has multiple unit and integration tests

These test can be run with the following commands
| Command | Description |
| --- | ---|
| go test .\cmd\finance-mngr\.. | Run Tests |
| go test -coverpkg=.\cmd\finance-mngr\... .\cmd\finance-mngr\... | Run Tests and Generate Coverage Reports |

`Note: The integration tests run in containers and as such require the docker engine to be running in order to work`

## API Documentation
### Viewing Documentation
This application uses swagger to provide API documentation
While the app is running, documentation can be accessed at `http://localhost:8080/swagger/index.html`

<b>Note: You may need to update the port in the above url if you started the API on a different port</b> 

### Updating or Generating Documentation
To update or generate the swagger documentation, first cd into the finance-mngr directory, then run the following command:

```powershell
swag init -g internal/handlers/fmhandler/FinanceManagerHandler.go
```

If the command fails you may need to first download the swaggo package. This can be done with the following commands

```powershell
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
```