curl --request POST \
  --url http://localhost:9999/v1/objects/bird \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"attrs": {"color": "red","specie": "aigle","height": 106,"owner": 1}}'

curl --request POST \
  --url http://localhost:9999/v1/objects/bird \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"attrs": {"color": "orange","specie": "aigle","height": 105,"owner": 1}}'

curl --request POST \
  --url http://localhost:9999/v1/objects/bird \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"attrs": {"color": "blue","specie": "mesange","height": 9,"owner": 1}}'



curl --request POST \
  --url http://localhost:9999/v1/objects/bird \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"attrs": {"color": "orange","specie": "bird with null field","height": 50,"weight": 5.2}}'

curl --request POST \
  --url http://localhost:9999/v1/objects/bird \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
  --data '{"attrs": {"color": "orange","specie": "bird with null field","height": 50,"weight": 10.8}}'

curl --request POST \
  --url http://localhost:9999/v1/objects/bird \
  --header 'content-type: application/json' \
  --header 'user-agent: vscode-restclient' \
    --data '{"attrs": {"color": "orange","specie": "bird with null field","height": 50,"weight": 100.8}}'