# Invoice System

This project is the backend component of invoice system, implemented using Docker, Go, and Postgresql. It provides a basic backend and transaction functionality.


## Features

- Invoice Management

## Prerequisites

Make sure you have the following prerequisites installed:

- Docker (version >= 27.1.2)
- Docker Compose (version >= v1.29.2)

## Configuration

Modified Create ```app.yml``` file. ```app.yml``` supposed to be like this:

```
env: dev # dev, prod

server:
  port: :8081
  timezone: "Asia/Jakarta"
  read_timeout: 2
  write_timeout: 3

log:
  level: -4 # -4 DEBUG, 0 INFO, 4 WARN, 8 ERROR

db:
  user: postgres
  password: haris123
  host: 172.17.0.1
  port: 5432
  name: postgres
  pool:
    min: 2
    max: 5
    max_idle_time: 5m
    max_life_time: 1h


```

## How To Run

#### 1. Fun Way:


Thanks to docker-compose.

1. Clone repo

```
git clone git@github.com:Harisatul/invoice-test.git
```

2. Suppose you are in root folder. running docker compose up command:

```
docker-compose up -d
```

3. Execute sql migration file on sql folder. you can use go migrate:

```
migrate -source file://sql -database "postgres://db_username:db_password@db_host:sb_port/db_name?sslmode=disable" up
```

#### 2. Also Fun Way:

build binary

1. Clone repo

```
git clone git@github.com:Harisatul/invoice-test.git
```

2. Build Binary

```
go build .
```

3. Execute Binary

```
./invoice-test serve-http
```
note: service has argument ```serve-http```

4. Execute sql migration file on sql folder. you can use go migrate:

```
migrate -source file://sql -database "postgres://db_username:db_password@db_host:sb_port/db_name?sslmode=disable" up
```


## API USAGE

The backend component provides the API endpoints for Invoice System. To interact with the backend, you
can use an API testing tool such as Postman.

Import postman collection  to your Postman.

## API Documentation

### Health Check

- URL : ```/api/v1/health-check```
- Method : ```GET```
- success response:
```azure
{
    "status": 200,
    "message": "success",
    "data": "hello from server 👋"
}
```


### Create Invoice

- URL : ```/api/v1/invoice```
- Method : POST
- Request Body:

```azure
{
"customer_name": "John Doe",
"sales_person_name": "Jane Smith",
"payment_type": "CREDII",
"notes": "Thank you for your purchase!",
"product": [
    {
      "item_name": "Laptop",
      "quantity": 2,
      "total_cogs": 1000000,
      "total_price_sold": 1200000
    },
    {
      "item_name": "Wireless Mouse",
      "quantity": 1,
      "total_cogs": 200000,
      "total_price_sold": 250000
    }
  ]
}
```
 - success response
 ```azure
{
"status": 200,
"message": "success create invoice",
"data": "INV-7187-520093"
}
```

### Get Invoice

- URL : ```/api/v1/invoice```
- Method : ```GET```
- Request Params:
  - page: int (eg: 1) 
  - size: int (eg: 3)
  - start_date: YYYY-MM-DD (eg: 2024-11-26)
  - end_date: YYYY-MM-DD (eg: 2024-11-28)
  
- Success Response:
```azure
{
    "status": 200,
    "message": "success update invoice",
    "data": {
        "Invoice": [
            {
                "InvoiceNumber": "INV-7187-520093",
                "Date": "2024-11-30T00:00:00Z",
                "CustomerName": "John Doe",
                "Salesperson": "Jane Smith",
                "Notes": "Thank you for your purchase!",
                "PaymentType": "CREDIT",
                "Products": null
            }
        ],
        "InvoiceAggregateResponse": {
            "total_profit": 250000,
            "total_of_cash_transaction": 0
        }
    },
    "pagination_index": {
        "page": 1,
        "page_size": 2,
        "total_count": 1,
        "total_pages": 1,
        "has_previous": false,
        "has_next": false
    }
}
```

### Update Invoice

- URL : ```/api/invoice?id=INV-3912-923979```
- Method : ```PUT```
- Query Params:
    - id (invoice_number) : int (eg: INV-3912-923979)
- Request Body :
```azure
{
"customer_name": "John Doe",
"sales_person_name": "Jane Smith",
"payment_type": "CREDII",
"notes": "Thank you for your purchase!",
"product": [
    {
      "item_name": "Laptop",
      "quantity": 2,
      "total_cogs": 1000000,
      "total_price_sold": 1200000
    },
    {
      "item_name": "Wireless Mouse",
      "quantity": 1,
      "total_cogs": 200000,
      "total_price_sold": 250000
    }
  ]
}
```
- success response
 ```azure
{
"status": 200,
"message": "success update invoice",
"data": "INV-7187-520093"
}
```

### Delete Invoice

- URL : ```api/invoice?id=INV-5004-407761```
- Method : ```DELETE```
- Query Params:
    - id (invoice_number): int (eg: INV-5004-407761)
- success response
```azure
{
    "status": 200,
    "message": "success delete invoice"
}
```
- error response
```azure
{
    "status": 400,
    "message": "failed to delete invoice",
    "error": "given id not found"
}
```

### Import Invoice

- URL : ```/api/invoice/import```
- Method : ```POST```
- Form :
    - file :XLSX file with two sheet


