#! /usr/bin/env bash

# Parameters to set
function_name="addi64";
function_arguments="1, 2";
host="localhost:8080";

echo "SELECT $function_name($function_arguments) AS RESULT;"

curl -i -H "Content-Type: application/json" $host/q -d@- <<EndOfMessage
{"query": "SELECT $function_name($function_arguments) AS RESULT;"}
EndOfMessage

