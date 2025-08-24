#!/bin/bash

# Script to enable/disable cron jobs for concert manager cache refreshing
# Usage: ./manage_cron.sh [enable|disable|status] [server_url]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CRON_COMMENT="# Concert Manager Cache Refresh Jobs"

# Default server URL (can be overridden)
SERVER_URL="${2:-http://localhost:3001}"

show_usage() {
    echo "Usage: $0 [enable|disable|status] [server_url]"
    echo ""
    echo "Commands:"
    echo "  enable   - Enable cron jobs for cache refreshing"
    echo "  disable  - Disable cron jobs for cache refreshing"
    echo "  status   - Show current cron job status"
    echo ""
    echo "Arguments:"
    echo "  server_url - Base URL of the concert manager server (default: http://localhost:3001)"
    echo ""
    echo "Examples:"
    echo "  $0 enable"
    echo "  $0 enable http://your-server.com:3001"
    echo "  $0 disable"
    echo "  $0 status"
}

check_server_health() {
    local url="$1"
    echo "Checking server health at $url..."
    
    if ! curl -s -f "$url/v1/artists" > /dev/null 2>&1; then
        echo "Warning: Server at $url appears to be unreachable"
        echo "The cron jobs will still be configured, but may fail until the server is running"
        return 1
    else
        echo "Server health check passed"
        return 0
    fi
}

enable_cron() {
    local server_url="$1"
    
    echo "Enabling cron jobs for cache refresh..."
    echo "Server URL: $server_url"
    
    # Check if server is reachable (but don't fail if not)
    check_server_health "$server_url" || true
    
    # Create temporary cron file with existing crontab
    local temp_cron=$(mktemp)
    crontab -l 2>/dev/null | grep -v "Concert Manager Cache Refresh" > "$temp_cron" || true
    
    # Add our cron jobs
    echo "" >> "$temp_cron"
    echo "$CRON_COMMENT" >> "$temp_cron"
    echo "# Refresh artist ranks every Friday at 5:00 AM EST" >> "$temp_cron"
    echo "0 10 * * 5 curl -X POST -s $server_url/v1/ranks/refresh > /dev/null 2>&1" >> "$temp_cron"
    echo "# Refresh upcoming events daily at 6:00 AM EST" >> "$temp_cron"
    echo "0 11 * * * curl -X POST -s $server_url/v1/events/upcoming/refresh > /dev/null 2>&1" >> "$temp_cron"
    
    # Install the new crontab
    crontab "$temp_cron"
    rm "$temp_cron"
    
    echo "Cron jobs enabled successfully!"
    echo ""
    echo "Scheduled jobs:"
    echo "  - Artist ranks refresh: Fridays at 5:00 AM EST (10:00 AM UTC)"
    echo "  - Events refresh: Daily at 6:00 AM EST (11:00 AM UTC)"
    echo ""
    echo "Note: Times shown are in UTC (server time). Adjust if your server timezone differs."
}

disable_cron() {
    echo "Disabling cron jobs for cache refresh..."
    
    # Create temporary cron file without our jobs
    local temp_cron=$(mktemp)
    crontab -l 2>/dev/null | grep -v "Concert Manager Cache Refresh" | grep -v "/v1/ranks/refresh" | grep -v "/v1/events/upcoming/refresh" > "$temp_cron" || true
    
    # Remove any empty lines at the end and our comment blocks
    sed -i '/^# Refresh artist ranks/d; /^# Refresh upcoming events/d' "$temp_cron" 2>/dev/null || true
    
    # Install the cleaned crontab
    crontab "$temp_cron"
    rm "$temp_cron"
    
    echo "Cron jobs disabled successfully!"
}

show_status() {
    echo "Current cron job status:"
    echo ""
    
    local cron_output=$(crontab -l 2>/dev/null || echo "No crontab found")
    
    if echo "$cron_output" | grep -q "Concert Manager Cache Refresh"; then
        echo "✓ Concert Manager cron jobs are ENABLED"
        echo ""
        echo "Active jobs:"
        echo "$cron_output" | grep -A 10 "Concert Manager Cache Refresh" | grep -E "(ranks|events)/refresh"
    else
        echo "✗ Concert Manager cron jobs are DISABLED"
    fi
    
    echo ""
    echo "Full crontab:"
    echo "$cron_output"
}

# Main script logic
case "$1" in
    enable)
        enable_cron "$SERVER_URL"
        ;;
    disable)
        disable_cron
        ;;
    status)
        show_status
        ;;
    *)
        show_usage
        exit 1
        ;;
esac