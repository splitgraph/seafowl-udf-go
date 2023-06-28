#! /usr/bin/env bash
# from: https://dev.to/meleu/how-to-join-array-elements-in-a-bash-script-303a#join-elements-with-a-string
joinByString() {
  local separator="$1"
  shift
  local first="$1"
  shift
  printf "%s" "$first" "${@/#/$separator}"
}

filename="udf.wasm"
function_name="add_ints"
wasm_export="add_ints"
return_type="BIGINT"
input_types=("BIGINT" "BIGINT")

curl -i -H "Content-Type: application/json" $host/q -d@- <<EOF
{"query": "CREATE FUNCTION $function_name AS '{
  \"entrypoint\": \"$wasm_export\",
  \"language\": \"wasmMessagePack\",
  \"input_types\": [\"$(joinByString '\", \"' "${input_types[@]}")\"],
  \"return_type\": \"$return_type\",
  \"data\": \"$(base64 -i $filename)\"
}';"}
EOF