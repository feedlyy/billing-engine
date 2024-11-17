# billing-engine

This application is designed to manage loan billing with a special rule for first-time loaners: the first billing does not occur in the week the loan is requested but starts the following week.

## Getting Started


### Prerequisites
Before you begin, ensure you have:
- Go installed (version 1.16 or above is recommended).

### Installation
```azure
git clone [https://github.com/feedlyy/billing-engine.git]
cd billing-engine
```

### Running the application
```azure
go run main.go
```

### API Endpoints
All routes are defined in the router folder, particularly in the api.go file. Below is a brief overview:

#### Authentication
- Login: POST /login
    - Description: Endpoint to authenticate users.
    - Request form-data: { "username": "string", "password": "string" }
    - Response: JWT token for further authenticated requests.

#### Loans
- Check is Delinquent: GET /loan/check [admin role]
    - Description: For checking if the inputted user is delinquent or not
    - Request query param: {user: string}
    - Response: status user is delinquent or not
- Check is current outstanding: GET /loan/outstanding [customer role]
    - Description: For checking how many the current debt user have
    - Request: nil (get from current logged user) 
    - Response: int how many the current debt is (0 if it's closed or empty)
- Check Schedule Loan: GET /loan
    - Description: Check when user have to pay in week
    - Request: nil (get from current logged user)
    - Response: []string of list weeks
- Make Payment for loan: POST /loan/payment
    - Description: Pay for current week debt
    - Request form-data: {amount: string}
    - Response: empty if it's success, error if there's error (like already paid)

### Note
This application didn't rely on any database, so pre-defined dummy data are already exists in the project, on file:
```azure
util/helper.go
```