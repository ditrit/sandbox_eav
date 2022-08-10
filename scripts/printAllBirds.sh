curl --request GET \
  --url http://localhost:9999/v1/objects/bird/ \
  --header 'content-type: application/json' \
| jq .