

### Requirement

* Elasticsearch 7.x
* Postgresql 9.4
* RabbitMQ

### Preparation

### Elastic Mapping
```json
curl --location --request PUT 'http://localhost:9200/news/_mapping' \
--header 'Accept: application/json' \
--data-raw '{
    "properties": {
        "id": {
            "type": "integer"
        },
        "created": {
        	"type": "date"
        }
    }
}'
```

## Elasticsearch Query

### Ordering by created and pagination

```json
{
  "from": 0, // offset
  "size": 20,
  "query": {
      "match_all": {}
  },
  "sort": [
    {
      "created": {
          "order": "desc"
      }
    }
  ]
}
```

## Example

### Create News

```json
curl --location --request POST 'https://kumparan.nias.dev/news' \
--header 'Content-Type: application/json' \
--data-raw '{
	"author": "Kumparan",
	"body": "Body"
}'
```

### Get News

```json
curl -XGET https://kumparan.nias.dev/news
```
