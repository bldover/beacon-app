#!/bin/bash

PROJECT_ID="concert-manager-dev"
INSTANCE_NAME="concert-manager-dev"
INSTANCE_ZONE="us-east1-b"
SSH_USER="dover_bradley"

LOCAL_PROJECT_ROOT="/mnt/c/Users/User/workspace/beacon-app/concert-manager"
LOCAL_BUILD_DIR="${LOCAL_PROJECT_ROOT}/build"
BINARY_NAME="cm-server"
LOCAL_ENV_FILE="${LOCAL_PROJECT_ROOT}/resources/dev.env"

REMOTE_INSTALL_DIR="/opt/concert-manager"
REMOTE_BINARY_PATH="${REMOTE_INSTALL_DIR}/${BINARY_NAME}"
REMOTE_ENV_FILE="${REMOTE_INSTALL_DIR}/environment.conf"
REMOTE_SERVICE_NAME="concert-manager"
REMOTE_LOG_FILE="${REMOTE_BINARY_PATH}.log"

SYSTEMD_SERVICE_FILE="/etc/systemd/system/${REMOTE_SERVICE_NAME}.service"

set -e
set -x

# Build the binary
cd "${LOCAL_PROJECT_ROOT}"
make server

# Create systemd service file
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

gcloud config set project "${PROJECT_ID}"

# Create remote directories
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mkdir -p ${REMOTE_INSTALL_DIR}"

# Upload environment variables
gcloud compute scp "${LOCAL_ENV_FILE}" "${SSH_USER}@${INSTANCE_NAME}:/tmp/environment.conf" --zone "${INSTANCE_ZONE}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mv /tmp/environment.conf ${REMOTE_ENV_FILE}"

# Upload systemd service file
gcloud compute scp "/tmp/${REMOTE_SERVICE_NAME}.service" "${SSH_USER}@${INSTANCE_NAME}:/tmp/${REMOTE_SERVICE_NAME}.service" --zone "${INSTANCE_ZONE}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mv /tmp/${REMOTE_SERVICE_NAME}.service ${SYSTEMD_SERVICE_FILE} && sudo systemctl daemon-reload"

# Kill existing instance
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo systemctl stop ${REMOTE_SERVICE_NAME} || true"

# Upload new binary
gcloud compute scp "${LOCAL_BUILD_DIR}/${BINARY_NAME}" "${SSH_USER}@${INSTANCE_NAME}:/tmp/${BINARY_NAME}" --zone "${INSTANCE_ZONE}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mv /tmp/${BINARY_NAME} ${REMOTE_BINARY_PATH} && sudo chmod +x ${REMOTE_BINARY_PATH}"

# Start new binary
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo systemctl start ${REMOTE_SERVICE_NAME}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo systemctl status ${REMOTE_SERVICE_NAME}"

# Verify the service
sleep 5
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo tail -10 ${REMOTE_LOG_FILE}"
if gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo grep -q 'Starting server on port :3001' ${REMOTE_LOG_FILE}"; then
    echo "${REMOTE_SERVICE_NAME} deployment completed successfully"
else
    echo "Server may not have started properly - check logs"
    exit 1
fi
