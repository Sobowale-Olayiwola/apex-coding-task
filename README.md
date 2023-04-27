# SIMPLE DICE PROJECT

This project is focused on building a simple dice game REST API

## Assumptions made in the project
- All users are authenticated and must be registered on the system
- A holding account to credit new users wallet is added to the system which comes after successful DB migration is run.

## Tools required to run the project

- Golang Migrate a CLI tool is installed to run migrations (https://github.com/golang-migrate/migrate)
- Code generation tool sqlc is used to generate type safe code from SQL (https://github.com/kyleconroy/sqlc)

- A mocking framework for Golang (https://github.com/golang/mock)

- Docker, to pull postgres image

## How to run the project
 Follow the steps sequentially and the listed tools above should be installed to run the
 project successfully. The commands listed below should be run in the root folder of the project
 - make postgres
 - make createdb
 - make migrateup
 - make server

 ## To run the tests

 A part of the application is tested which involves the transfer of sats, creation of user 
 - make test

## Endpoints
A postman collection export is in the root of the application folder to be imported to easily interact with the endpoints
- POST /users
    - To create a new user
- POST /users/start-game
    - To start a new game
- PUT /users/:sessionID/end-game
    - To end a user session
- GET /users/:sessionID/check-session
    - To check if a user session is still active
- POST /users/roll-dice
    - To roll a dice in the game when a user have an active session
- GET /users/transaction-logs/:username
    - To get the transaction log for a particular user
- POST /transfers/fund-wallet
    - To fund the user wallet
- GET wallets/:username/balance
    - To get the user's wallet balance


