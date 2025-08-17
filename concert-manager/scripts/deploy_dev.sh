#!/bin/bash

PROJECT_ID="concert-manager-dev"
INSTANCE_NAME="concert-manager-dev"
INSTANCE_ZONE="us-east1-b"
LOCAL_USER=bldover
SSH_USER="dover_bradley"

LOCAL_PROJECT_ROOT="/mnt/c/Users/User/workspace/beacon-app/concert-manager"
LOCAL_BUILD_DIR="${LOCAL_PROJECT_ROOT/build}"
BINARY_NAME="cm-server"
LOCAL_ENV_FILE="${LOCAL_PROJECT_ROOT}/resources/dev.env"

SYSTEMD_SERVICE_FILE="/etc/systemd/system/${REMOTE_SERVICE_NAME}.service"

REMOTE_INSTALL_DIR="/opt/concert-manager"
REMOTE_BINARY_PATH="${REMOTE_INSTALL_DIR}/${BINARY_NAME}"
REMOTE_ENV_FILE="${REMOTE_INSTALL_DIR}/environment.conf"
REMOTE_SERVICE_NAME="concert-manager"

set -e
set -x

# Build the binary
make server

# Create systemd service file
cat > "/tmp/${REMOTE_SERVICE_NAME}.service" << EOF
[Unit]
Description=Beacon - Concert Manager and Finder
After=network.target

[Service]
Type=simple
User=${SSH_USER}
EnvironmentFile=${REMOTE_ENV_FILE}
ExecStart=${REMOTE_BINARY_PATH}
Restart=no

[Install]
WantedBy=multi-user.target
EOF

gcloud config set project "${PROJECT_ID}"

# Upload environment variables
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mkdir -p $(dirname ${REMOTE_ENV_FILE}) $(dirname ${REMOTE_LOG_PATH})"
gcloud compute scp "${LOCAL_ENV_FILE}" "${SSH_USER}@${INSTANCE_NAME}:/tmp/${REMOTE_ENV_FILE}" --zone "${INSTANCE_ZONE}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mv /tmp/${REMOTE_SERVICE_NAME}.service ${SYSTEMD_SERVICE_FILE} && sudo systemctl daemon-reload"

# Upload systemd service file
gcloud compute scp "/tmp/${REMOTE_SERVICE_NAME}.service" "${SSH_USER}@${INSTANCE_NAME}:/tmp/${REMOTE_SERVICE_NAME}.service" --zone "${INSTANCE_ZONE}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo mv /tmp/${REMOTE_SERVICE_NAME}.service ${SYSTEMD_SERVICE_FILE} && sudo systemctl daemon-reload"

# Kill existing instance
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo systemctl stop ${REMOTE_SERVICE_NAME} || true"

# Upload new binary
gcloud compute scp "${LOCAL_BUILD_DIR}/${BINARY_NAME}" "${SSH_USER}@${INSTANCE_NAME}:${REMOTE_BINARY_PATH}" --zone "${INSTANCE_ZONE}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "chmod +x ${REMOTE_BINARY_PATH}"
rm "${LOCAL_BUILD_DIR}/${BINARY_NAME}"

# Start new binary
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo systemctl start ${REMOTE_SERVICE_NAME}"
gcloud compute ssh "${SSH_USER}@${INSTANCE_NAME}" --zone "${INSTANCE_ZONE}" -- "sudo systemctl status ${REMOTE_SERVICE_NAME}"

echo "${REMOTE_SERVICE_NAME} started"
