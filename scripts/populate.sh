
ENDPOINT="http://localhost:9999/v1/objects/bird"
CONTENTTYPE="content-type: application/json"
DATAARRAY=('{"attrs":{"color":"red","specie":"aigle","height":10,"owner":1,"weight":5.2}}'\
    '{"attrs":{"color":"orange","specie":"mésange","height":6,"owner":1,"weight":10.2}}'\
    '{"attrs":{"color":"red","specie":"rouge_gorge","height":60,"owner":1,"weight":20.55}}'\
    '{"attrs":{"color":"orange","specie":"pinson","height":5,"owner":1,"weight":2.5}}'\
    '{"attrs":{"color":"red","specie":"aigle","height":1,"owner":2,"weight":25.1}}'\
    '{"attrs":{"color":"yellow","specie":"vautour","height":58,"owner":2,"weight":1.2}}'\
    '{"attrs":{"color":"blue","specie":"aigle","height":20,"owner":2,"weight":12525.112524211}}'\
    '{"attrs":{"color":"purple","specie":"mésange","height":21,"owner":2,"weight":12.12}}'\
'{"attrs":{"color":"blue","specie":"pinson","height":25785,"owner":2,"weight":3.2}}')


for data in ${DATAARRAY[@]}; do
    echo "Sending data: " $(echo $data | jq .)
    echo "Received: "
    curl --request POST --url $ENDPOINT --header CONTENTTYPE --data ${data} --silent | jq .
done