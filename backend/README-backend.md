# Backend API

The Go service defaults to SQLite at `./data/xpense-trackr.db`. Set `DB_DRIVER=mysql` and `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, and `DB_NAME` to use MySQL. `DB_PATH` selects the SQLite file. The service applies additive schema migrations on startup.

The API is available at the legacy `/api/...` paths. `/api/v1/...` currently aliases the same handlers to provide a versioned path while preserving compatibility; it is not yet a separately versioned contract.

## Local accounts and transaction ledger

- `GET /api/accounts`: List accounts and calculated balances in each account's currency.
- `POST /api/accounts`: Create an account with `name`, `type`, `currency`, and `opening_balance_minor`.
- `GET /api/transactions`: List transactions. Optional query parameters: `account_id`, `type`, `from`, `to`, `limit` (maximum 100), and `offset`.
- `POST /api/transactions`: Create an `income`, `expense`, or same-currency `transfer`. Amounts use integer minor units and dates use `YYYY-MM-DD`.
- `GET /api/transactions/{id}`, `PUT /api/transactions/{id}`, `DELETE /api/transactions/{id}`: Read, update, and delete a transaction.
- `GET /api/transactions/export.csv`: Export ledger transactions in the importer's CSV format.
- `POST /api/transactions/import.csv`: Upload multipart field `file` to preview a transaction CSV. Add `?commit=true` only after review; invalid rows prevent the whole batch from committing, and repeated imports are deduplicated.
- `GET|POST /api/category-rules`, `PUT|DELETE /api/category-rules/{id}`: Manage case-insensitive payee-contains rules. Higher priority rules apply first, and an explicitly chosen transaction category is never overwritten.
- `GET|POST /api/scheduled-flows`, `PUT|DELETE /api/scheduled-flows/{id}`: Manage expected income and bills using `once`, `weekly`, `biweekly`, `monthly`, or `yearly` frequency.
- `GET /api/forecast/cashflow?days=90`: Project account balances from scheduled items (1–365 days). Projections are estimates and never create actual transactions.
- `GET /api/reports/cashflow?from=YYYY-MM-DD&to=YYYY-MM-DD`: Summarize bounded cash flow by currency and expense category; periods are limited to one year.
- `GET /api/reports/cashflow.csv`: Download the cash-flow report as CSV.
- `GET /api/net-worth`: Summarize positive and negative account balances separately by currency without assuming foreign-exchange rates.
- `GET|POST /api/receipts`: List receipts or upload multipart fields `transaction_id` and `file`; uploads allow PDF, PNG, and JPEG up to 10 MB.
- `GET /api/receipts/{id}/file`, `DELETE /api/receipts/{id}`: Download or delete a receipt attachment.

Currencies are currently limited to the supported set shown in the account form, with 0-, 2-, or 3-digit minor-unit precision handled explicitly. Cross-currency transfers and exchange-rate conversion are not yet supported. Net-worth summaries are current snapshots only; historical valuation, securities holdings, and live market prices are not provided.

Receipt metadata is stored in SQLite/MySQL; receipt content is stored in `APP_DATA_DIR/receipts` (default `./data/receipts`). Include both the database and this directory in a backup.

## Existing feature APIs

- Expenses: `GET|POST /api/expenses`, `GET|PUT|DELETE /api/expenses/{id}`, `GET /api/expenses/cumulative`
- Categories: `GET|POST /api/categories`, `GET|PUT|DELETE /api/categories/{id}`
- Budgets: `GET|POST /api/budgets`, `GET|PUT|DELETE /api/budgets/{id}`
- Goals: `GET|POST /api/goals`, `GET|PUT|DELETE /api/goals/{id}`
- Investments: `GET|POST /api/investments`, `GET|PUT|DELETE /api/investments/{id}`, `GET /api/investments/cumulative`
- Recurring items: `GET|POST /api/recurring/expenses`, `GET|PUT|DELETE /api/recurring/expenses/{id}`, with matching `/api/recurring/investments` paths
- Trends: `GET /api/trends/spending`

`GET /healthz` is a liveness check and `GET /readyz` checks database connectivity. Account login and registration endpoints currently return `501 Not Implemented`; the local single-user app does not claim to provide authentication.
