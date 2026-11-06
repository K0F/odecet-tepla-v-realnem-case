#!/usr/bin/env bash
# sync_heat_to_boxc.sh - Synchronizes heating database and jsonl to boxC
set -euo pipefail

SRC_DIR="$HOME/src/observator/data"

if [ -f "$SRC_DIR/heat_meters.jsonl" ]; then
    rsync -az "$SRC_DIR/heat_meters.jsonl" "$SRC_DIR/heat_meters.db" boxC:/tmp/
    ssh boxC "sudo mv /tmp/heat_meters.jsonl /tmp/heat_meters.db /root/rtl_sdr-heat_wmbus_parser/"
fi
