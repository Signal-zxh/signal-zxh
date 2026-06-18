#!/bin/bash

BASE="http://localhost:8080"

echo "===================="
echo "1. LOGIN"
echo "===================="

RESP=$(curl -s $BASE/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123"}')

TOKEN=$(echo $RESP | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

echo "TOKEN: $TOKEN"

if [ -z "$TOKEN" ]; then
  echo "login failed"
  exit 1
fi

echo ""
echo "===================="
echo "2. CREATE POST"
echo "===================="

CREATE_RESP=$(curl -s -X POST $BASE/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"hello","content":"from script"}')

POST_ID=$(echo $CREATE_RESP | sed -n 's/.*"id":\([0-9]*\).*/\1/p')

echo "POST ID: $POST_ID"

echo ""
echo "===================="
echo "3. GET POSTS"
echo "===================="

curl -s $BASE/posts \
  -H "Authorization: Bearer $TOKEN" | head -c 200

echo ""

echo ""
echo "===================="
echo "4. GET BY ID"
echo "===================="

curl -s $BASE/posts/$POST_ID \
  -H "Authorization: Bearer $TOKEN"

echo ""

echo ""
echo "===================="
echo "5. UPDATE"
echo "===================="

curl -s -X PUT $BASE/posts/$POST_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"updated","content":"updated content"}'

echo ""
echo "update done"

echo ""
echo "===================="
echo "6. DELETE"
echo "===================="

curl -s -X DELETE $BASE/posts/$POST_ID \
  -H "Authorization: Bearer $TOKEN"

echo ""
echo "delete done"

echo ""
echo "ALL DONE"