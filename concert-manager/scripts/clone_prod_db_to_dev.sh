#!/bin/bash

# Clone production Firestore database to development environment
# Usage: ./clone_prod_db_to_dev.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROD_ENV_FILE="$SCRIPT_DIR/env/vars_prod.sh"
DEV_ENV_FILE="$SCRIPT_DIR/env/vars_dev.sh"
EXPORT_DIR="/tmp/firestore-export-$(date +%Y%m%d-%H%M%S)"
EXPORT_BUCKET_SUFFIX="-firestore-exports"

show_usage() {
    echo "Usage: $0"
    echo ""
    echo "This script will:"
    echo "  1. Clear the development Firestore database"
    echo "  2. Export the entire production Firestore database"
    echo "  3. Import the production data into the development database"
}

load_environments() {
    if [[ ! -f "$PROD_ENV_FILE" ]]; then
        echo "Error: Production environment file not found: $PROD_ENV_FILE"
        exit 1
    fi

    if [[ ! -f "$DEV_ENV_FILE" ]]; then
        echo "Error: Development environment file not found: $DEV_ENV_FILE"
        exit 1
    fi

    echo "Loading production environment..."
    source "$PROD_ENV_FILE"
    PROD_PROJECT_ID="$PROJECT_ID"

    echo "Loading development environment..."
    source "$DEV_ENV_FILE"
    DEV_PROJECT_ID="$PROJECT_ID"

    if [[ "$PROD_PROJECT_ID" == "$DEV_PROJECT_ID" ]]; then
        echo "Error: Production and development project IDs are the same!"
        echo "This would result in data loss. Please check your environment files."
        exit 1
    fi

    echo "Configuration:"
    echo "  Production project:  $PROD_PROJECT_ID"
    echo "  Development project: $DEV_PROJECT_ID"
    echo "  Export directory:    $EXPORT_DIR"
}

check_gcp_access() {
    echo "Checking GCP access..."

    if ! gcloud config set project "$PROD_PROJECT_ID" &>/dev/null; then
        echo "Error: Cannot access production project $PROD_PROJECT_ID"
        exit 1
    fi

    if ! gcloud firestore databases list --project="$PROD_PROJECT_ID" &>/dev/null; then
        echo "Error: Cannot access Firestore in production project $PROD_PROJECT_ID"
        exit 1
    fi

    if ! gcloud config set project "$DEV_PROJECT_ID" &>/dev/null; then
        echo "Error: Cannot access development project $DEV_PROJECT_ID"
        exit 1
    fi

    if ! gcloud firestore databases list --project="$DEV_PROJECT_ID" &>/dev/null; then
        echo "Error: Cannot access Firestore in development project $DEV_PROJECT_ID"
        exit 1
    fi

    echo "GCP access verified for both projects"
}

ensure_export_bucket_exists() {
    local project_id="$1"
    local bucket_name="$2"

    echo "Setting up export bucket: $bucket_name"

    gcloud config set project "$project_id" &>/dev/null

    if ! gcloud storage ls "gs://$bucket_name" &>/dev/null; then
        echo "Creating export bucket..."
        gcloud storage buckets create "gs://$bucket_name" --project="$project_id"
    else
        echo "Export bucket already exists"
    fi
}

clear_dev_database() {
    echo "Clearing development Firestore database..."

    gcloud config set project "$DEV_PROJECT_ID" &>/dev/null

    echo "Deleting all documents from development database collections..."
    gcloud firestore bulk-delete --collection-ids='artists','events','venues' --quiet || {
        echo "Warning: Bulk delete failed, collections may already be empty"
    }

    echo "Development database cleared"
}

copy_export_data() {
    local prod_bucket="$1"
    local dev_bucket="$2"
    local export_dir="$3"

    echo "Copying export data from production to development bucket..."
    echo "Source: gs://$prod_bucket$export_dir"
    echo "Destination: gs://$dev_bucket$export_dir"

    gcloud storage cp -r "gs://$prod_bucket$export_dir" "gs://$dev_bucket$(dirname $export_dir)/" || {
        echo "Error: Failed to copy export data"
        exit 1
    }

    echo "Export data copied successfully"
}

export_production_data() {
    local bucket_name="$1"
    local export_path="$2"

    echo "Exporting production Firestore data..."
    echo "Export path: $export_path"

    gcloud config set project "$PROD_PROJECT_ID" &>/dev/null

    local operation_name
    operation_name=$(gcloud firestore export "$export_path" --format="value(name)")

    if [[ -z "$operation_name" ]]; then
        echo "Error: Failed to start export operation"
        exit 1
    fi

    echo "Export operation started: $operation_name"
    echo "Waiting for export to complete..."

    while true; do
        local state
        state=$(gcloud firestore operations describe "$operation_name" --format="value(done)" || echo "false")

        if [[ "$state" == "True" ]]; then
            echo "Export completed successfully"
            break
        elif [[ "$state" == "false" ]]; then
            echo "Export in progress..."
            sleep 30
        else
            echo "Error: Cannot determine export status"
            exit 1
        fi
    done
}

import_to_development() {
    local export_path="$1"

    echo "Importing data to development Firestore..."
    echo "Import path: $export_path"

    gcloud config set project "$DEV_PROJECT_ID" &>/dev/null

    local operation_name
    operation_name=$(gcloud firestore import "$export_path" --format="value(name)")

    if [[ -z "$operation_name" ]]; then
        echo "Error: Failed to start import operation"
        exit 1
    fi

    echo "Import operation started: $operation_name"
    echo "Waiting for import to complete..."

    while true; do
        local state
        state=$(gcloud firestore operations describe "$operation_name" --format="value(done)" || echo "false")

        if [[ "$state" == "True" ]]; then
            echo "Import completed successfully"
            break
        elif [[ "$state" == "false" ]]; then
            echo "Import in progress..."
            sleep 30
        else
            echo "Error: Cannot determine import status"
            exit 1
        fi
    done
}

main() {
    case "${1:-}" in
        --help|-h)
            show_usage
            exit 0
            ;;
        "")
            # No arguments, continue
            ;;
        *)
            echo "Error: Unknown argument $1"
            show_usage
            exit 1
            ;;
    esac

    echo "=== Firestore Database Clone: Production → Development ==="
    echo ""

    load_environments
    check_gcp_access

    echo ""
    echo "WARNING: This operation will:"
    echo "   1. PERMANENTLY DELETE all data in development database ($DEV_PROJECT_ID)"
    echo "   2. Replace it with data from production database ($PROD_PROJECT_ID)"
    echo ""
    read -p "Are you sure you want to continue? (yes/no): " confirm

    if [[ "$confirm" != "yes" ]]; then
        echo "Operation cancelled by user"
        exit 0
    fi

    echo ""
    echo "Starting database clone operation..."

    local prod_bucket="${PROD_PROJECT_ID}${EXPORT_BUCKET_SUFFIX}"
    local dev_bucket="${DEV_PROJECT_ID}${EXPORT_BUCKET_SUFFIX}"

    ensure_export_bucket_exists "$PROD_PROJECT_ID" "$prod_bucket"
    ensure_export_bucket_exists "$DEV_PROJECT_ID" "$dev_bucket"

    clear_dev_database

    local prod_export_path="gs://$prod_bucket$EXPORT_DIR"
    local dev_export_path="gs://$dev_bucket$EXPORT_DIR"

    export_production_data "$prod_bucket" "$prod_export_path"
    copy_export_data "$prod_bucket" "$dev_bucket" "$EXPORT_DIR"
    import_to_development "$dev_export_path"

    echo ""
    echo "Database clone completed successfully"
    echo "   Production ($PROD_PROJECT_ID) → Development ($DEV_PROJECT_ID)"
}

main "$@"
