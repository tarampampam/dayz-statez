#!/usr/bin/env bash

SERVER_DIR="$HOME/.local/share/Steam/steamapps/common/DayZ Server Exp"
MOD_DIR="$SERVER_DIR/@metricz/metricz"
PROFILES="$SERVER_DIR/test_profiles"
SCRIPT_ROOT="$(cd "$(dirname "$0")" && pwd)"

# Keep symlinks fresh
mkdir -p "$MOD_DIR"
ln -sfn "$SCRIPT_ROOT/config.cpp" "$MOD_DIR/config.cpp"
ln -sfn "$SCRIPT_ROOT/scripts" "$MOD_DIR/scripts"

# Clean old logs
rm -f "$PROFILES"/DayZServer_*.RPT "$PROFILES"/script_*.log "$PROFILES"/error_*.log

# DayZ server must run from its own directory (relative config/profiles paths)
cd "$SERVER_DIR"

"$SERVER_DIR/DayZServer" \
    -config=test_serverDZ.cfg \
    -mod=@metricz \
    -filePatching \
    -profiles=test_profiles \
    -port=12302 \
    -nologs \
    >/dev/null 2>&1 &
SRV_PID=$!

echo "Server PID $SRV_PID — waiting for script compilation..."

SCRIPT_LOG=""
for i in $(seq 1 60); do
    sleep 2
    SCRIPT_LOG=$(ls "$PROFILES"/script_*.log 2>/dev/null | tail -1 || true)
    if [[ -n "$SCRIPT_LOG" ]]; then
        MODULES=$(grep -c "SCRIPT.*Module:.*loaded" "$SCRIPT_LOG" 2>/dev/null || true)
        if [[ "$MODULES" -ge 4 ]]; then
            break
        fi
    fi
done

kill "$SRV_PID" 2>/dev/null || true
wait "$SRV_PID" 2>/dev/null || true

if [[ -z "$SCRIPT_LOG" ]]; then
    echo "ERROR: script log not found — server may have crashed"
    exit 1
fi

echo ""
echo "=== Script compilation results ==="
grep "Module:.*loaded" "$SCRIPT_LOG" || true

# Script errors (E) — these are real compile failures
ERRORS=$(grep "SCRIPT.*(E)" "$SCRIPT_LOG" 2>/dev/null || true)
METRICZ_ISSUES=$(grep -i "metricz" "$SCRIPT_LOG" 2>/dev/null | grep -iv "Module:.*loaded" || true)

if [[ -n "$ERRORS" ]]; then
    echo ""
    echo "SCRIPT ERRORS:"
    echo "$ERRORS"
    exit 1
fi

if [[ -n "$METRICZ_ISSUES" ]]; then
    echo ""
    echo "MetricZ log entries:"
    echo "$METRICZ_ISSUES"
fi

echo ""
echo "OK — no script compilation errors"
