# xpense-trackr

A local-first personal finance app for expenses, budgets, goals, recurring transactions, investments, and spending trends. The application stores its data locally by default and does not require a cloud account or separately installed database.

## Run locally with Docker

Prerequisite: Docker Desktop (or Docker Engine with the Compose plugin).

```sh
docker compose up --build
```

Open [http://localhost:8080](http://localhost:8080). The app binds to the local machine, creates a SQLite database automatically, and persists data in the `xpense-data` Docker volume across restarts. To use another local port, set `APP_PORT`, for example `APP_PORT=8090 docker compose up --build`.

Stop with `Ctrl+C`; run `docker compose down` to stop and remove the container. The named data volume is kept. **Do not run `docker compose down -v` unless you intend to delete the app's stored data.**

## Back up and restore your data

Backups include a consistent SQLite snapshot, receipt files, and a versioned manifest. Create a backup while the app is running, then copy it to a location outside the Docker volume:

```sh
docker compose exec app /app/xpense-trackr backup --output /data/xpense-backup.zip
docker compose cp app:/data/xpense-backup.zip ./xpense-backup.zip
```

To restore, copy the archive into the volume, stop the app, and run the restore command. Restart the app afterward:

```sh
docker compose cp ./xpense-backup.zip app:/data/restore.zip
docker compose stop app
docker compose run --rm --no-deps app restore --input /data/restore.zip
docker compose start app
```

Restore validates the archive and database before changing data. It preserves the pre-restore database in `/data/pre-restore-<timestamp>.zip` and existing receipts in a `receipts.pre-restore-<timestamp>` directory. Keep an additional copy of the backup outside the machine running Docker.

## Run the services separately for development

Prerequisites: Go 1.27+ and Node.js 26+.

```sh
# Terminal 1: API (creates ./data/xpense-trackr.db on first run)
cd backend
go run ./cmd/expense-tracker

# Terminal 2: React development server
cd frontend/xpense-trackr-frontend
npm ci
npm start
```

The development frontend uses `http://localhost:8080` for API requests. The production build uses same-origin `/api` requests. The frontend test/build commands are `npm run test:ci` and `npm run build`.

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `APP_ADDR` | `127.0.0.1:8080` | HTTP listen address. Compose sets this to `:8080` inside its isolated container. |
| `DB_DRIVER` | `sqlite` | `sqlite` or `mysql`. |
| `DB_PATH` | `./data/xpense-trackr.db` | SQLite database path. |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME` | unset | Required MySQL connection settings when `DB_DRIVER=mysql`. |
| `APP_DATA_DIR` | `./data` | Local app data directory; receipts are stored under `receipts/`. |
| `WEB_DIR` | `./web` | Directory containing the production frontend build. |

The legacy `mysql/init.sql` file is not used by the local Compose setup. It contains demo fixtures for disposable legacy MySQL environments; do not use it to initialize an existing database.

## Import an existing MySQL installation

The importer previews and validates legacy categories, expenses, budgets, goals, investments, and recurring rules before writing to an initialized local SQLite database. Make and verify a backup first. Start the app once to create its schema, stop it, then run a preview and the explicit apply command:

```sh
export LEGACY_DB_PASS='your-legacy-database-password'
docker compose exec app /app/xpense-trackr backup --output /data/before-mysql-import.zip
docker compose cp app:/data/before-mysql-import.zip ./before-mysql-import.zip
docker compose stop app
docker compose run --rm --no-deps -e LEGACY_DB_PASS app import-mysql \
  --source-host host.docker.internal --source-port 3306 \
  --source-user xpense_trackr --source-database xpense_trackr_db \
  --currency USD
docker compose run --rm --no-deps -e LEGACY_DB_PASS app import-mysql \
  --source-host host.docker.internal --source-port 3306 \
  --source-user xpense_trackr --source-database xpense_trackr_db \
  --currency USD --apply
docker compose start app
```

The local source host address may differ outside Docker Desktop. The app must have been started once to initialize its schema and must remain stopped while importing. The importer refuses a non-empty destination and records successful source imports to prevent accidental repeats; it never changes the source. Because legacy records have no currency, `--currency` explicitly assigns one to all imported amounts. Legacy budgets become monthly envelopes, investments remain contribution history (not current market valuations), and recurring rows stay in their legacy tables for review. Legacy user/password records are intentionally excluded; this app is single-user and local-first.

## Current feature areas

- Local financial accounts and an exact-minor-unit transaction ledger for income, expenses, and same-currency transfers
- Merchant categorization rules for new uncategorized transactions
- Scheduled income and bills with daily balance projections; forecasts never write actual transactions
- Bill reminders, exact-match price-change signals, and confirmable, idempotent payment/skip actions
- Currency-separated cash-flow reports, CSV export, and ledger-derived historical net worth
- Manual asset/liability values with dated sources and historical net-worth rollups
- Savings-goal contributions with progress and monthly budget envelopes with optional rollover
- CSV transaction import with row validation, duplicate detection, preview, and transactional commit
- Local receipt attachments for PDF, PNG, and JPEG files (10 MB maximum)
- Expense and category management
- Category budgets and savings goals
- Recurring expenses and investments
- Manual investments
- A first-run dashboard, spending trends, and accessible tabular summaries
- An installable PWA app shell with static-only offline caching; financial API responses are never cached

The local account ledger is additive to the existing expense and investment screens while the product data model is being consolidated. Legacy `/api/...` endpoints remain available; the initial `/api/v1/...` path is a compatibility alias for them. The current account login/registration endpoints are intentionally unavailable (`501 Not Implemented`); this local single-user release does not claim to provide authentication.

Install the app from the browser's install menu after opening it over HTTP or HTTPS. Offline mode can load the app shell and cached static assets, but financial data and writes are unavailable until the local app server is reachable again.

The CSV importer accepts `account_id,type,amount,currency,date,payee,description,category_id,transfer_account_id` columns. Use **Import / export** to preview and validate before confirming an import; transaction amounts must use the selected currency's precision. Receipt files are stored alongside the SQLite database in the local data directory (or persistent Docker volume), not in the database itself. Back up both the database and `receipts/` together when backing up manually.

## Health endpoints

- `GET /healthz` reports that the HTTP server is running.
- `GET /readyz` checks that the database is reachable.

