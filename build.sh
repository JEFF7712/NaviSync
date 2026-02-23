#!/bin/bash
set -e

echo "Checking for Go..."
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.23+."
    exit 1
fi

echo "Checking for TinyGo..."
if ! command -v tinygo &> /dev/null; then
    echo "TinyGo is not installed. Please install TinyGo 0.34+."
    echo "See https://tinygo.org/getting-started/install/"
    exit 1
fi

echo "Downloading dependencies..."
go mod tidy

echo "Building plugin..."
make build

echo "Done! Plugin is at navisync.ndp"
