#!/bin/bash
# Setup Clay pricing plans and features via Portcall API

set -e

# Read API key from .env file
if [ ! -f ".env" ]; then
  echo "Error: .env file not found. Please run 'npx portcall generate' first."
  exit 1
fi

API_KEY=$(grep '^PC_API_SECRET=' .env | cut -d'=' -f2 | tr -d '"' | tr -d "'")
if [ -z "$API_KEY" ]; then
  echo "Error: PC_API_SECRET not found in .env file. Please run 'npx portcall generate' first."
  exit 1
fi

API_URL="http://localhost:9080"

echo "Using API key from .env: ${API_KEY:0:20}..."
echo ""

# Helper function to make API calls
api() {
  local method=$1
  local endpoint=$2
  local data=$3
  
  if [ -n "$data" ]; then
    curl -s -X "$method" -H "x-api-key: $API_KEY" -H "Content-Type: application/json" \
      -d "$data" "$API_URL$endpoint"
  else
    curl -s -X "$method" -H "x-api-key: $API_KEY" "$API_URL$endpoint"
  fi
}

echo "=== Creating Features ==="

# Metered features (credits, users, people_company_searches)
echo "Creating metered features..."
api POST "/v1/features" '{"feature_id": "credits", "is_metered": true}' || true
api POST "/v1/features" '{"feature_id": "users", "is_metered": true}' || true
api POST "/v1/features" '{"feature_id": "people_company_searches", "is_metered": true}' || true

# On/off features
echo "Creating on/off features..."
FEATURES=(
  "sculptor"
  "sequencer"
  "exporting"
  "ai_claygent"
  "rollover_credits"
  "integration_providers"
  "scheduling"
  "phone_number_enrichments"
  "use_your_own_api_keys"
  "signals"
  "integrate_with_any_http_api"
  "webhooks"
  "email_sequencing_integrations"
  "exclude_people_company_filters"
  "web_intent"
  "crm_integrations"
)

for f in "${FEATURES[@]}"; do
  api POST "/v1/features" "{\"feature_id\": \"$f\", \"is_metered\": false}" 2>/dev/null || true
done

echo ""
echo "=== Creating Plans ==="

# Clay Free Plan - $0/mo, 100 credits
echo "Creating Clay Free plan..."
FREE_PLAN=$(api POST "/v1/plans" '{"name": "Clay Free", "currency": "USD", "interval": "month", "unit_amount": 0}')
FREE_PLAN_ID=$(echo "$FREE_PLAN" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "$FREE_PLAN" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "Free Plan ID: $FREE_PLAN_ID"

# Clay Starter Plan - $229/mo, 2000 credits
echo "Creating Clay Starter plan..."
STARTER_PLAN=$(api POST "/v1/plans" '{"name": "Clay Starter", "currency": "USD", "interval": "month", "unit_amount": 22900}')
STARTER_PLAN_ID=$(echo "$STARTER_PLAN" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "$STARTER_PLAN" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "Starter Plan ID: $STARTER_PLAN_ID"

# Clay Explorer Plan - $349/mo, 10000 credits
echo "Creating Clay Explorer plan..."
EXPLORER_PLAN=$(api POST "/v1/plans" '{"name": "Clay Explorer", "currency": "USD", "interval": "month", "unit_amount": 34900}')
EXPLORER_PLAN_ID=$(echo "$EXPLORER_PLAN" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "$EXPLORER_PLAN" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "Explorer Plan ID: $EXPLORER_PLAN_ID"

# Clay Pro Plan - $800/mo, 50000 credits
echo "Creating Clay Pro plan..."
PRO_PLAN=$(api POST "/v1/plans" '{"name": "Clay Pro", "currency": "USD", "interval": "month", "unit_amount": 80000}')
PRO_PLAN_ID=$(echo "$PRO_PLAN" | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['id'])" 2>/dev/null || echo "$PRO_PLAN" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "Pro Plan ID: $PRO_PLAN_ID"

echo ""
echo "=== Attaching Features to Plans ==="

# Function to add plan feature
add_feature() {
  local plan_id=$1
  local feature_id=$2
  local quota=$3
  local interval=${4:-"none"}
  
  if [ "$quota" = "null" ]; then
    api POST "/v1/plan-features" "{\"plan_id\": \"$plan_id\", \"feature_id\": \"$feature_id\", \"interval\": \"$interval\"}"
  else
    api POST "/v1/plan-features" "{\"plan_id\": \"$plan_id\", \"feature_id\": \"$feature_id\", \"quota\": $quota, \"interval\": \"$interval\"}"
  fi
}

# Free plan features: basic features + 100 credits
if [ -n "$FREE_PLAN_ID" ]; then
  echo "Adding features to Free plan..."
  add_feature "$FREE_PLAN_ID" "credits" 100 "month"
  add_feature "$FREE_PLAN_ID" "users" -1 "none"
  add_feature "$FREE_PLAN_ID" "people_company_searches" 100 "none"
  add_feature "$FREE_PLAN_ID" "sculptor" 1 "none"
  add_feature "$FREE_PLAN_ID" "sequencer" 1 "none"
  add_feature "$FREE_PLAN_ID" "exporting" 1 "none"
  add_feature "$FREE_PLAN_ID" "ai_claygent" 1 "none"
  add_feature "$FREE_PLAN_ID" "rollover_credits" 1 "none"
  add_feature "$FREE_PLAN_ID" "integration_providers" 1 "none"
fi

# Starter plan features: Free + scheduling, phone, api keys, 2000 credits
if [ -n "$STARTER_PLAN_ID" ]; then
  echo "Adding features to Starter plan..."
  add_feature "$STARTER_PLAN_ID" "credits" 2000 "month"
  add_feature "$STARTER_PLAN_ID" "users" -1 "none"
  add_feature "$STARTER_PLAN_ID" "people_company_searches" 5000 "none"
  add_feature "$STARTER_PLAN_ID" "sculptor" 1 "none"
  add_feature "$STARTER_PLAN_ID" "sequencer" 1 "none"
  add_feature "$STARTER_PLAN_ID" "exporting" 1 "none"
  add_feature "$STARTER_PLAN_ID" "ai_claygent" 1 "none"
  add_feature "$STARTER_PLAN_ID" "rollover_credits" 1 "none"
  add_feature "$STARTER_PLAN_ID" "integration_providers" 1 "none"
  add_feature "$STARTER_PLAN_ID" "scheduling" 1 "none"
  add_feature "$STARTER_PLAN_ID" "phone_number_enrichments" 1 "none"
  add_feature "$STARTER_PLAN_ID" "use_your_own_api_keys" 1 "none"
fi

# Explorer plan: Starter + HTTP API, webhooks, email seq, 10k credits
if [ -n "$EXPLORER_PLAN_ID" ]; then
  echo "Adding features to Explorer plan..."
  add_feature "$EXPLORER_PLAN_ID" "credits" 10000 "month"
  add_feature "$EXPLORER_PLAN_ID" "users" -1 "none"
  add_feature "$EXPLORER_PLAN_ID" "people_company_searches" 10000 "none"
  add_feature "$EXPLORER_PLAN_ID" "sculptor" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "sequencer" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "exporting" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "ai_claygent" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "rollover_credits" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "integration_providers" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "scheduling" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "phone_number_enrichments" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "use_your_own_api_keys" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "integrate_with_any_http_api" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "webhooks" 1 "none"
  add_feature "$EXPLORER_PLAN_ID" "email_sequencing_integrations" 1 "none"
fi

# Pro plan: Explorer + CRM integrations, 50k credits
if [ -n "$PRO_PLAN_ID" ]; then
  echo "Adding features to Pro plan..."
  add_feature "$PRO_PLAN_ID" "credits" 50000 "month"
  add_feature "$PRO_PLAN_ID" "users" -1 "none"
  add_feature "$PRO_PLAN_ID" "people_company_searches" 25000 "none"
  add_feature "$PRO_PLAN_ID" "sculptor" 1 "none"
  add_feature "$PRO_PLAN_ID" "sequencer" 1 "none"
  add_feature "$PRO_PLAN_ID" "exporting" 1 "none"
  add_feature "$PRO_PLAN_ID" "ai_claygent" 1 "none"
  add_feature "$PRO_PLAN_ID" "rollover_credits" 1 "none"
  add_feature "$PRO_PLAN_ID" "integration_providers" 1 "none"
  add_feature "$PRO_PLAN_ID" "scheduling" 1 "none"
  add_feature "$PRO_PLAN_ID" "phone_number_enrichments" 1 "none"
  add_feature "$PRO_PLAN_ID" "use_your_own_api_keys" 1 "none"
  add_feature "$PRO_PLAN_ID" "integrate_with_any_http_api" 1 "none"
  add_feature "$PRO_PLAN_ID" "webhooks" 1 "none"
  add_feature "$PRO_PLAN_ID" "email_sequencing_integrations" 1 "none"
  add_feature "$PRO_PLAN_ID" "crm_integrations" 1 "none"
fi

echo ""
echo "=== Publishing Plans ==="

# Publish all plans
if [ -n "$FREE_PLAN_ID" ]; then
  echo "Publishing Free plan..."
  api POST "/v1/plans/$FREE_PLAN_ID/publish" "{}"
fi

if [ -n "$STARTER_PLAN_ID" ]; then
  echo "Publishing Starter plan..."
  api POST "/v1/plans/$STARTER_PLAN_ID/publish" "{}"
fi

if [ -n "$EXPLORER_PLAN_ID" ]; then
  echo "Publishing Explorer plan..."
  api POST "/v1/plans/$EXPLORER_PLAN_ID/publish" "{}"
fi

if [ -n "$PRO_PLAN_ID" ]; then
  echo "Publishing Pro plan..."
  api POST "/v1/plans/$PRO_PLAN_ID/publish" "{}"
fi

echo ""
echo "=== Done! ==="
echo "Free Plan ID: $FREE_PLAN_ID"
echo "Starter Plan ID: $STARTER_PLAN_ID"
echo "Explorer Plan ID: $EXPLORER_PLAN_ID"
echo "Pro Plan ID: $PRO_PLAN_ID"
echo ""
echo "Run 'npx portcall generate' to regenerate the SDK client"
