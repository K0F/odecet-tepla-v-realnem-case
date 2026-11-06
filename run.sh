#!/usr/bin/env bash
# run.sh - Runs the Go heating dashboard locally
go run main.go -p 8082 -db "$HOME/src/observator/data/heat_meters.db" -jsonl "$HOME/src/observator/data/heat_meters.jsonl"
