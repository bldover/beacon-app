# Concert Manager Scripts

Management scripts for deploying and maintaining the Beacon Concert Manager backend service on GCP.

## Prerequisites

- GCP CLI (`gcloud`) installed and authenticated
- Environment files configured in `env/` directory (`vars_dev.sh`, `vars_prod.sh`)

## Scripts

### manage_svc.sh

Manages the Concert Manager service lifecycle on GCP compute instances.

**Usage:**
```bash
./manage_svc.sh [command] [environment]
```

**Commands:** `deploy`, `start`, `stop`, `restart`, `status`
**Environment:** `dev` (default), `prod`

**Examples:**
```bash
./manage_svc.sh deploy          # Deploy to dev
./manage_svc.sh start prod      # Start production service
./manage_svc.sh status dev      # Check dev service status
```

---

### getlogs.sh

Downloads log files from the deployed service.

**Usage:**
```bash
./getlogs.sh [environment] [lines]
```

**Arguments:**
- `environment` - `dev` (default) or `prod`
- `lines` - Optional number of lines from end of log (retrieves entire file if omitted)

**Examples:**
```bash
./getlogs.sh                # Get all dev logs
./getlogs.sh prod           # Get all prod logs
./getlogs.sh dev 100        # Get last 100 lines from dev
```

---

### manage_cron.sh

Manages cron jobs for automated cache refreshing on remote GCP servers.

**Usage:**
```bash
./manage_cron.sh [command] [environment]
```

**Commands:** `enable`, `disable`, `status`
**Environment:** `dev` (default), `prod`

**Examples:**
```bash
./manage_cron.sh enable          # Enable for dev
./manage_cron.sh enable dev      # Enable for dev
./manage_cron.sh enable prod     # Enable for prod
./manage_cron.sh status prod     # Check prod status
```

**Scheduled jobs (runs on remote server, calling localhost:3001):**
- Artist ranks refresh: Fridays at 5:00 AM EST
- Events refresh: Daily at 6:00 AM EST

---

### clone_prod_db_to_dev.sh

Clones production Firestore database to development environment.

**Usage:**
```bash
./clone_prod_db_to_dev.sh
```

**What it does:**
1. Clears development database
2. Exports production data
3. Imports into development

**⚠️ Warning:** Permanently deletes all dev database data. Requires confirmation.

---

## Common Workflows

**Deploy new version:**
```bash
./manage_svc.sh deploy dev      # Test in dev
./getlogs.sh dev 100            # Check logs
./manage_svc.sh deploy prod     # Deploy to prod
```

**Debug production:**
```bash
./manage_svc.sh status prod     # Check status
./getlogs.sh prod 500           # Get recent logs
./manage_svc.sh restart prod    # Restart if needed
```

**Refresh dev database:**
```bash
./clone_prod_db_to_dev.sh       # Clone prod data
./manage_svc.sh restart dev     # Restart service
```
