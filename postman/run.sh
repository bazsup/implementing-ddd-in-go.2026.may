#!/usr/bin/env bash
# Run DDD implementation scenarios against the API.
# Requirements: bash, curl, python3, awk — all pre-installed on macOS.
#
# Usage:
#   ./postman/run.sh                          # run all scenarios
#   ./postman/run.sh session-1-no-fractions   # run one scenario
#   BASE_URL=http://staging:8080 ./postman/run.sh
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
INVOICE_URL="${INVOICE_URL:-http://localhost:9000}"

PASS=0
FAIL=0

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

post_json() {
  local url="$1"
  local body="${2:-{}}"
  local response
  if ! response=$(curl -sf -X POST \
    -H "Content-Type: application/json" \
    -d "$body" \
    "$url" 2>&1); then
    echo "  ERROR: POST $url failed" >&2
    echo "  Response: $response" >&2
    exit 1
  fi
  echo "$response"
}

get_json() {
  local url="$1"
  local response
  if ! response=$(curl -sf -X GET "$url" 2>&1); then
    echo "  ERROR: GET $url failed" >&2
    echo "  Response: $response" >&2
    exit 1
  fi
  echo "$response"
}

# Extract a field from a JSON string. Supports dotted paths like "body.invoice_amount".
extract() {
  local field="$1"
  local json="$2"
  # Build python accessor from dot-notation: "body.invoice_amount" -> ['body']['invoice_amount']
  local accessor
  accessor=$(python3 -c "
parts = '${field}'.split('.')
print(''.join('[\"' + p + '\"]' for p in parts))
")
  echo "$json" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d${accessor})"
}

assert_field() {
  local label="$1"
  local actual="$2"
  local expected="$3"
  if [ "$actual" = "$expected" ]; then
    echo "  [PASS] $label"
    PASS=$((PASS + 1))
  else
    echo "  [FAIL] $label  expected='$expected'  actual='$actual'"
    FAIL=$((FAIL + 1))
  fi
}

assert_number() {
  local label="$1"
  local actual="$2"
  local expected="$3"
  if awk "BEGIN { exit !($actual == $expected) }"; then
    echo "  [PASS] $label"
    PASS=$((PASS + 1))
  else
    echo "  [FAIL] $label  expected=$expected  actual=$actual"
    FAIL=$((FAIL + 1))
  fi
}

# ---------------------------------------------------------------------------
# Scenarios
# ---------------------------------------------------------------------------

scenario_session_1_no_fractions() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [],
    "person_id": "Bald Eagle",
    "visit_id": "1"
  }')
  assert_number "price_amount"   "$(extract price_amount   "$resp")" 0
  assert_field  "person_id"      "$(extract person_id      "$resp")" "Bald Eagle"
  assert_field  "visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_1_some_fractions() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [
      {"amount_dropped": 15, "fraction_type": "Green waste"},
      {"amount_dropped": 39, "fraction_type": "Construction waste"}
    ],
    "person_id": "Bald Eagle",
    "visit_id": "1"
  }')
  assert_number "price_amount"   "$(extract price_amount   "$resp")" 7.35
  assert_field  "person_id"      "$(extract person_id      "$resp")" "Bald Eagle"
  assert_field  "visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_2_some_fractions_in_oak_city() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [
      {"amount_dropped": 83, "fraction_type": "Green waste"},
      {"amount_dropped": 18, "fraction_type": "Construction waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "1"
  }')
  assert_number "price_amount"   "$(extract price_amount   "$resp")" 10.06
  assert_field  "person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_2_5_additional_fee_when_3_deliveries_or_more_in_same_month() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [
      {"amount_dropped": 83, "fraction_type": "Green waste"},
      {"amount_dropped": 18, "fraction_type": "Construction waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "1"
  }')
  assert_number "visit1 price_amount"   "$(extract price_amount   "$resp")" 10.06
  assert_field  "visit1 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit1 visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "visit1 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-24",
    "dropped_fractions": [
      {"amount_dropped": 49, "fraction_type": "Construction waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "2"
  }')
  assert_number "visit2 price_amount"   "$(extract price_amount   "$resp")" 9.31
  assert_field  "visit2 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit2 visit_id"       "$(extract visit_id       "$resp")" "2"
  assert_field  "visit2 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-25",
    "dropped_fractions": [
      {"amount_dropped": 103, "fraction_type": "Green waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "3"
  }')
  assert_number "visit3 price_amount"   "$(extract price_amount   "$resp")" 8.65
  assert_field  "visit3 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit3 visit_id"       "$(extract visit_id       "$resp")" "3"
  assert_field  "visit3 price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_2_5_complex_scenario() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [
      {"amount_dropped": 83, "fraction_type": "Green waste"},
      {"amount_dropped": 18, "fraction_type": "Construction waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "1"
  }')
  assert_number "visit1 price_amount"   "$(extract price_amount   "$resp")" 10.06
  assert_field  "visit1 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit1 visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "visit1 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-24",
    "dropped_fractions": [
      {"amount_dropped": 49, "fraction_type": "Construction waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "2"
  }')
  assert_number "visit2 price_amount"   "$(extract price_amount   "$resp")" 9.31
  assert_field  "visit2 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit2 visit_id"       "$(extract visit_id       "$resp")" "2"
  assert_field  "visit2 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-24",
    "dropped_fractions": [
      {"amount_dropped": 10, "fraction_type": "Construction waste"}
    ],
    "person_id": "Bald Eagle",
    "visit_id": "3"
  }')
  assert_number "visit3 price_amount"   "$(extract price_amount   "$resp")" 1.5
  assert_field  "visit3 person_id"      "$(extract person_id      "$resp")" "Bald Eagle"
  assert_field  "visit3 visit_id"       "$(extract visit_id       "$resp")" "3"
  assert_field  "visit3 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-25",
    "dropped_fractions": [
      {"amount_dropped": 103, "fraction_type": "Green waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "4"
  }')
  assert_number "visit4 price_amount"   "$(extract price_amount   "$resp")" 8.65
  assert_field  "visit4 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit4 visit_id"       "$(extract visit_id       "$resp")" "4"
  assert_field  "visit4 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-08-12",
    "dropped_fractions": [
      {"amount_dropped": 33, "fraction_type": "Green waste"}
    ],
    "person_id": "Squirrel Gus",
    "visit_id": "5"
  }')
  assert_number "visit5 price_amount"   "$(extract price_amount   "$resp")" 2.64
  assert_field  "visit5 person_id"      "$(extract person_id      "$resp")" "Squirrel Gus"
  assert_field  "visit5 visit_id"       "$(extract visit_id       "$resp")" "5"
  assert_field  "visit5 price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_3_business_customers() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [
      {"amount_dropped": 83, "fraction_type": "Green waste"},
      {"amount_dropped": 18, "fraction_type": "Construction waste"}
    ],
    "person_id": "Beaver Bertha",
    "visit_id": "1"
  }')
  assert_number "visit1 price_amount"   "$(extract price_amount   "$resp")" 10.42
  assert_field  "visit1 person_id"      "$(extract person_id      "$resp")" "Beaver Bertha"
  assert_field  "visit1 visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "visit1 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-09-30",
    "dropped_fractions": [
      {"amount_dropped": 134, "fraction_type": "Green waste"},
      {"amount_dropped": 201, "fraction_type": "Construction waste"}
    ],
    "person_id": "Bear Billy",
    "visit_id": "2"
  }')
  assert_number "visit2 price_amount"   "$(extract price_amount   "$resp")" 42.21
  assert_field  "visit2 person_id"      "$(extract person_id      "$resp")" "Bear Billy"
  assert_field  "visit2 visit_id"       "$(extract visit_id       "$resp")" "2"
  assert_field  "visit2 price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_3_tier_based_pricing() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [{"amount_dropped": 597, "fraction_type": "Construction waste"}],
    "person_id": "Beaver Bertha",
    "visit_id": "1"
  }')
  assert_number "visit1 price_amount"   "$(extract price_amount   "$resp")" 125.37
  assert_field  "visit1 person_id"      "$(extract person_id      "$resp")" "Beaver Bertha"
  assert_field  "visit1 visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "visit1 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [{"amount_dropped": 1803, "fraction_type": "Construction waste"}],
    "person_id": "Beaver Bertha",
    "visit_id": "2"
  }')
  assert_number "visit2 price_amount"   "$(extract price_amount   "$resp")" 490.63
  assert_field  "visit2 person_id"      "$(extract person_id      "$resp")" "Beaver Bertha"
  assert_field  "visit2 visit_id"       "$(extract visit_id       "$resp")" "2"
  assert_field  "visit2 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [{"amount_dropped": 1228, "fraction_type": "Construction waste"}],
    "person_id": "Peppa Python",
    "visit_id": "3"
  }')
  assert_number "visit3 price_amount"   "$(extract price_amount   "$resp")" 276.12
  assert_field  "visit3 person_id"      "$(extract person_id      "$resp")" "Peppa Python"
  assert_field  "visit3 visit_id"       "$(extract visit_id       "$resp")" "3"
  assert_field  "visit3 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-25",
    "dropped_fractions": [{"amount_dropped": 901, "fraction_type": "Construction waste"}],
    "person_id": "Beaver Bertha",
    "visit_id": "4"
  }')
  assert_number "visit4 price_amount"   "$(extract price_amount   "$resp")" 261.29
  assert_field  "visit4 person_id"      "$(extract person_id      "$resp")" "Beaver Bertha"
  assert_field  "visit4 visit_id"       "$(extract visit_id       "$resp")" "4"
  assert_field  "visit4 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2024-01-15",
    "dropped_fractions": [{"amount_dropped": 600, "fraction_type": "Construction waste"}],
    "person_id": "Beaver Bertha",
    "visit_id": "5"
  }')
  assert_number "visit5 price_amount"   "$(extract price_amount   "$resp")" 126
  assert_field  "visit5 person_id"      "$(extract person_id      "$resp")" "Beaver Bertha"
  assert_field  "visit5 visit_id"       "$(extract visit_id       "$resp")" "5"
  assert_field  "visit5 price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_4_business_exemption_per_business() {
  post_json "$BASE_URL/startScenario" > /dev/null

  local resp
  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [{"amount_dropped": 597, "fraction_type": "Construction waste"}],
    "person_id": "Beaver Bertha",
    "visit_id": "1"
  }')
  assert_number "visit1 price_amount"   "$(extract price_amount   "$resp")" 125.37
  assert_field  "visit1 person_id"      "$(extract person_id      "$resp")" "Beaver Bertha"
  assert_field  "visit1 visit_id"       "$(extract visit_id       "$resp")" "1"
  assert_field  "visit1 price_currency" "$(extract price_currency "$resp")" "USD"

  resp=$(post_json "$BASE_URL/calculatePrice" '{
    "date": "2023-07-23",
    "dropped_fractions": [{"amount_dropped": 1803, "fraction_type": "Construction waste"}],
    "person_id": "Beaver Bruce",
    "visit_id": "2"
  }')
  assert_number "visit2 price_amount"   "$(extract price_amount   "$resp")" 490.63
  assert_field  "visit2 person_id"      "$(extract person_id      "$resp")" "Beaver Bruce"
  assert_field  "visit2 visit_id"       "$(extract visit_id       "$resp")" "2"
  assert_field  "visit2 price_currency" "$(extract price_currency "$resp")" "USD"
}

scenario_session_5_send_invoice_after_payment() {
  post_json "$INVOICE_URL/api/invoice-inspector/reset" > /dev/null
  post_json "$BASE_URL/startScenario" > /dev/null

  post_json "$INVOICE_URL/api/invoice" '{
    "email": "beavers@dam-building.com",
    "invoice_amount": 125.37,
    "invoice_currency": "USD"
  }' > /dev/null

  local resp
  resp=$(get_json "$INVOICE_URL/api/invoice-inspector/beavers@dam-building.com")
  assert_number "invoice_amount"   "$(extract body.invoice_amount   "$resp")" 125.37
  assert_field  "invoice_currency" "$(extract body.invoice_currency "$resp")" "USD"
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

ALL_SCENARIOS=(
  session-1-no-fractions
  session-1-some-fractions
  session-2-some-fractions-in-oak-city
  session-2.5-additional-fee-when-3-deliveries-or-more-in-same-month
  session-2.5-complex-scenario
  session-3-business-customers
  session-3-tier-based-pricing
  session-4-business-exemption-per-business
  session-5-send-invoice-after-payment
)

SELECTED=("${@:-${ALL_SCENARIOS[@]}}")

for scenario in "${SELECTED[@]}"; do
  echo ""
  echo "=== $scenario ==="
  # Map scenario name to function: replace - and . with _
  fn="scenario_$(echo "$scenario" | tr '.-' '_')"
  "$fn"
done

echo ""
echo "Results: $PASS passed, $FAIL failed"
[ "$FAIL" -eq 0 ]
