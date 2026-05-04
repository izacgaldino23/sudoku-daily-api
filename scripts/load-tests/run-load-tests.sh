#!/bin/sh
# run-load-tests.sh - Orchestrates Vegeta load tests for high-load endpoints

set -e

API_URL="${API_URL:-http://api-load:8080}"
REPORTS_DIR="/reports"

echo "=========================================="
echo "Starting Load Tests for Sudoku Daily API"
echo "=========================================="

# Setup authentication
echo "\n[1/6] Setting up authentication..."
/scripts/setup-auth.sh

# Clean previous reports
rm -f "$REPORTS_DIR"/*.txt "$REPORTS_DIR"/*.html

# Test 1: GET /api/sudoku (public, highest volume)
# echo "\n[2/6] Load testing GET /api/sudoku (1000 req/s for 30s)..."
# vegeta attack -targets /scripts/targets/get-sudoku.txt -rate 1000 -duration 30s -output "$REPORTS_DIR/sudoku.bin"
# vegeta report -type text < "$REPORTS_DIR/sudoku.bin" > "$REPORTS_DIR/get-sudoku.txt"
# vegeta plot < "$REPORTS_DIR/sudoku.bin" > "$REPORTS_DIR/get-sudoku.html"
# echo "  Report saved: $REPORTS_DIR/get-sudoku.txt"
# echo "  Plot saved: $REPORTS_DIR/get-sudoku.html"

# Test 2: POST /api/sudoku/submit/guest (guest submissions)
echo "\n[3/6] Load testing POST /api/sudoku/submit/guest (500 req/s for 30s)..."
vegeta attack -targets /scripts/targets/submit-guest.txt -rate 500 -duration 30s -output "$REPORTS_DIR/submit-guest.bin"
vegeta report -type text < "$REPORTS_DIR/submit-guest.bin" > "$REPORTS_DIR/submit-guest.txt"
vegeta plot < "$REPORTS_DIR/submit-guest.bin" > "$REPORTS_DIR/submit-guest.html"
echo "  Report saved: $REPORTS_DIR/submit-guest.txt"
echo "  Plot saved: $REPORTS_DIR/submit-guest.html"

# Test 3: POST /api/sudoku/submit (authenticated submissions)
echo "\n[4/6] Load testing POST /api/sudoku/submit (authenticated, 500 req/s for 30s)..."
if [ -f /scripts/targets/submit-auth.txt ]; then
    vegeta attack -targets /scripts/targets/submit-auth.txt -rate 500 -duration 30s -output "$REPORTS_DIR/submit-auth.bin"
    vegeta report -type text < "$REPORTS_DIR/submit-auth.bin" > "$REPORTS_DIR/submit-auth.txt"
    vegeta plot < "$REPORTS_DIR/submit-auth.bin" > "$REPORTS_DIR/submit-auth.html"
    echo "  Report saved: $REPORTS_DIR/submit-auth.txt"
else
    echo "  Skipped: submit-auth.txt not found (auth setup may have failed)"
fi

# Test 4: GET /api/leaderboard (public, heavy queries)
echo "\n[5/6] Load testing GET /api/leaderboard (300 req/s for 30s)..."
vegeta attack -targets /scripts/targets/leaderboard.txt -rate 300 -duration 30s -output "$REPORTS_DIR/leaderboard.bin"
vegeta report -type text < "$REPORTS_DIR/leaderboard.bin" > "$REPORTS_DIR/leaderboard.txt"
vegeta plot < "$REPORTS_DIR/leaderboard.bin" > "$REPORTS_DIR/leaderboard.html"
echo "  Report saved: $REPORTS_DIR/leaderboard.txt"

# Test 5: POST /api/auth/login (authentication endpoint)
echo "\n[6/6] Load testing POST /api/auth/login (200 req/s for 30s)..."
if [ -f /scripts/targets/login.txt ]; then
    vegeta attack -targets /scripts/targets/login.txt -rate 200 -duration 30s -output "$REPORTS_DIR/login.bin"
    vegeta report -type text < "$REPORTS_DIR/login.bin" > "$REPORTS_DIR/login.txt"
    vegeta plot < "$REPORTS_DIR/login.bin" > "$REPORTS_DIR/login.html"
    echo "  Report saved: $REPORTS_DIR/login.txt"
fi

# Generate summary
echo "\n=========================================="
echo "Load Testing Complete!"
echo "=========================================="
echo "\nSummary of Results:"
echo "------------------------------------------"
for report in "$REPORTS_DIR"/*.txt; do
    if [ -f "$report" ]; then
        echo "\n--- $(basename "$report") ---"
        head -n 10 "$report"
    fi
done

echo "\nAll reports saved to: $REPORTS_DIR"
echo "Access HTML plots in browser or copy from container."
