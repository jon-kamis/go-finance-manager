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

Once installed, the app can be started with the following command from the directory this file is located in:

`docker compose up --build`

## Test Accounts Included by default
| Account | Role | Username | Password|
| --- | --- | --- | --- |
| Admin | Administrator | admin | admin|