#!/usr/bin/env bash

set -euo pipefail

datastore_base="http://localhost:8082/twirp/decode.iot.policystore.PolicyStore/"

function list_policies {
  echo "--> list policies"
  curl --request "POST" \
       --silent \
       --location "${datastore_base}ListEntitlementPolicies" \
       --header "Content-Type: application/json" \
       --data "{}" \
       | jq "."
}

function create_policy {
  echo "--> create policy"
  curl --request "POST" \
       --silent \
       --location "${datastore_base}CreateEntitlementPolicy" \
       --header "Content-Type: application/json" \
       --data "{\"public_key\":\"$1\",\"label\":\"first policy\",\"operations\":[{\"action\":\"SHARE\",\"sensor_id\":2}]}" \
       | jq "."
}

function delete_policy {
  echo "--> delete policy"
  curl --request "POST"
       --silent \
       --location "${datastore_base}DeleteEntitlementPolicy" \
       --header "Content-Type: application/json" \
       --data "{\"policy_id\": \"$1\", \"token\": \"$2\"}" \
       | jq "."
}

#list_policies
#delete_policy
#create_policy
"$@"