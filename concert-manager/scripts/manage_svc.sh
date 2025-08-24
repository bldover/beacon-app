#!/bin/bash

# GCP deployment script for concert manager
# Usage: ./deploy_gcp.sh [deploy|kill|status] [env_vars_file]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Service configuration
LOCAL_BUILD_DIR="${LOCAL_PROJECT_ROOT}/build"
BINARY_NAME="cm-server"
REMOTE_INSTALL_DIR="/opt/concert-manager"
REMOTE_BINARY_PATH="${REMOTE_INSTALL_DIR}/${BINARY_NAME}"
REMOTE_ENV_FILE="${REMOTE_INSTALL_DIR}/environment.conf"
REMOTE_SERVICE_NAME="concert-manager"
REMOTE_LOG_FILE="${REMOTE_BINARY_PATH}.log"
SYSTEMD_SERVICE_FILE="/etc/systemd/system/${REMOTE_SERVICE_NAME}.service"

show_usage() {
    echo "Usage: $0 [deploy|start|stop|restart|status] [environment]"
    echo ""
    echo "Commands:"
    echo "  deploy   - Build, upload, and start the service"
    echo "  start    - Start the service"
    echo "  stop     - Stop the service"
    echo "  restart  - Restart the service"
    echo "  status   - Show service status and recent logs"
    echo ""
    echo "Arguments:"
    echo "  environment - Environment to manage (optional)"
    echo "                Options: dev, prod"
    echo "                Default: dev"
    echo ""
    echo "Examples:"
    echo "  $0 deploy       # Deploy to dev environment"
    echo "  $0 start prod   # Start production service"
    echo "  $0 restart dev  # Restart dev service"
    echo "  $0 stop prod    # Stop production service"
    echo "  $0 status dev   # Check dev service status"
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
            # If it looks like a file path, use it as-is
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

    local required_vars=("PROJECT_ID" "INSTANCE_NAME" "INSTANCE_ZONE" "SSH_USER" "LOCAL_PROJECT_ROOT")
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

create_systemd_service() {
    echo "Creating systemd service file..."
    cat > "/tmp/${REMOTE_SERVICE_NAME}.service" << EOF
[Unit]
Description=Beacon - Concert Manager
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=${REMOTE_INSTALL_DIR}
EnvironmentFile=${REMOTE_ENV_FILE}
ExecStart=/bin/bash -c 'source ${REMOTE_ENV_FILE} && ${REMOTE_BINARY_PATH}'
StandardOutput=append:${REMOTE_BINARY_PATH}.log
StandardError=append:${REMOTE_BINARY_PATH}.log
Restart=no

[Install]
WantedBy=multi-user.target
EOF
}

deploy_service() {
    echo "Starting deployment to $INSTANCE_NAME..."

    echo "Building application..."
    cd "$LOCAL_PROJECT_ROOT"
    make server

    if [[ ! -f "$LOCAL_BUILD_DIR/$BINARY_NAME" ]]; then
        echo "Error: Build failed - binary not found at $LOCAL_BUILD_DIR/$BINARY_NAME"
        exit 1
    fi

    create_systemd_service
    check_gcp_connectivity

    echo "Setting up remote directories..."
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo mkdir -p ${REMOTE_INSTALL_DIR}"

    if [[ -n "$LOCAL_ENV_FILE" && -f "$LOCAL_ENV_FILE" ]]; then
        echo "Uploading environment configuration..."
        gcloud compute scp "$LOCAL_ENV_FILE" "${SSH_USER}@${INSTANCE_NAME}:/tmp/environment.conf" --zone "${INSTANCE_ZONE}"
        gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
            "sudo mv /tmp/environment.conf ${REMOTE_ENV_FILE}"
    else
        echo "Warning: No LOCAL_ENV_FILE specified or file not found"
    fi

    echo "Installing systemd service..."
    gcloud compute scp "/tmp/${REMOTE_SERVICE_NAME}.service" "${SSH_USER}@${INSTANCE_NAME}:/tmp/${REMOTE_SERVICE_NAME}.service" --zone "${INSTANCE_ZONE}"
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo mv /tmp/${REMOTE_SERVICE_NAME}.service ${SYSTEMD_SERVICE_FILE} && sudo systemctl daemon-reload"

    echo "Stopping existing service..."
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl stop ${REMOTE_SERVICE_NAME} || true"

    echo "Uploading application binary..."
    gcloud compute scp "${LOCAL_BUILD_DIR}/${BINARY_NAME}" "${SSH_USER}@${INSTANCE_NAME}:/tmp/${BINARY_NAME}" --zone "${INSTANCE_ZONE}"
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo mv /tmp/${BINARY_NAME} ${REMOTE_BINARY_PATH} && sudo chmod +x ${REMOTE_BINARY_PATH}"

    echo "Starting service..."
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl enable ${REMOTE_SERVICE_NAME} && sudo systemctl start ${REMOTE_SERVICE_NAME}"

    echo "Verifying deployment..."
    sleep 5
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl status ${REMOTE_SERVICE_NAME}"
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo tail -10 ${REMOTE_LOG_FILE}" 2>/dev/null || echo "No logs available yet"
    if gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl is-active ${REMOTE_SERVICE_NAME}" | grep -q "active"; then
        echo "Deployment completed successfully"
        echo "Service is running on $INSTANCE_NAME"
    else
        echo "Service may not have started properly"
        exit 1
    fi

    # Clean up
    rm -f "/tmp/${REMOTE_SERVICE_NAME}.service"
}

start_service() {
    echo "Starting service on $INSTANCE_NAME..."

    check_gcp_connectivity

    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl start ${REMOTE_SERVICE_NAME}"

    sleep 2

    if gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl is-active ${REMOTE_SERVICE_NAME}" | grep -q "active"; then
        echo "Service started successfully"
    else
        echo "Service may not have started properly"
        echo "Check status with: $0 status"
        exit 1
    fi
}

stop_service() {
    echo "Stopping service on $INSTANCE_NAME..."

    check_gcp_connectivity

    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl stop ${REMOTE_SERVICE_NAME}"

    echo "Service stopped"
}

restart_service() {
    echo "Restarting service on $INSTANCE_NAME..."

    check_gcp_connectivity

    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl restart ${REMOTE_SERVICE_NAME}"

    sleep 3

    if gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl is-active ${REMOTE_SERVICE_NAME}" | grep -q "active"; then
        echo "Service restarted successfully"
        
        echo ""
        echo "Recent startup logs:"
        gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
            "sudo journalctl -u ${REMOTE_SERVICE_NAME} --since '1 minute ago' --no-pager -q" || true
    else
        echo "Service may not have restarted properly"
        echo "Check status with: $0 status"
        exit 1
    fi
}

show_status() {
    echo "Checking service status on $INSTANCE_NAME..."

    check_gcp_connectivity

    echo "=== Service Status ==="
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo systemctl status ${REMOTE_SERVICE_NAME}" || echo "Service not found or inactive"

    echo ""
    echo "=== Recent Logs (last 20 lines) ==="
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "sudo tail -20 ${REMOTE_LOG_FILE}" 2>/dev/null || echo "No logs available"

    echo ""
    echo "=== Process Information ==="
    gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- \
        "pgrep -f ${BINARY_NAME} && echo 'Service process is running' || echo 'Service process not found'"
}

ENVIRONMENT="${2:-dev}"
case "$1" in
    deploy)
        load_environment "$ENVIRONMENT"
        deploy_service
        ;;
    start)
        load_environment "$ENVIRONMENT"
        start_service
        ;;
    stop)
        load_environment "$ENVIRONMENT"
        stop_service
        ;;
    restart)
        load_environment "$ENVIRONMENT"
        restart_service
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
