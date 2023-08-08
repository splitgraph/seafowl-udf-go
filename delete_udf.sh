# NOTE: As of August 2023, DROP FUNCTION support was added.

filename="seafowl-udf-go.wasm"
function_name="addi64"
wasm_export="addi64"
return_type="BIGINT"
input_types=("BIGINT" "BIGINT")
host="localhost:8080"

IFS=", "
joined_input_types="${input_types[*]}"

curl -i -H "Content-Type: application/json" $host/q -d@- <<EOF
{"query": "DROP FUNCTION $function_name;"}
EOF