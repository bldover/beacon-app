#!/bin/bash

PROJECT_ID="concert-manager-dev"
INSTANCE_NAME="concert-manager-dev"
INSTANCE_ZONE="us-east1-b"
SSH_USER="dover_bradley"

BINARY_NAME="cm-server"
REMOTE_INSTALL_DIR="/opt/concert-manager"
REMOTE_BINARY_PATH="${REMOTE_INSTALL_DIR}/${BINARY_NAME}"
REMOTE_LOG_FILE="${REMOTE_BINARY_PATH}.log"
LOCAL_LOG_FILE="./${BINARY_NAME}.log"

set -e
set -x

gcloud config set project "${PROJECT_ID}"
gcloud compute scp "${SSH_USER}@${INSTANCE_NAME}:${REMOTE_LOG_FILE}" "${LOCAL_LOG_FILE}" --zone "${INSTANCE_ZONE}"
