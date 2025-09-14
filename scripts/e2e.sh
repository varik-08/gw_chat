#!/usr/bin/env bash
set -euo pipefail

GATEWAY=${GATEWAY:-http://localhost:8080}

echo "# Login"
TOKENS=$(curl -s -X POST "$GATEWAY/api/auth/login" -H 'Content-Type: application/json' -d '{"username":"test","password":"test"}')
echo "$TOKENS" | jq . >/dev/null 2>&1 || { echo "Gateway not ready or jq missing"; exit 1; }
ACCESS=$(echo "$TOKENS" | jq -r .tokens.access_token)

echo "# Create chat"
CHAT=$(curl -s -X POST "$GATEWAY/api/chats" -H "Authorization: Bearer $ACCESS" -H 'Content-Type: application/json' -d '{"title":"demo","isPublic":true,"memberIds":[]}')
CHAT_ID=$(echo "$CHAT" | jq -r .chatId)
echo "Chat: $CHAT_ID"

echo "# Post message"
MSG=$(curl -s -X POST "$GATEWAY/api/messages" -H "Authorization: Bearer $ACCESS" -H 'Content-Type: application/json' -d '{"chatId":'$CHAT_ID',"senderId":1,"text":"hello"}')
echo "$MSG"

echo "# Get messages"
curl -s "$GATEWAY/api/chats/$CHAT_ID/messages" -H "Authorization: Bearer $ACCESS" | jq . || true

echo "OK"

