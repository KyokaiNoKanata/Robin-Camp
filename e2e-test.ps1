# PowerShell E2E Test Script for Movies API

# Configuration
$BASE_URL = "http://localhost:9090"
$AUTH_TOKEN = "Bearer test-token"
$TESTS_PASSED = 0
$TESTS_FAILED = 0

# Utility functions
function Write-ColorOutput {
    param(
        [string]$Message,
        [string]$Color
    )
    Write-Host $Message -ForegroundColor $Color
}

function Write-Success {
    param([string]$Message)
    Write-ColorOutput "[SUCCESS] $Message" "Green"
    $script:TESTS_PASSED++
}

function Write-ErrorOutput {
    param([string]$Message)
    Write-ColorOutput "[ERROR] $Message" "Red"
    $script:TESTS_FAILED++
}

function Write-Info {
    param([string]$Message)
    Write-ColorOutput "[INFO] $Message" "Blue"
}

function Test-Endpoint {
    param(
        [string]$Method,
        [string]$Endpoint,
        [hashtable]$Headers = @{},
        [string]$Body = $null,
        [int]$ExpectedStatus = 200
    )
    
    $url = "$BASE_URL$Endpoint"
    Write-Info "Making $Method request to $url"
    
    try {
        $params = @{
            Uri = $url
            Method = $Method
            UseBasicParsing = $true
            TimeoutSec = 10
        }
        
        if ($Headers.Count -gt 0) {
            $params.Headers = $Headers
        }
        
        if ($Body) {
            $params.Body = $Body
        }
        
        $response = Invoke-WebRequest @params
        $status = $response.StatusCode
        
        if ($status -eq $ExpectedStatus) {
            Write-Success "$Endpoint - Status $status (expected $ExpectedStatus)"
            return $response
        } else {
            Write-ErrorOutput "$Endpoint - Expected status $ExpectedStatus, got $status"
            Write-ErrorOutput "Response: $($response.Content)"
            return $null
        }
    } catch {
        $status = $_.Exception.Response.StatusCode.value__
        Write-ErrorOutput "$Endpoint - Error: $status"
        if ($_.Exception.Response) {
            $content = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream()).ReadToEnd()
            Write-ErrorOutput "Response: $content"
        }
        return $null
    }
}

# Test Suite
Write-Info "Starting E2E Tests for Movie Rating API"

# Test 1: Health Check
Write-Info "=== Test 1: Health Check ==="
$healthResponse = Test-Endpoint -Method Get -Endpoint "/healthz" -ExpectedStatus 200
if ($healthResponse) {
    Write-Success "Health check passed"
}

# Test 2: Create a Movie
Write-Info "=== Test 2: Create Movie ==="
$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$uniqueMovieTitle = "E2E Test Movie $timestamp"
$movieBody = @{
    title = $uniqueMovieTitle
    releaseDate = "2024-01-01"
    genre = "Action"
} | ConvertTo-Json

$movieHeaders = @{
    "Content-Type" = "application/json"
    "Authorization" = $AUTH_TOKEN
}

$createResponse = Test-Endpoint -Method Post -Endpoint "/movies" -Headers $movieHeaders -Body $movieBody -ExpectedStatus 201

# Test 3: List Movies
Write-Info "=== Test 3: List Movies ==="
$listHeaders = @{
    "Authorization" = $AUTH_TOKEN
}

$listResponse = Test-Endpoint -Method Get -Endpoint "/movies" -Headers $listHeaders -ExpectedStatus 200

# Test 4: Submit Rating
Write-Info "=== Test 4: Submit Rating ==="
if ($createResponse) {
    $ratingBody = @{
        score = 5
        comment = "Great movie!"
    } | ConvertTo-Json
    
    $ratingHeaders = @{
        "Content-Type" = "application/json"
        "Authorization" = $AUTH_TOKEN
        "X-Rater-Id" = "test-user-123"
    }
    
    # URL encode the unique movie title
    $encodedMovieTitle = [System.Web.HttpUtility]::UrlEncode($uniqueMovieTitle)
    Test-Endpoint -Method Post -Endpoint "/movies/$encodedMovieTitle/ratings" -Headers $ratingHeaders -Body $ratingBody -ExpectedStatus 201
}

# Test 5: Get Movie Ratings
Write-Info "=== Test 5: Get Movie Ratings ==="
Test-Endpoint -Method Get -Endpoint "/movies/$encodedMovieTitle/ratings" -Headers $listHeaders -ExpectedStatus 200

# Summary
Write-Host "`n=== TEST SUMMARY ===" -ForegroundColor Yellow
Write-Host "Tests Passed: $TESTS_PASSED" -ForegroundColor Green
Write-Host "Tests Failed: $TESTS_FAILED" -ForegroundColor Red

if ($TESTS_FAILED -eq 0) {
    Write-Host "All tests passed!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "Some tests failed!" -ForegroundColor Red
    exit 1
}