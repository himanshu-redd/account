# Valuelabs Assignment

## Instructions to Install, Setup & Run

 ### Install postgres using docker
-  Goto project root directory and run command `docker compose up -d`. 
- Open Postgres CLI : `docker exec -it postgres_db psql -U root -d mydatabse`
- Create a new DB: ` create database bank`.
- Choose the created db: `\c bank`
- Create `accounts` table: 
    ```
    CREATE TABLE accounts (
        id BIGSERIAL PRIMARY KEY,
        balance DECIMAL(18,6) NOT NULL DEFAULT 0.00,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        deleted_at TIMESTAMPTZ
    );
    ```
### Run Project
- Goto project root directory and run `go mod tidy`
- Run project using: `go run cmd/main.go` 

### Endpoints

1. Create Account
    ```
    curl --location 'localhost:8080/accounts' \
    --header 'Content-Type: application/json' \
    --data '{
        "account_id" : 137224,
        "initial_balance" : "12398.76"
    }'
    ```
2. Get/Fetch Account
    ```
    curl --location 'localhost:8080/accounts/137224' \ --data ''
    ```
3. Perform a transaction
    ```
    curl --location 'localhost:8080/transactions' \
    --header 'Content-Type: application/json' \
    --data '{
        "source_account_id" : 137224,
        "destination_account_id" : 1334,
        "amount" : "1010924.7656"
    }'
    ```
