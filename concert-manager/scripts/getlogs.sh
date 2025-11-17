#!/bin/bash

# Script to download log files from GCP deployment
# Usage: ./getlogs.sh [environment] [lines]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BINARY_NAME="cm-server"
REMOTE_INSTALL_DIR="/opt/concert-manager"
REMOTE_BINARY_PATH="${REMOTE_INSTALL_DIR}/${BINARY_NAME}"
REMOTE_LOG_FILE="${REMOTE_BINARY_PATH}.log"
LOCAL_LOG_FILE="./${BINARY_NAME}.log"

show_usage() {
    echo "Usage: $0 [environment] [lines]"
    echo ""
    echo "Arguments:"
    echo "  environment - Environment to get logs from (optional)"
    echo "                Options: dev, prod"
    echo "                Default: dev"
    echo ""
    echo "  lines       - Number of lines from end of log to retrieve (optional)"
    echo "                If not specified, retrieves entire log file"
    echo "                Must be a positive integer"
    echo ""
    echo "Examples:"
    echo "  $0              # Get all logs from dev environment"
    echo "  $0 dev          # Get all logs from dev environment"
    echo "  $0 prod         # Get all logs from prod environment"
    echo "  $0 dev 100      # Get last 100 lines from dev environment"
    echo "  $0 prod 500     # Get last 500 lines from prod environment"
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

get_logs() {
    local num_lines="$1"

    echo "Downloading logs from $INSTANCE_NAME..."

    check_gcp_connectivity

    if [[ -n "$num_lines" ]]; then
        echo "Retrieving last $num_lines lines..."

        gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
            "sudo tail -n ${num_lines} ${REMOTE_LOG_FILE}" > "${LOCAL_LOG_FILE}" 2>&1

        if [[ $? -eq 0 ]]; then
            echo "Last $num_lines lines downloaded to: ${LOCAL_LOG_FILE}"
        else
            echo "Error: Failed to retrieve logs"
            exit 1
        fi
    else
        echo "Retrieving entire log file..."

        gcloud compute scp "${SSH_USER}@${INSTANCE_NAME}:${REMOTE_LOG_FILE}" "${LOCAL_LOG_FILE}" --zone "${INSTANCE_ZONE}"

        echo "Logs downloaded to: ${LOCAL_LOG_FILE}"
    fi
}

ENVIRONMENT="${1:-dev}"
NUM_LINES="${2:-}"

if [[ "$ENVIRONMENT" == "-h" || "$ENVIRONMENT" == "--help" ]]; then
    show_usage
    exit 0
fi

if [[ -n "$NUM_LINES" ]]; then
    if ! [[ "$NUM_LINES" =~ ^[0-9]+$ ]] || [[ "$NUM_LINES" -le 0 ]]; then
        echo "Error: lines argument must be a positive integer"
        echo ""
        show_usage
        exit 1
    fi
fi

load_environment "$ENVIRONMENT"
get_logs "$NUM_LINES"
