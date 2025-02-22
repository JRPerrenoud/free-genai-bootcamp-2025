#!/bin/bash

# Kill any running Next.js processes
echo "Killing any existing Next.js processes..."
pkill -f "node.*next dev" || true

# Function to kill process on a specific port (more aggressive)
kill_port() {
    local port=$1
    echo "Checking port $port..."
    
    # Try lsof first
    pid=$(lsof -ti:$port 2>/dev/null)
    if [ ! -z "$pid" ]; then
        echo "Found process $pid on port $port, killing it..."
        kill -9 $pid 2>/dev/null
    fi
    
    # Also try fuser as backup
    fuser -k $port/tcp 2>/dev/null
    
    # Wait a moment to ensure the port is freed
    sleep 2
    
    # Verify the port is actually free
    if lsof -i:$port >/dev/null 2>&1; then
        echo "Warning: Port $port is still in use!"
        return 1
    else
        echo "Port $port is now free"
        return 0
    fi
}

# Array of common Next.js development ports
ports=(3000 3001 3002 3003 3004 3005)

# Stop any running processes on these ports
for port in "${ports[@]}"; do
    kill_port $port
done

# Double check port 3000 specifically
if ! kill_port 3000; then
    echo "Failed to free port 3000. Please check manually with 'lsof -i:3000'"
    exit 1
fi

# Start the Next.js development server
echo "Starting Next.js development server..."
npm run dev
