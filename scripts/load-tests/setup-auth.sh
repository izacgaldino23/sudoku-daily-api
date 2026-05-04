#!/bin/sh
# setup-auth.sh - Gets JWT token for load testing protected endpoints

API_URL="${API_URL:-http://api-load:8080}"

# Wait for API to be ready
echo "Waiting for API to be ready..."
while ! wget -q --spider "$API_URL/health" 2>/dev/null; do
    sleep 1
done

# Register a test user (idempotent - will fail if exists, that's ok)
echo "Registering test user..."
curl -s -X POST "$API_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"loadtest@sudoku.test","username":"loadtest","password":"LoadTest123!"}' > /dev/null 2>&1

# Login to get token
echo "Logging in to get JWT token..."
RESPONSE=$(curl -s -X POST "$API_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"loadtest@sudoku.test","password":"LoadTest123!"}')

# Extract access token
TOKEN=$(echo "$RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "Failed to get JWT token. Response: $RESPONSE"
    exit 1
fi

echo "JWT Token obtained successfully"
echo "$TOKEN" > /scripts/.jwt_token

# Update token in body files
sed -i "s/TOKEN_PLACEHOLDER/${TOKEN}/g" /scripts/targets/bodies/submit-auth.json

# Create authenticated target files (using @/absolute/path.json syntax)
cat > /scripts/targets/submit-auth.txt <<EOF
POST http://api-load:8080/api/sudoku/submit
Authorization: Bearer ${TOKEN}
Content-Type: application/json

@/scripts/targets/bodies/submit-auth.json
EOF

cat > /scripts/targets/login.txt <<EOF
POST http://api-load:8080/api/auth/login
Content-Type: application/json

@/scripts/targets/bodies/login.json
EOF

echo "Auth setup complete. Token saved and target files created."
