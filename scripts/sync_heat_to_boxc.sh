#!/usr/bin/env bash
# sync_heat_to_boxc.sh - Synchronizes heating database and jsonl to boxC
set -euo pipefail

SRC_DIR="$HOME/src/observator/data"
DEST_HOST="boxC"
DEST_DIR="/root/rtl_sdr-heat_wmbus_parser/data"

if [ -f "$SRC_DIR/heat_meters.db" ]; then
    ssh "$DEST_HOST" "mkdir -p $DEST_DIR"
    rsync -az "$SRC_DIR/heat_meters.db" "$SRC_DIR/heat_meters.jsonl" "$DEST_HOST:$DEST_DIR/"
fi
