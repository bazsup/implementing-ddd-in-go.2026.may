# Run DDD implementation scenarios against the API.
# Requirements: PowerShell (built into Windows 10+) — no installs needed.
#
# Usage:
#   .\postman\run.ps1                                             # run all scenarios
#   .\postman\run.ps1 -Scenarios session-1-no-fractions          # run one scenario
#   .\postman\run.ps1 -Scenarios s1,s2 -BaseUrl http://staging:8080
param(
  [string[]]$Scenarios,
  [string]$BaseUrl    = if ($env:BASE_URL)    { $env:BASE_URL }    else { "http://localhost:8080" },
  [string]$InvoiceUrl = if ($env:INVOICE_URL) { $env:INVOICE_URL } else { "http://localhost:9000" }
)

$script:Pass = 0
$script:Fail = 0

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

function Invoke-Post {
  param([string]$Url, [object]$Body = @{})
  try {
    $json = if ($Body -is [string]) { $Body } else { $Body | ConvertTo-Json -Depth 10 }
    Invoke-RestMethod -Method Post -Uri $Url `
      -ContentType "application/json" `
      -Body $json `
      -ErrorAction Stop
  } catch {
    Write-Error "POST $Url failed: $_"
    exit 1
  }
}

function Invoke-Get {
  param([string]$Url)
  try {
    Invoke-RestMethod -Method Get -Uri $Url -ErrorAction Stop
  } catch {
    Write-Error "GET $Url failed: $_"
    exit 1
  }
}

function Assert-Field {
  param([string]$Label, [string]$Actual, [string]$Expected)
  if ($Actual -eq $Expected) {
    Write-Host "  [PASS] $Label" -ForegroundColor Green
    $script:Pass++
  } else {
    Write-Host "  [FAIL] $Label  expected='$Expected'  actual='$Actual'" -ForegroundColor Red
    $script:Fail++
  }
}

function Assert-Number {
  param([string]$Label, $Actual, $Expected)
  if ([double]$Actual -eq [double]$Expected) {
    Write-Host "  [PASS] $Label" -ForegroundColor Green
    $script:Pass++
  } else {
    Write-Host "  [FAIL] $Label  expected=$Expected  actual=$Actual" -ForegroundColor Red
    $script:Fail++
  }
}

# ---------------------------------------------------------------------------
# Scenarios
# ---------------------------------------------------------------------------

function Invoke-Scenario-session-1-no-fractions {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @()
    person_id = "Bald Eagle"
    visit_id = "1"
  }
  Assert-Number "price_amount"   $r.price_amount   0
  Assert-Field  "person_id"      $r.person_id      "Bald Eagle"
  Assert-Field  "visit_id"       $r.visit_id       "1"
  Assert-Field  "price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-1-some-fractions {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(
      @{ amount_dropped = 15; fraction_type = "Green waste" }
      @{ amount_dropped = 39; fraction_type = "Construction waste" }
    )
    person_id = "Bald Eagle"
    visit_id = "1"
  }
  Assert-Number "price_amount"   $r.price_amount   7.35
  Assert-Field  "person_id"      $r.person_id      "Bald Eagle"
  Assert-Field  "visit_id"       $r.visit_id       "1"
  Assert-Field  "price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-2-some-fractions-in-oak-city {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(
      @{ amount_dropped = 83; fraction_type = "Green waste" }
      @{ amount_dropped = 18; fraction_type = "Construction waste" }
    )
    person_id = "Squirrel Gus"
    visit_id = "1"
  }
  Assert-Number "price_amount"   $r.price_amount   10.06
  Assert-Field  "person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit_id"       $r.visit_id       "1"
  Assert-Field  "price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-2.5-additional-fee-when-3-deliveries-or-more-in-same-month {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(
      @{ amount_dropped = 83; fraction_type = "Green waste" }
      @{ amount_dropped = 18; fraction_type = "Construction waste" }
    )
    person_id = "Squirrel Gus"; visit_id = "1"
  }
  Assert-Number "visit1 price_amount"   $r.price_amount   10.06
  Assert-Field  "visit1 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit1 visit_id"       $r.visit_id       "1"
  Assert-Field  "visit1 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-24"
    dropped_fractions = @(@{ amount_dropped = 49; fraction_type = "Construction waste" })
    person_id = "Squirrel Gus"; visit_id = "2"
  }
  Assert-Number "visit2 price_amount"   $r.price_amount   9.31
  Assert-Field  "visit2 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit2 visit_id"       $r.visit_id       "2"
  Assert-Field  "visit2 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-25"
    dropped_fractions = @(@{ amount_dropped = 103; fraction_type = "Green waste" })
    person_id = "Squirrel Gus"; visit_id = "3"
  }
  Assert-Number "visit3 price_amount"   $r.price_amount   8.65
  Assert-Field  "visit3 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit3 visit_id"       $r.visit_id       "3"
  Assert-Field  "visit3 price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-2.5-complex-scenario {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(
      @{ amount_dropped = 83; fraction_type = "Green waste" }
      @{ amount_dropped = 18; fraction_type = "Construction waste" }
    )
    person_id = "Squirrel Gus"; visit_id = "1"
  }
  Assert-Number "visit1 price_amount"   $r.price_amount   10.06
  Assert-Field  "visit1 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit1 visit_id"       $r.visit_id       "1"
  Assert-Field  "visit1 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-24"
    dropped_fractions = @(@{ amount_dropped = 49; fraction_type = "Construction waste" })
    person_id = "Squirrel Gus"; visit_id = "2"
  }
  Assert-Number "visit2 price_amount"   $r.price_amount   9.31
  Assert-Field  "visit2 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit2 visit_id"       $r.visit_id       "2"
  Assert-Field  "visit2 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-24"
    dropped_fractions = @(@{ amount_dropped = 10; fraction_type = "Construction waste" })
    person_id = "Bald Eagle"; visit_id = "3"
  }
  Assert-Number "visit3 price_amount"   $r.price_amount   1.5
  Assert-Field  "visit3 person_id"      $r.person_id      "Bald Eagle"
  Assert-Field  "visit3 visit_id"       $r.visit_id       "3"
  Assert-Field  "visit3 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-25"
    dropped_fractions = @(@{ amount_dropped = 103; fraction_type = "Green waste" })
    person_id = "Squirrel Gus"; visit_id = "4"
  }
  Assert-Number "visit4 price_amount"   $r.price_amount   8.65
  Assert-Field  "visit4 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit4 visit_id"       $r.visit_id       "4"
  Assert-Field  "visit4 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-08-12"
    dropped_fractions = @(@{ amount_dropped = 33; fraction_type = "Green waste" })
    person_id = "Squirrel Gus"; visit_id = "5"
  }
  Assert-Number "visit5 price_amount"   $r.price_amount   2.64
  Assert-Field  "visit5 person_id"      $r.person_id      "Squirrel Gus"
  Assert-Field  "visit5 visit_id"       $r.visit_id       "5"
  Assert-Field  "visit5 price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-3-business-customers {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(
      @{ amount_dropped = 83; fraction_type = "Green waste" }
      @{ amount_dropped = 18; fraction_type = "Construction waste" }
    )
    person_id = "Beaver Bertha"; visit_id = "1"
  }
  Assert-Number "visit1 price_amount"   $r.price_amount   10.42
  Assert-Field  "visit1 person_id"      $r.person_id      "Beaver Bertha"
  Assert-Field  "visit1 visit_id"       $r.visit_id       "1"
  Assert-Field  "visit1 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-09-30"
    dropped_fractions = @(
      @{ amount_dropped = 134; fraction_type = "Green waste" }
      @{ amount_dropped = 201; fraction_type = "Construction waste" }
    )
    person_id = "Bear Billy"; visit_id = "2"
  }
  Assert-Number "visit2 price_amount"   $r.price_amount   42.21
  Assert-Field  "visit2 person_id"      $r.person_id      "Bear Billy"
  Assert-Field  "visit2 visit_id"       $r.visit_id       "2"
  Assert-Field  "visit2 price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-3-tier-based-pricing {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(@{ amount_dropped = 597; fraction_type = "Construction waste" })
    person_id = "Beaver Bertha"; visit_id = "1"
  }
  Assert-Number "visit1 price_amount"   $r.price_amount   125.37
  Assert-Field  "visit1 person_id"      $r.person_id      "Beaver Bertha"
  Assert-Field  "visit1 visit_id"       $r.visit_id       "1"
  Assert-Field  "visit1 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(@{ amount_dropped = 1803; fraction_type = "Construction waste" })
    person_id = "Beaver Bertha"; visit_id = "2"
  }
  Assert-Number "visit2 price_amount"   $r.price_amount   490.63
  Assert-Field  "visit2 person_id"      $r.person_id      "Beaver Bertha"
  Assert-Field  "visit2 visit_id"       $r.visit_id       "2"
  Assert-Field  "visit2 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(@{ amount_dropped = 1228; fraction_type = "Construction waste" })
    person_id = "Peppa Python"; visit_id = "3"
  }
  Assert-Number "visit3 price_amount"   $r.price_amount   276.12
  Assert-Field  "visit3 person_id"      $r.person_id      "Peppa Python"
  Assert-Field  "visit3 visit_id"       $r.visit_id       "3"
  Assert-Field  "visit3 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-25"
    dropped_fractions = @(@{ amount_dropped = 901; fraction_type = "Construction waste" })
    person_id = "Beaver Bertha"; visit_id = "4"
  }
  Assert-Number "visit4 price_amount"   $r.price_amount   261.29
  Assert-Field  "visit4 person_id"      $r.person_id      "Beaver Bertha"
  Assert-Field  "visit4 visit_id"       $r.visit_id       "4"
  Assert-Field  "visit4 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2024-01-15"
    dropped_fractions = @(@{ amount_dropped = 600; fraction_type = "Construction waste" })
    person_id = "Beaver Bertha"; visit_id = "5"
  }
  Assert-Number "visit5 price_amount"   $r.price_amount   126
  Assert-Field  "visit5 person_id"      $r.person_id      "Beaver Bertha"
  Assert-Field  "visit5 visit_id"       $r.visit_id       "5"
  Assert-Field  "visit5 price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-4-business-exemption-per-business {
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(@{ amount_dropped = 597; fraction_type = "Construction waste" })
    person_id = "Beaver Bertha"; visit_id = "1"
  }
  Assert-Number "visit1 price_amount"   $r.price_amount   125.37
  Assert-Field  "visit1 person_id"      $r.person_id      "Beaver Bertha"
  Assert-Field  "visit1 visit_id"       $r.visit_id       "1"
  Assert-Field  "visit1 price_currency" $r.price_currency "USD"

  $r = Invoke-Post "$BaseUrl/calculatePrice" @{
    date = "2023-07-23"
    dropped_fractions = @(@{ amount_dropped = 1803; fraction_type = "Construction waste" })
    person_id = "Beaver Bruce"; visit_id = "2"
  }
  Assert-Number "visit2 price_amount"   $r.price_amount   490.63
  Assert-Field  "visit2 person_id"      $r.person_id      "Beaver Bruce"
  Assert-Field  "visit2 visit_id"       $r.visit_id       "2"
  Assert-Field  "visit2 price_currency" $r.price_currency "USD"
}

function Invoke-Scenario-session-5-send-invoice-after-payment {
  Invoke-Post "$InvoiceUrl/api/invoice-inspector/reset" | Out-Null
  Invoke-Post "$BaseUrl/startScenario" | Out-Null

  Invoke-Post "$InvoiceUrl/api/invoice" @{
    email = "beavers@dam-building.com"
    invoice_amount = 125.37
    invoice_currency = "USD"
  } | Out-Null

  $r = Invoke-Get "$InvoiceUrl/api/invoice-inspector/beavers@dam-building.com"
  Assert-Number "invoice_amount"   $r.body.invoice_amount   125.37
  Assert-Field  "invoice_currency" $r.body.invoice_currency "USD"
}

# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

$allScenarios = @(
  "session-1-no-fractions"
  "session-1-some-fractions"
  "session-2-some-fractions-in-oak-city"
  "session-2.5-additional-fee-when-3-deliveries-or-more-in-same-month"
  "session-2.5-complex-scenario"
  "session-3-business-customers"
  "session-3-tier-based-pricing"
  "session-4-business-exemption-per-business"
  "session-5-send-invoice-after-payment"
)

$selected = if ($Scenarios) { $Scenarios } else { $allScenarios }

foreach ($name in $selected) {
  Write-Host ""
  Write-Host "=== $name ===" -ForegroundColor Cyan
  $fn = "Invoke-Scenario-$name"
  & $fn
}

Write-Host ""
Write-Host "Results: $($script:Pass) passed, $($script:Fail) failed"
if ($script:Fail -gt 0) { exit 1 }
