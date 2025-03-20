# wyvern-api
Sample project for handle concurrent in Go

# API List
## Credit
`POST /api/transactions/credit`
### Request
    {
        "user_id": 1,
        "amount": 100000
    }
### Response

    {
        "code":200,
        "status":"success",
        "message":"",
        "data":{
            "transaction_id":20003,
            "new_balance":9850000
        }
    }
## Debit
`POST /api/transactions/debit`
### Request
    {
        "user_id": 1,
        "amount": 100000
    }
### Response

    {
        "code":200,
        "status":"success",
        "message":"",
        "data":{
            "transaction_id":20002,
            "new_balance":9850000
        }
    }

# Unit Test
`services > transaction_service_test.go`