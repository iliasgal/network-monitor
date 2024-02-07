#!/bin/bash

# Prompt for sudo password upfront
echo "Please enter your sudo password to continue..."
sudo -v

# Keep-alive: update existing sudo time stamp if set, otherwise do nothing.
while true; do sudo -n true; sleep 60; kill -0 "$$" || exit; done 2>/dev/null &


# Load environment variables from .env file
if [ -f ".env" ]; then
    echo "Loading environment variables from .env file..."
    set -a  # Automatically export all variables
    source .env
    set +a
else 
    echo "No .env file found"
fi

# Define the temporary directory for the build
BUILD_DIR=$(mktemp -d)

# Ensure the temporary directory is removed on script exit
trap "rm -rf $BUILD_DIR" EXIT

# Navigate to your Go project directory
cd cmd

# Build the Go program in the temporary directory
echo "Building Go program in $BUILD_DIR..."
go build -o "$BUILD_DIR/main"
echo "Go program built successfully!"

# Ensure InfluxDB container is running
echo "Starting InfluxDB container..."
docker-compose up -d influxdb

# Wait a bit for InfluxDB to initialize (optional)
sleep 5

# Run the Go program from the temporary directory
echo "Running the network monitor program..."
sudo -E "$BUILD_DIR/main"
