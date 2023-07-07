# NOTE: at the time of writing, DROP FUNCTION is not supported.
# See https://seafowl.io/docs/guides/custom-udf-wasm#limitations
# 

filename="seafowl-udf-go.wasm"
function_name="addi64"
wasm_export="addi64"
return_type="BIGINT"
input_types=("BIGINT" "BIGINT")
host="localhost:8080"

IFS=", "
joined_input_types="${input_types[*]}"

curl -i -H "Content-Type: application/json" $host/q -d@- <<EOF
{"query": "DROP FUNCTION IF EXISTS $function_name($joined_input_types);"}
EOF