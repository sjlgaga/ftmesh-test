# ftmesh-dynamodb-test
## Data:(microseconds)
only request/       with attached state/      synchronize through separate request

0.5K(+0.5K)       3648.81               3787.34                     7984.52

1.5K(+0.5K)       3686.92               3751.70                     8380.76

10K(+0.5K)
        
10K(+10K)

## File Description
Data in `data.txt`, config are `yaml` files and `resource.go`, client is `client.go`

## Database Table Creation
``` bash
aws dynamodb create-table \
    --table-name movieTable \
    --attribute-definitions \
        AttributeName=title,AttributeType=S \
        AttributeName=year,AttributeType=N \
    --key-schema \
        AttributeName=title,KeyType=HASH \
        AttributeName=year,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=10 \
    --endpoint-url http://127.0.0.1:8000


aws dynamodb put-item \
    --table-name movieTable \
    --item '{
        "title": {"S": "FooBar"},
        "year": {"N": "2024"},
        "info": {
            "M": {
                "rating": {"S": "5.0"},
                "plot": {"S": "An epic saga of adventure and intrigue."}
            }
        }
    }' \
    --endpoint-url http://127.0.0.1:8000
```