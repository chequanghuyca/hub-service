#!/bin/bash

set -e

RENDER_API_KEY="$1"
SERVICE_ID="srv-cjfbgshod3po73f0hklg"

if [ -z "$RENDER_API_KEY" ]; then
  echo "Error: Render API Key not provided."
  exit 1
fi

echo "--- Starting deployment for service: ${SERVICE_ID} ---"

# Step 1: Trigger a new deploy and get the deploy ID
echo "1. Triggering a new deploy..."
DEPLOY_RESPONSE=$(curl --silent --request POST \
  --url "https://api.render.com/v1/services/${SERVICE_ID}/deploys" \
  --header "Accept: application/json" \
  --header "Authorization: Bearer ${RENDER_API_KEY}")

DEPLOY_ID=$(echo "$DEPLOY_RESPONSE" | jq --raw-output '.id')

if [ -z "$DEPLOY_ID" ] || [ "$DEPLOY_ID" == "null" ]; then
  echo "Error: Failed to trigger deploy. Response from Render:"
  echo "$DEPLOY_RESPONSE"
  exit 1
fi

echo "   -> Deploy triggered with ID: ${DEPLOY_ID}"

# Step 2: Loop and check the deployment status until it's finished
echo "2. Waiting for deployment to complete..."
while true; do
  STATUS_RESPONSE=$(curl --silent --request GET \
    --url "https://api.render.com/v1/services/${SERVICE_ID}/deploys/${DEPLOY_ID}" \
    --header "Accept: application/json" \
    --header "Authorization: Bearer ${RENDER_API_KEY}")
  
  STATUS=$(echo "$STATUS_RESPONSE" | jq --raw-output '.status')

  echo "   -> Current status: ${STATUS}"

  if [ "$STATUS" == "live" ]; then
    echo "--- ✅ Deployment successful! ---"
    exit 0
  elif [ "$STATUS" == "build_failed" ] || [ "$STATUS" == "canceled" ] || [ "$STATUS" == "deactivated" ]; then
    echo "--- ❌ Deployment failed with status: ${STATUS} ---"
    exit 1
  fi

  # Wait for 20 seconds before checking again
  sleep 20
done