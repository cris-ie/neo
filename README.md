# neo

Chart based application to query a weeks worth ov nasas neo api

# Requirements 
- linux (tested with ubuntu 20.10)
- build-essentials
- docker
- k3s
- helm

# Installation

1. Copy `values.yaml` to ./chart/
2. Run `make deploy` in .

The service should now be reachable under http://localhost

# Endpoints
The following endpoints are provided by the application: 

- http://localhost/status - reports that the server is ready to serve and a db connection is established
- http://localhost/liveness - reports that the server is ready to serve html files (a db connection might not yet be established or guaranteed)
- http://localhost/neo/week - reports the number of NEOs for the next week
- http://localhost/neo/next - reports the next NEO (Optional Query Parameter: ?hazardous=true - if set to true the next NEO that has the flag IsPotentiallyHazardousAstroid set to true is returned)

# Running the Application

The application can be started using the provided helm chart or locally

## Running the Application using helm

Assuming you have the prerequisites running run the following command to start the Application and its dependencies

`make build`

And after a succsessfull build run

`make deploy`

## Running the go application directly

Assuming you have a Postgresql datbase running you need to specify some environment variables for the application to run

- DB_NAME: The name of the database to use. E.g.: neo
- DB_HOST: The host where postgresql runs. E.g.: localhost
- DB_USER: A user name that has access rights to the Database. E.g.: postgres
- DB_PASSWORD: The `DB_USER`s password
- NASA_KEY: An API Key for NASA JPLs NEO Api