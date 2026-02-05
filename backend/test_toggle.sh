#!/bin/bash
# 获取一个有效的 token
TOKEN=$(curl -s 'http://localhost:8084/api/v1/auth/wechat/callback?code=test&type=open' | jq -r '.data.token // empty')

if [ "$TOKEN" != "empty" ] && [ -n "$TOKEN" ]; then
  echo "Testing with valid token..."
  curl -s -X POST 'http://localhost:8084/api/v1/user-settings/toggle-collection' \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"enabled":false}'
else
  echo "No valid token available, testing without auth..."
  curl -s -X POST 'http://localhost:8084/api/v1/user-settings/toggle-collection' \
    -H 'Content-Type: application/json' \
    -d '{"enabled":false}'
fi
