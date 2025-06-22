if [ -z "$SYSTEM_RENDER_API_KEY" ]; then
  echo "SYSTEM_RENDER_API_KEY is required"
  exit 1
fi

echo "Deploying to Render..."
curl -X POST https://api.render.com/v1/services/hub-service/deployments \
  -H "Authorization: Bearer $SYSTEM_RENDER_API_KEY"