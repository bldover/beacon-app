#!/bin/bash

# Script to enable/disable cron jobs for concert manager cache refreshing
# Usage: ./manage_cron.sh [enable|disable|status] [environment]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CRON_COMMENT="# Concert Manager Cache Refresh Jobs"
SERVER_URL="http://localhost:3001"

show_usage() {
    echo "Usage: $0 [enable|disable|status] [environment]"
    echo ""
    echo "Commands:"
    echo "  enable   - Enable cron jobs for cache refreshing on remote server"
    echo "  disable  - Disable cron jobs for cache refreshing on remote server"
    echo "  status   - Show current cron job status on remote server"
    echo ""
    echo "Arguments:"
    echo "  environment - Environment to manage (optional)"
    echo "                Options: dev, prod"
    echo "                Default: dev"
    echo ""
    echo "Examples:"
    echo "  $0 enable        # Enable for dev environment"
    echo "  $0 enable dev    # Enable for dev environment"
    echo "  $0 enable prod   # Enable for prod environment"
    echo "  $0 disable dev   # Disable dev cron jobs"
    echo "  $0 status prod   # Check prod status"
}

resolve_environment_file() {
    local env_name="$1"

    case "$env_name" in
        dev)
            echo "$SCRIPT_DIR/env/vars_dev.sh"
            ;;
        prod)
            echo "$SCRIPT_DIR/env/vars_prod.sh"
            ;;
        *)
            if [[ "$env_name" == *"/"* || "$env_name" == *".sh" ]]; then
                echo "$env_name"
            else
                echo "Error: Invalid environment '$env_name'. Valid options: dev, prod"
                exit 1
            fi
            ;;
    esac
}

load_environment() {
    local env_name="$1"
    local env_file

    env_file=$(resolve_environment_file "$env_name")

    if [[ ! -f "$env_file" ]]; then
        echo "Error: Environment file not found: $env_file"
        echo "Available environments: dev, prod"
        echo "Corresponding files:"
        echo "  dev  -> $SCRIPT_DIR/env/vars_dev.sh"
        echo "  prod -> $SCRIPT_DIR/env/vars_prod.sh"
        exit 1
    fi

    echo "Loading environment: $env_name"
    echo "From file: $env_file"
    source "$env_file"

    local required_vars=("PROJECT_ID" "INSTANCE_NAME" "INSTANCE_ZONE" "SSH_USER")
    for var in "${required_vars[@]}"; do
        if [[ -z "${!var}" ]]; then
            echo "Error: Required environment variable $var is not set in $env_file"
            exit 1
        fi
    done

    echo "Environment loaded:"
    echo "  Project: $PROJECT_ID"
    echo "  Instance: $INSTANCE_NAME"
    echo "  Zone: $INSTANCE_ZONE"
    echo "  User: $SSH_USER"
}

check_gcp_connectivity() {
    echo "Checking GCP connectivity..."

    if ! gcloud config set project "$PROJECT_ID" 2>/dev/null; then
        echo "Error: Failed to set GCP project $PROJECT_ID"
        echo "Please check your GCP authentication and project access"
        exit 1
    fi

    if ! gcloud compute instances describe "$INSTANCE_NAME" --zone "$INSTANCE_ZONE" &>/dev/null; then
        echo "Error: Cannot access instance $INSTANCE_NAME in zone $INSTANCE_ZONE"
        echo "Please check instance name and zone, and ensure the instance exists"
        exit 1
    fi

    echo "GCP connectivity verified"
}

enable_cron() {
    echo "Enabling cron jobs on remote server..."
    echo "Server: $INSTANCE_NAME"
    echo "Server URL: $SERVER_URL"

    check_gcp_connectivity

    local temp_cron=$(mktemp)

    cat > "$temp_cron" << 'EOF'
# Concert Manager Cache Refresh Jobs
# Refresh artist ranks every Friday at 5:00 AM EST
0 10 * * 5 curl -X POST -s http://localhost:3001/v1/ranks/refresh > /dev/null 2>&1
# Refresh upcoming events daily at 6:00 AM EST
0 11 * * * curl -X POST -s http://localhost:3001/v1/events/upcoming/refresh > /dev/null 2>&1
EOF

    echo "Uploading cron configuration to server..."
    gcloud compute scp "$temp_cron" "${SSH_USER}@${INSTANCE_NAME}:/tmp/cm-cron.txt" --zone "${INSTANCE_ZONE}"

    echo "Installing cron jobs on server..."
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        'crontab -l 2>/dev/null | grep -v "Concert Manager Cache Refresh" | grep -v "/v1/ranks/refresh" | grep -v "/v1/events/upcoming/refresh" > /tmp/existing-cron.txt || true; cat /tmp/existing-cron.txt /tmp/cm-cron.txt | crontab -; rm /tmp/cm-cron.txt /tmp/existing-cron.txt'

    rm "$temp_cron"

    echo "Cron jobs enabled successfully!"
    echo ""
    echo "Scheduled jobs:"
    echo "  - Artist ranks refresh: Fridays at 5:00 AM EST (10:00 AM UTC)"
    echo "  - Events refresh: Daily at 6:00 AM EST (11:00 AM UTC)"
}

disable_cron() {
    echo "Disabling cron jobs on remote server..."
    echo "Server: $INSTANCE_NAME"

    check_gcp_connectivity

    echo "Removing cron jobs from server..."
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        'crontab -l 2>/dev/null | grep -v "Concert Manager Cache Refresh" | grep -v "/v1/ranks/refresh" | grep -v "/v1/events/upcoming/refresh" | sed "/^# Refresh artist ranks/d; /^# Refresh upcoming events/d" | crontab - 2>/dev/null || crontab -r 2>/dev/null || true'

    echo "Cron jobs disabled successfully!"
}

show_status() {
    echo "Checking cron job status on remote server..."
    echo "Server: $INSTANCE_NAME"

    check_gcp_connectivity

    echo ""
    echo "=== Cron Job Status ==="

    local cron_output
    cron_output=$(gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        'crontab -l 2>/dev/null || echo "No crontab found"')

    if echo "$cron_output" | grep -q "Concert Manager Cache Refresh"; then
        echo "✓ Concert Manager cron jobs are ENABLED"
        echo ""
        echo "Active jobs:"
        echo "$cron_output" | grep -E "(ranks|events)/refresh"
    else
        echo "✗ Concert Manager cron jobs are DISABLED"
    fi

    echo ""
    echo "Full crontab:"
    echo "$cron_output"
}

ENVIRONMENT="${2:-dev}"

if [[ "$1" == "-h" || "$1" == "--help" ]]; then
    show_usage
    exit 0
fi

case "$1" in
    enable)
        load_environment "$ENVIRONMENT"
        enable_cron
        ;;
    disable)
        load_environment "$ENVIRONMENT"
        disable_cron
        ;;
    status)
        load_environment "$ENVIRONMENT"
        show_status
        ;;
    *)
        show_usage
        exit 1
        ;;
esac