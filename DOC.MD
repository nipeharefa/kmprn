

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
