# Xpense Trackr Backend API Reference

This reference describes the HTTP routes registered by the Go backend in
`backend/internal/routes.go`, along with the health routes in
`backend/cmd/expense-tracker/main.go`. It covers both the newer account/ledger
API and the older expense, investment, budget, goal, and recurring-item APIs.
Those generations intentionally have different payload and status conventions;
they are documented as implemented.

## Contents

- [Connecting to the API](#connecting-to-the-api)
- [Conventions and common responses](#conventions-and-common-responses)
- [Authentication](#authentication)
- [Accounts and transactions](#accounts-and-transactions)
- [Automatic category rules](#automatic-category-rules)
- [Scheduled flows, bills, and forecasts](#scheduled-flows-bills-and-forecasts)
- [Reports and net worth](#reports-and-net-worth)
- [Manual assets and portfolio holdings](#manual-assets-and-portfolio-holdings)
- [Planning goals and budget envelopes](#planning-goals-and-budget-envelopes)
- [CSV import and export](#csv-import-and-export)
- [Receipt files](#receipt-files)
- [Legacy expense and category APIs](#legacy-expense-and-category-apis)
- [Legacy budgets and goals](#legacy-budgets-and-goals)
- [Legacy investments and recurring items](#legacy-investments-and-recurring-items)
- [Spending trends](#spending-trends)
- [Health and readiness](#health-and-readiness)

## Connecting to the API

The local server listens on `http://127.0.0.1:8080` by default. Docker Compose
publishes it at `http://localhost:8080`. The development frontend uses this
same origin for API requests.

All API routes are available under `/api/...`. The server also aliases
`/api/v1/<path>` to `/api/<path>` for compatibility. This is a path alias, not
a separately versioned contract: both prefixes use the same handlers and
payloads. For example, `GET /api/v1/accounts` is handled as
`GET /api/accounts`.

Unless noted otherwise:

- Send JSON request bodies with `Content-Type: application/json`.
- JSON field names are `snake_case`.
- `{id}` denotes a positive numeric resource ID.
- Date-only values use `YYYY-MM-DD`; month values use `YYYY-MM`.
- The newer ledger APIs represent money as integer minor units, such as
  `12345` for USD 123.45. They never convert between currencies.
- Older APIs use their legacy decimal `amount` fields. See their individual
  schemas below.
- Lists are returned as JSON arrays; an empty list is generally `[]` for the
  newer APIs. Some older handlers encode a nil Go slice as `null`.

## Conventions and common responses

| Status | Meaning in this API |
| --- | --- |
| `200 OK` | Successful read or update; some legacy create/delete routes also return 200. |
| `201 Created` | Resource or receipt successfully created. |
| `204 No Content` | Successful delete or update with no response body. |
| `400 Bad Request` | Malformed JSON, invalid query/path values, or invalid resource data. |
| `404 Not Found` | Resource does not exist. |
| `409 Conflict` | A scheduled occurrence conflicts with one already processed, or a withdrawal exceeds the goal balance. |
| `500 Internal Server Error` | Unexpected storage, encoding, or server error. |
| `501 Not Implemented` | Authentication is intentionally unavailable. |
| `503 Service Unavailable` | Database readiness check failed. |

Most handler errors use Go's `http.Error`: the response body is plain text
(with a trailing newline), not a JSON error object. Error messages are listed
when the handler defines a stable message; validation messages may vary by
field or include a decoding error.

The newer ledger, planning, and receipt handlers set `Content-Type:
application/json` for JSON responses. Several legacy handlers encode JSON
without explicitly setting that header. Clients should parse successful
response bodies by route contract rather than rely on a uniform header across
legacy and newer routes.

Unsupported paths return `404`; unsupported methods are not registered. No
authentication middleware is applied to the finance routes in this local,
single-user application.

## Authentication

These routes are registered for compatibility, but account login and
registration are currently unimplemented.

### `POST /api/auth/login`

Request:

```json
{
  "email": "person@example.com",
  "password": "example-password"
}
```

Response: `501 Not Implemented`, plain-text body indicating the operation is
not implemented. Malformed JSON returns `400`. The handler's `401` invalid
credentials response is only relevant if authentication is implemented; the
current authentication service returns `501`.

### `POST /api/auth/register`

Request:

```json
{
  "username": "person",
  "email": "person@example.com",
  "password": "example-password"
}
```

Response: `501 Not Implemented`, plain-text body. Malformed JSON returns `400`.
If registration is implemented in a future version, the handler is structured
to return `201 Created` with a plain-text success message.

## Accounts and transactions

The account ledger is the current source for account balances, transaction
history, cash-flow reporting, and several planning features.

### Accounts

#### `GET /api/accounts`

No request body or query parameters.

`200 OK` returns accounts ordered with unarchived accounts first and then by
name. `balance_minor` is calculated from the opening balance and ledger
activity.

```json
[
  {
    "id": 1,
    "name": "Everyday checking",
    "type": "checking",
    "currency": "USD",
    "opening_balance_minor": 125000,
    "archived": false,
    "balance_minor": 118450,
    "created_at": "2026-01-10T12:00:00Z",
    "updated_at": "2026-01-10T12:00:00Z"
  }
]
```

Errors: `500` if accounts cannot be read.

#### `POST /api/accounts`

Request fields:

| Field | Type | Required | Rules |
| --- | --- | --- | --- |
| `name` | string | Yes | Trimmed; 1–120 characters. |
| `type` | string | Yes | `cash`, `checking`, `savings`, `investment`, `credit_card`, or `loan`. |
| `currency` | string | Yes | One of the supported three-letter currency codes; normalized to uppercase. |
| `opening_balance_minor` | integer | No | Defaults to zero. A negative opening balance is allowed. |

Example:

```json
{
  "name": "Everyday checking",
  "type": "checking",
  "currency": "USD",
  "opening_balance_minor": 125000
}
```

`201 Created` returns the created account, including its generated `id`,
timestamps, and `balance_minor`. Errors: `400` for malformed/invalid payloads;
`500` for storage errors.

### Transactions

#### `GET /api/transactions`

Optional query parameters:

| Parameter | Default | Rules |
| --- | --- | --- |
| `account_id` | none | Positive integer; includes transfers where the account is either side. |
| `type` | none | `income`, `expense`, or `transfer`. |
| `from` | none | Inclusive start date in `YYYY-MM-DD`. |
| `to` | none | Inclusive end date in `YYYY-MM-DD`. |
| `limit` | `50` | Integer from 1 to 100. |
| `offset` | `0` | Integer from 0 to 10,000. |

`200 OK` returns matching transactions sorted newest first by date and ID.
Errors: `400` for an invalid filter or pagination value; `500` for a storage
error.

#### `POST /api/transactions`

Request fields:

| Field | Type | Required | Rules |
| --- | --- | --- | --- |
| `account_id` | integer | Yes | Existing account. |
| `type` | string | Yes | `income`, `expense`, or `transfer`. |
| `amount_minor` | integer | Yes | Must be positive. |
| `currency` | string | Yes | Supported code matching the source account currency. |
| `date` | string | Yes | Valid `YYYY-MM-DD` date. |
| `transfer_account_id` | integer or null | For transfer | Existing, different account with the same currency. Only valid for `transfer`. |
| `category_id` | integer or null | No | Existing category; transfers cannot have a category. |
| `payee` | string | No | Up to 180 characters. |
| `description` | string | No | Up to 500 characters. |

For `income` and `expense`, omit `transfer_account_id`. For a transfer, include
`transfer_account_id` and omit `category_id`.

```json
{
  "account_id": 1,
  "type": "expense",
  "amount_minor": 4250,
  "currency": "USD",
  "date": "2026-11-08",
  "category_id": 4,
  "payee": "Market",
  "description": "Groceries"
}
```

`201 Created` returns the transaction with generated `id` and timestamps. If no
category is provided, an enabled matching category rule may populate
`category_id` and `category_rule_id`. The internal CSV fingerprint is not
returned.

Errors: `400` for malformed/unknown JSON fields, invalid amount/type/date,
currency mismatch, invalid account/category, or invalid transfer. `500` for
unexpected storage errors.

#### `GET /api/transactions/{id}`

No request body. `200 OK` returns one transaction. Errors: `400` for an invalid
ID; `404` with `Requested finance record not found` if no transaction exists;
`500` for other errors.

#### `PUT /api/transactions/{id}`

Accepts the same fields and validation as `POST /api/transactions`; the
request represents the transaction's complete intended values. The path ID
selects the transaction and the returned `id` is set from that path.

`200 OK` returns the updated transaction. Errors: `400` invalid ID or
transaction; `404` if no transaction exists; `500` otherwise.

#### `DELETE /api/transactions/{id}`

`204 No Content` on success. Errors: `400` invalid ID; `404` if no transaction
exists; `500` otherwise.

### Account-balance effects

Balances begin at `opening_balance_minor`. Income adds to the transaction's
account; expenses and outgoing transfers subtract from it; incoming transfers
add to the receiving account. Transactions and transfers cannot cross
currencies. Archived accounts remain in the API and their balances are still
calculated.

## Automatic category rules

Rules apply to new or updated ledger transactions only when no explicit
category was supplied. Matching is a case-insensitive substring search against
the payee; transfer transactions are not auto-categorized. Enabled rules are
evaluated by descending `priority`, then ascending ID. An explicit category
always takes precedence.

Rule representation:

```json
{
  "id": 1,
  "name": "Market purchases",
  "match_value": "market",
  "category_id": 4,
  "priority": 10,
  "enabled": true,
  "created_at": "2026-01-10T12:00:00Z",
  "updated_at": "2026-01-10T12:00:00Z"
}
```

### `GET /api/category-rules`

Returns all rules, including disabled rules, sorted by priority descending
then ID ascending. Success: `200 OK`; errors: `500`.

### `POST /api/category-rules`

Required: non-empty `name` (maximum 120 characters), non-empty `match_value`
(maximum 180 characters), and an existing numeric `category_id`. Optional:
integer `priority` and boolean `enabled` (defaults to true).

```json
{
  "name": "Market purchases",
  "match_value": "market",
  "category_id": 4,
  "priority": 10,
  "enabled": true
}
```

Success: `201 Created` with the created rule. Errors: `400` invalid JSON or
rule; `404` if `category_id` does not exist; `500` otherwise.

### `PUT /api/category-rules/{id}`

Accepts the same rule fields and validation as create. Success: `200 OK` with
the updated rule. Errors: `400` invalid ID/body/rule; `404` if the rule or
selected category does not exist; `500` otherwise.

### `DELETE /api/category-rules/{id}`

Success: `204 No Content`. Errors: `400` invalid ID; `404` if no rule exists;
`500` otherwise.

## Scheduled flows, bills, and forecasts

Scheduled flows describe expected income or expenses. They are estimates until
an occurrence is explicitly recorded. Only recording an occurrence writes a
transaction to the ledger.

### Scheduled flows

Flow representation:

```json
{
  "id": 2,
  "account_id": 1,
  "category_id": 4,
  "label": "Rent",
  "type": "expense",
  "amount_minor": 180000,
  "currency": "USD",
  "frequency": "monthly",
  "next_date": "2026-12-01",
  "end_date": null,
  "reminder_days": 3,
  "enabled": true,
  "created_at": "2026-11-08T12:00:00Z",
  "updated_at": "2026-11-08T12:00:00Z"
}
```

`category_id` and `end_date` may be omitted or null. `category_id` is only
allowed for expenses. Supported frequencies are `once`, `weekly`, `biweekly`,
`monthly`, and `yearly`. `reminder_days` is from 0 to 30. The account and
currency must match.

#### `GET /api/scheduled-flows`

Returns all flows ordered by next date then ID. Success: `200 OK`; errors:
`500`.

#### `POST /api/scheduled-flows`

Requires `account_id`, `label`, `type` (`income` or `expense`),
`amount_minor` (positive integer), `currency`, `frequency`, and `next_date`.
Optional fields: `category_id`, `end_date`, `reminder_days` (defaults to 3),
and `enabled` (defaults to true).

Success: `201 Created` with the created flow. Errors: `400` invalid body or
flow; `404` if a referenced account/category does not exist; `500` otherwise.

#### `PUT /api/scheduled-flows/{id}`

Accepts the same fields and validation as create. Success: `200 OK` with the
updated flow. Errors: `400` invalid ID/body/flow; `404` if the flow or a
referenced record does not exist; `500` otherwise.

#### `DELETE /api/scheduled-flows/{id}`

Success: `204 No Content`. Errors: `400` invalid ID; `404` if no flow exists;
`500` otherwise.

### Scheduled occurrences

#### `POST /api/scheduled-flows/{id}/occurrences`

Request:

```json
{
  "scheduled_date": "2026-12-01",
  "action": "record",
  "amount_minor": 180000
}
```

`scheduled_date` must equal the flow's current `next_date`. `action` is
`record` or `skip`. `amount_minor` is optional for `record` (zero or omitted
uses the scheduled amount); it must be zero for `skip`.

For `record`, the handler creates exactly one income/expense ledger
transaction, stores the occurrence, and advances the schedule atomically. For
`skip`, it stores the skipped occurrence and advances the schedule without
creating a ledger transaction. A `once` flow is disabled after processing.
Repeating an identical confirmation is idempotent and returns the existing
occurrence.

Response: `200 OK` with an occurrence:

```json
{
  "id": 9,
  "scheduled_flow_id": 2,
  "scheduled_date": "2026-12-01",
  "status": "record",
  "transaction_id": 35,
  "actual_amount_minor": 180000,
  "created_at": "2026-11-08T12:00:00Z"
}
```

`transaction_id` is omitted for skipped occurrences. Errors: `400` invalid
body/action/date, amount, or occurrence not matching the flow's next date;
`404` if the flow does not exist; `409` if the date was already processed
with different details; `500` otherwise.

### Bills

#### `GET /api/bills`

Returns enabled scheduled expense flows with account names and due indicators.
Each item contains the scheduled-flow fields plus `account_name`,
`status` (`overdue`, `due_soon`, or `upcoming`), `days_until_due`,
`has_price_change`, and, when there is a signal, `previous_amount_minor`,
`latest_amount_minor`, and `change_percent`.

Price-change comparisons use the two most recent expense transactions in the
same currency whose payee or description exactly matches the bill label,
ignoring case and surrounding spaces. A signal is set when the change is at
least 5%. Success: `200 OK`; errors: `500`.

### Cash-flow forecast

#### `GET /api/forecast/cashflow`

Optional query parameter: `days`, default `90`, integer from 1 to 365.

`200 OK` response:

```json
{
  "start_date": "2026-11-08",
  "end_date": "2027-02-06",
  "days": 90,
  "accounts": [
    {
      "account_id": 1,
      "name": "Everyday checking",
      "currency": "USD",
      "current_balance_minor": 118450,
      "projected_balance_minor": 105000,
      "points": [
        {"date": "2026-11-08", "balance_minor": 118450},
        {"date": "2026-11-09", "balance_minor": 118450}
      ]
    }
  ],
  "events": [
    {
      "scheduled_flow_id": 2,
      "account_id": 1,
      "label": "Rent",
      "type": "expense",
      "amount_minor": 180000,
      "currency": "USD",
      "date": "2026-12-01"
    }
  ]
}
```

The values are projections only; this endpoint never writes transactions.
Errors: `400` if `days` is not an integer from 1 to 365; `500` otherwise.

## Reports and net worth

### Cash-flow report

#### `GET /api/reports/cashflow`

Optional query parameters `from` and `to`, each `YYYY-MM-DD`. Defaults are
the first day of the current month and today. The inclusive range must be
valid, non-reversed, and no longer than one calendar year.

`200 OK` response:

```json
{
  "from": "2026-11-01",
  "to": "2026-11-08",
  "currencies": [
    {
      "currency": "USD",
      "income_minor": 250000,
      "expense_minor": 72500,
      "net_minor": 177500,
      "categories": [
        {
          "category_id": 4,
          "category_name": "Groceries",
          "amount_minor": 72500
        }
      ]
    }
  ]
}
```

Currencies are reported separately. `categories` contains expense totals;
uncategorized expenses have an empty `category_name` and omit `category_id`.
Transfers are excluded from cash-flow totals. Errors: `400` invalid date/range;
`500` otherwise.

#### `GET /api/reports/cashflow.csv`

Accepts the same `from` and `to` parameters and validation. Success is
`200 OK`, `Content-Type: text/csv; charset=utf-8`, with a download filename of
`cashflow-report.csv`. Columns:

```csv
from,to,currency,type,category,amount_minor
```

The export emits income, expense-by-category, and net rows. Errors: `400`
invalid date/range; `500` if the report cannot be built.

### Period-comparison insights

#### `GET /api/reports/insights`

Accepts the same date parameters and range rules as the cash-flow report. The
backend compares the selected inclusive period to the immediately preceding
period of equal length, separately per currency.

`200 OK` response:

```json
{
  "from": "2026-11-01",
  "to": "2026-11-08",
  "previous_from": "2026-10-24",
  "previous_to": "2026-10-31",
  "insights": [
    {
      "kind": "spending_change",
      "title": "Spending increased by 15% in USD",
      "explanation": "Recorded spending changed from 500.00 to 575.00 compared with the previous period.",
      "currency": "USD",
      "current_amount_minor": 57500,
      "previous_amount_minor": 50000,
      "change_percent": 15
    }
  ]
}
```

`category_name` is present for category-specific insights. The backend filters
out small changes; category insights are limited to the three largest absolute
changes per currency. Errors: `400` invalid date/range; `500` otherwise.

### Current net worth

#### `GET /api/net-worth`

No parameters. `200 OK` returns an `as_of` date and `currencies` array. Each
currency contains `assets_minor`, `liabilities_minor`, `net_worth_minor`,
`portfolio_minor`, `accounts`, `manual_items`, and `holdings`. Account balances
are calculated as of the report date. Manual valuations and portfolio
valuations are included only when dated on or before that date. No currency
conversion or live market prices are applied.

```json
{
  "as_of": "2026-11-08",
  "currencies": [
    {
      "currency": "USD",
      "assets_minor": 150000,
      "liabilities_minor": 20000,
      "net_worth_minor": 130000,
      "portfolio_minor": 50000,
      "accounts": [
        {
          "account_id": 1,
          "name": "Everyday checking",
          "type": "checking",
          "balance_minor": 100000
        }
      ],
      "manual_items": [],
      "holdings": []
    }
  ]
}
```

Success: `200 OK`; errors: `500`.

### Net-worth history

#### `GET /api/net-worth/history`

Optional `months` query parameter; default `12`, integer from 1 to 36.
Returns `months` month-end snapshots; the current month's entry uses today's
date rather than a future month-end.

```json
{
  "months": 2,
  "entries": [
    {"as_of": "2026-10-31", "currencies": []},
    {"as_of": "2026-11-08", "currencies": []}
  ]
}
```

Each entry's `currencies` objects use the same structure as
`GET /api/net-worth`. Success: `200 OK`; errors: `400` if `months` is outside
1–36 or not an integer; `500` otherwise.

## Manual assets and portfolio holdings

### Manual assets and liabilities

#### `GET /api/assets`

Returns manual items sorted by name. Each item includes `id`, `name`, `kind`,
`currency`, `latest_value_minor`, `valuation_date`, `valuation_source`,
`has_valuation`, and timestamps. Items without a valuation have
`has_valuation: false`, zero `latest_value_minor`, and empty valuation
date/source values.

Success: `200 OK`; errors: `500`.

#### `POST /api/assets`

Required fields: `name` (1–120 characters), `kind` (`asset` or `liability`),
supported `currency`, non-negative `value_minor`, and valid `date`. Optional
`source` defaults to `"manual"` and is limited to 100 characters.

```json
{
  "name": "Car",
  "kind": "asset",
  "currency": "USD",
  "value_minor": 1800000,
  "date": "2026-11-08",
  "source": "Manual estimate"
}
```

Success: `201 Created` with the created asset and its initial valuation.
Errors: `400` invalid JSON or fields; `500` otherwise.

#### `PUT /api/assets/{id}`

Request body contains `name` (1–120 characters) and `kind` (`asset` or
`liability`). This changes metadata only; it does not change the currency or
valuation history.

Success: `204 No Content`; errors: `400` invalid ID/body; `404` if the asset
does not exist; `500` otherwise.

#### `DELETE /api/assets/{id}`

Deletes the asset and all of its valuations. Success: `204 No Content`;
errors: `400` invalid ID; `404` if it does not exist; `500` otherwise.

#### `GET /api/assets/{id}/valuations`

Returns valuations for the asset, newest date first. Each valuation contains
`id`, `asset_id`, `value_minor`, `date`, `source`, and `created_at`.
Success: `200 OK`; errors: `400` invalid ID; `404` if the asset does not
exist; `500` otherwise.

#### `POST /api/assets/{id}/valuations`

Requires non-negative `value_minor` and valid `date`. Optional `source`
defaults to `"manual"` and is limited to 100 characters.

```json
{
  "value_minor": 1750000,
  "date": "2026-11-08",
  "source": "Updated estimate"
}
```

Success: `201 Created` with the valuation. Errors: `400` invalid ID/body or
valuation; `404` if the asset does not exist; `500` otherwise.

### Portfolio holdings

Supported `asset_class` values are `cash`, `equity`, `fixed_income`,
`real_estate`, `commodity`, `crypto`, and `other`. `quantity` is a positive
decimal string with up to 8 fractional digits. `ticker` may contain letters,
digits, `.`, `^`, `_`, or `-`, up to 16 characters.

#### `GET /api/portfolio/holdings`

Returns holdings sorted by name. Each includes `id`, `name`, `ticker`,
`asset_class`, `currency`, `quantity`, `cost_basis_minor`,
`latest_value_minor`, `valuation_date`, `valuation_source`, `has_valuation`,
and timestamps. A holding is returned even if it has no valuation.

Success: `200 OK`; errors: `500`.

#### `POST /api/portfolio/holdings`

Required fields: `name`, `asset_class`, supported `currency`, positive
`quantity`, non-negative `cost_basis_minor`, non-negative `value_minor`, and
valid valuation `date`. Optional `ticker` defaults to empty and optional
`source` defaults to `"manual"`.

```json
{
  "name": "Index fund",
  "ticker": "INDEX",
  "asset_class": "equity",
  "currency": "USD",
  "quantity": "12.5",
  "cost_basis_minor": 250000,
  "value_minor": 275000,
  "date": "2026-11-08",
  "source": "Broker statement"
}
```

Creates both the holding and its initial valuation. Success: `201 Created`
with the holding. Errors: `400` invalid JSON or values; `500` otherwise.

#### `PUT /api/portfolio/holdings/{id}`

Requires `name`, `ticker`, `asset_class`, `quantity`, and
`cost_basis_minor`. The currency cannot be changed through this endpoint.
Success: `204 No Content`; errors: `400` invalid ID/body/holding; `404` if no
holding exists; `500` otherwise.

#### `DELETE /api/portfolio/holdings/{id}`

Deletes the holding and all its valuations. Success: `204 No Content`;
errors: `400` invalid ID; `404` if no holding exists; `500` otherwise.

#### `GET /api/portfolio/holdings/{id}/valuations`

Returns valuations newest date first. Each contains `id`, `holding_id`,
`value_minor`, `date`, `source`, and `created_at`. Success: `200 OK`; errors:
`400` invalid ID; `404` if no holding exists; `500` otherwise.

#### `POST /api/portfolio/holdings/{id}/valuations`

Requires non-negative `value_minor` and valid `date`. Optional `source`
defaults to `"manual"` and is limited to 100 characters.

```json
{
  "value_minor": 280000,
  "date": "2026-11-08",
  "source": "Broker statement"
}
```

Success: `201 Created` with the valuation. Errors: `400` invalid ID/body or
valuation; `404` if no holding exists; `500` otherwise.

## Planning goals and budget envelopes

These planning endpoints use the newer ledger's minor-unit currency model.
They are separate from the legacy `/api/goals` and `/api/budgets` resources.

### Savings goals

Goal representation:

```json
{
  "id": 3,
  "name": "Emergency fund",
  "currency": "USD",
  "target_amount_minor": 1000000,
  "target_date": "2027-12-31",
  "current_amount_minor": 150000,
  "remaining_minor": 850000,
  "required_monthly_minor": 70834,
  "created_at": "2026-11-08T12:00:00Z",
  "updated_at": "2026-11-08T12:00:00Z"
}
```

`target_date` and `required_monthly_minor` may be omitted when not applicable.
Progress is calculated from contribution history; it is not directly editable
on this goal resource.

#### `GET /api/planning/goals`

Returns goals sorted by target date, with undated goals last. Success:
`200 OK`; errors: `500`.

#### `POST /api/planning/goals`

Requires `name` (1–120 characters), supported `currency`, and positive
`target_amount_minor`. Optional `target_date` must be a valid `YYYY-MM-DD`.

```json
{
  "name": "Emergency fund",
  "currency": "USD",
  "target_amount_minor": 1000000,
  "target_date": "2027-12-31"
}
```

Success: `201 Created` with the goal and calculated progress fields. Errors:
`400` invalid JSON or fields; `500` otherwise.

#### `DELETE /api/planning/goals/{id}`

Deletes the goal and its contribution history. Success: `204 No Content`;
errors: `400` invalid ID; `404` if the goal does not exist; `500` otherwise.

### Goal contributions

#### `GET /api/planning/goals/{id}/contributions`

Returns contribution entries newest date first:

```json
[
  {
    "id": 8,
    "goal_id": 3,
    "amount_minor": 25000,
    "direction": "contribution",
    "date": "2026-11-08",
    "note": "Monthly savings",
    "created_at": "2026-11-08T12:00:00Z"
  }
]
```

Success: `200 OK`; errors: `400` invalid ID; `500` if the entries cannot be
read. A missing goal currently produces `500` on this GET route.

#### `POST /api/planning/goals/{id}/contributions`

Requires positive `amount_minor`, `direction` (`contribution` or
`withdrawal`), and valid `date`. Optional `note` is limited to 300 characters.
The goal's currency is used; no currency field is accepted.

```json
{
  "amount_minor": 25000,
  "direction": "contribution",
  "date": "2026-11-08",
  "note": "Monthly savings"
}
```

Success: `201 Created` with the contribution. Errors: `400` invalid ID/body or
entry; `404` if the goal does not exist; `409` if a withdrawal exceeds the
available contribution balance; `500` otherwise.

### Budget envelopes

Envelope summary fields include the stored envelope fields plus `period`,
`spent_minor`, `rollover_minor`, and `available_minor`.

#### `GET /api/planning/envelopes`

Optional `period` query parameter uses `YYYY-MM`; it defaults to the current
UTC month. Returns each envelope with its spending summary for that period.
When rollover is enabled, prior monthly underspend or overspend contributes to
`rollover_minor`. `available_minor` is monthly limit plus rollover minus
spending.

Success: `200 OK`; errors: `400` if `period` is invalid; `500` otherwise.

#### `POST /api/planning/envelopes`

Requires an existing `category_id`, supported `currency`, positive
`monthly_limit_minor`, and optional boolean `rollover` (defaults to false).
There can be only one envelope per category/currency pair.

```json
{
  "category_id": 4,
  "currency": "USD",
  "monthly_limit_minor": 60000,
  "rollover": true
}
```

Success: `201 Created` with the stored envelope. Errors: `400` invalid JSON or
fields; `500` for storage errors, including a missing category or duplicate
category/currency pair.

#### `PUT /api/planning/envelopes/{id}`

Accepts the same fields as create. Success: `200 OK` with the updated envelope
and path `id`. Errors: `400` invalid ID/body/fields; `404` if the envelope
does not exist; `500` for storage errors.

#### `DELETE /api/planning/envelopes/{id}`

Success: `204 No Content`; errors: `400` invalid ID; `404` if the envelope
does not exist; `500` otherwise.

## CSV import and export

### Export transactions

#### `GET /api/transactions/export.csv`

Downloads all ledger transactions in date/ID ascending order. Success:
`200 OK`, `Content-Type: text/csv; charset=utf-8`, filename `transactions.csv`.

Columns:

```csv
date,account_id,type,amount,currency,payee,description,category_id,transfer_account_id
```

`amount` is formatted as a decimal string using the currency's minor-unit
precision. Nullable category and transfer account IDs are empty cells.
Errors: `500` if the export cannot be read.

### Import transactions

#### `POST /api/transactions/import.csv`

Upload as `multipart/form-data` with a file field named `file`. Maximum file
size is 5 MiB and maximum data rows is 5,000. The CSV header must include
`account_id`, `type`, `amount`, `currency`, and `date`; optional columns are
`payee`, `description`, `category_id`, and `transfer_account_id`. Additional
columns are accepted.

Example CSV:

```csv
account_id,type,amount,currency,date,payee,description,category_id,transfer_account_id
1,expense,42.50,USD,2026-11-08,Market,Groceries,4,
```

Without a query parameter, the endpoint previews and does not write any
transactions. Use `?commit=true` only after reviewing the preview. Rows are
reported as `ready`, `duplicate`, or `invalid`; duplicates are identified
within the file and against prior imports. An import fingerprint prevents
re-importing identical transaction data.

Example response:

```json
{
  "total": 1,
  "ready": 1,
  "duplicates": 0,
  "invalid": 0,
  "imported": 0,
  "rows": [
    {
      "line": 2,
      "status": "ready",
      "transaction": {
        "account_id": 1,
        "type": "expense",
        "amount_minor": 4250,
        "currency": "USD",
        "date": "2026-11-08",
        "payee": "Market",
        "description": "Groceries",
        "category_id": 4
      }
    }
  ]
}
```

Preview success: `200 OK`, including rows marked invalid. Commit success:
`200 OK`; committed rows use status `imported`, `imported` is incremented, and
`ready` becomes zero. If commit is requested with invalid rows, the batch is
not committed and the endpoint returns `400 Bad Request` with the same JSON
report. Other `400` responses use plain text for malformed/missing CSV, a
file over 5 MiB, or more than 5,000 rows. `500` indicates a storage or
unexpected import error. A batch that has no ready rows is not written.

## Receipt files

Receipt metadata is stored in the database; file content is stored under the
configured application data directory's `receipts/` folder. Supported content
types are PDF, PNG, and JPEG; maximum file size is 10 MiB.

### `GET /api/receipts`

Optional `transaction_id` filters receipts to one transaction and must be a
positive integer. Returns newest receipts first. The response includes `id`,
`transaction_id`, `file_name`, `media_type`, `size_bytes`, and `created_at`;
the private storage key is never returned.

Success: `200 OK`; errors: `400` invalid filter; `500` otherwise.

### `POST /api/receipts`

Upload `multipart/form-data` with:

| Field | Type | Required |
| --- | --- | --- |
| `transaction_id` | positive integer | Yes |
| `file` | PDF, PNG, or JPEG file | Yes |

Success: `201 Created` with the receipt metadata. Errors: `400` malformed
multipart data, missing/invalid fields, empty or oversized file, unsupported
file content, or invalid filename; `404` if the transaction does not exist;
`500` for file/database errors.

### `GET /api/receipts/{id}/file`

Downloads the file as an attachment, with its stored media type and original
filename. Success: `200 OK` with binary content (range requests are supported
by the HTTP file-serving implementation). Errors: `400` invalid ID; `404` if
receipt metadata does not exist; `500` if the stored file cannot be opened.

### `DELETE /api/receipts/{id}`

Deletes the receipt metadata and file. Success: `204 No Content`; errors:
`400` invalid ID; `404` if the receipt does not exist; `500` if file or
database operations fail.

## Legacy expense and category APIs

These endpoints retain their older data model. Their request validation is
less strict than the newer ledger APIs; database and JSON decoding errors are
returned by the handler. Legacy error bodies are plain text.

### Expenses

Legacy expense representation:

```json
{
  "id": 12,
  "category_id": "4",
  "amount": "42.5",
  "description": "Groceries",
  "date": "2026-11-08"
}
```

`category_id` and `amount` are JSON strings because of the legacy struct's
encoding tags. `id` is numeric.

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/expenses` | None | `200`, JSON array of expenses. | `500` |
| `POST /api/expenses` | `category_id`, `amount`, `description`, `date` | `201`, created expense. | `400`, `500` |
| `GET /api/expenses/{id}` | None | `200`, one expense. | `400`, `404`, `500` |
| `PUT /api/expenses/{id}` | Expense object | `200`, updated expense. | `400`, `404`, `500` |
| `DELETE /api/expenses/{id}` | None | `204`, no body. | `400`, `404`, `500` |

### Categories

Category representation: `{"id": 4, "name": "Groceries"}`.

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/categories` | None | `200`, JSON array of categories. | `500` |
| `POST /api/categories` | `{"name":"Groceries"}` | `201`, created category. | `400`, `500` |
| `GET /api/categories/{id}` | None | `200`, one category. | `400`, `404`, `500` |
| `PUT /api/categories/{id}` | Category object | `200`, updated category. | `400`, `404`, `500` |
| `DELETE /api/categories/{id}` | None | `204`, no body. | `400`, `404`, `500` |

Category names are not validated by the handler beyond JSON decoding. Existing
database constraints may reject a value; such failures return `500`.

### Cumulative expenses

#### `GET /api/expenses/cumulative?lastMonths={count}`

`lastMonths` is required and must be a positive integer; there is no explicit
upper bound. The response is a cumulative total per category per month,
starting at the oldest returned month and accumulating forward:

```json
[
  {"month": "2026-10", "category": "Groceries", "amount": 125.5},
  {"month": "2026-11", "category": "Groceries", "amount": 168.0}
]
```

Success: `200 OK`; errors: `400` missing/invalid `lastMonths`; `500` if
categories or expense data cannot be read.

## Legacy budgets and goals

These legacy resources use decimal `amount`/`target_amount` values and remain
separate from ledger envelopes and minor-unit planning goals.

### Budgets

Budget representation: `{"id": 1, "category_id": 4, "amount": 600}`.

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/budgets` | None | `200`, JSON array. | `500` |
| `POST /api/budgets` | `category_id` (number), `amount` (number) | `200`, created budget. | `400`, `500` |
| `GET /api/budgets/{id}` | None | `200`, one budget. | `400`, `404` |
| `PUT /api/budgets/{id}` | Budget object | `200`, updated budget. | `400`, `404`, `500` |
| `DELETE /api/budgets/{id}` | None | `200`, empty body. | `400`, `404`, `500` |

The legacy handler does not validate amount/category semantics. `POST` returns
200 rather than 201, and successful delete returns 200 rather than 204.

### Legacy savings goals

Goal representation:

```json
{
  "id": 3,
  "name": "Emergency fund",
  "target_amount": 10000,
  "current_amount": 2500,
  "target_date": "2027-12-31"
}
```

`target_date` is optional. `name` must not be blank, `target_amount` must be
positive, `current_amount` must be non-negative, and a supplied target date
must use `YYYY-MM-DD`.

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/goals` | None | `200`, JSON array. | `500` |
| `POST /api/goals` | Goal fields except `id`/timestamps | `201`, created goal. | `400`, `500` |
| `GET /api/goals/{id}` | None | `200`, one goal. | `400`, `404`, `500` |
| `PUT /api/goals/{id}` | Goal object | `200`, updated goal. | `400`, `404`, `500` |
| `DELETE /api/goals/{id}` | None | `204`, no body. | `400`, `404`, `500` |

## Legacy investments and recurring items

### Legacy investments

Investment representation:

```json
{
  "id": 5,
  "amount": "500",
  "investment_type": "stocks",
  "date": "2026-11-08"
}
```

`amount` is encoded as a JSON string by the legacy model.

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/investments` | None | `200`, JSON array. | `500` |
| `POST /api/investments` | `amount`, `investment_type`, `date` | `200`, created investment. | `400`, `500` |
| `GET /api/investments/{id}` | None | `200`, one investment. | `400`, `404` |
| `PUT /api/investments/{id}` | Investment object | `200`, updated investment. | `400`, `404`, `500` |
| `DELETE /api/investments/{id}` | None | `200`, empty body. | `400`, `404`, `500` |

The legacy handler does not validate investment amounts, type, or date beyond
JSON decoding and database behavior. Create and delete return 200.

#### `GET /api/investments/cumulative?lastMonths={count}`

Requires a positive integer `lastMonths`. Returns cumulative amounts by
investment type per month:

```json
[
  {"month": "2026-10", "type": "stocks", "amount": 500},
  {"month": "2026-11", "type": "stocks", "amount": 750}
]
```

Success: `200 OK`; errors: `400` missing/invalid `lastMonths`; `500` otherwise.

### Recurring expenses

Recurring expense representation:

```json
{
  "id": 8,
  "category_id": 4,
  "amount": 42.5,
  "description": "Groceries",
  "start_date": "2026-11-01T00:00:00Z",
  "end_date": null,
  "frequency": "monthly"
}
```

Date/time values are Go `time.Time` JSON strings (RFC 3339); `end_date` may
be null. The handler itself does not validate frequency or amount values.

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/recurring/expenses` | None | `200`, JSON array. | `500` |
| `POST /api/recurring/expenses` | Recurring expense fields | `201`, created item. | `400`, `500` |
| `GET /api/recurring/expenses/{id}` | None | `200`, one item. | `400`, `404`, `500` |
| `PUT /api/recurring/expenses/{id}` | Recurring expense object | `200`, updated item. | `400`, `404`, `500` |
| `DELETE /api/recurring/expenses/{id}` | None | `204`, no body. | `400`, `404`, `500` |

### Recurring investments

Recurring investment representation:

```json
{
  "id": 9,
  "amount": 250,
  "investment_type": "stocks",
  "description": "Monthly contribution",
  "start_date": "2026-11-01T00:00:00Z",
  "end_date": null,
  "frequency": "monthly"
}
```

`amount` is a JSON number for recurring investments (unlike the older
one-time investment API, where it is encoded as a string).

| Method and path | Request | Success response | Error statuses |
| --- | --- | --- | --- |
| `GET /api/recurring/investments` | None | `200`, JSON array. | `500` |
| `POST /api/recurring/investments` | Recurring investment fields | `201`, created item. | `400`, `500` |
| `GET /api/recurring/investments/{id}` | None | `200`, one item. | `400`, `404`, `500` |
| `PUT /api/recurring/investments/{id}` | Recurring investment object | `200`, updated item. | `400`, `404`, `500` |
| `DELETE /api/recurring/investments/{id}` | None | `204`, no body. | `400`, `404`, `500` |

### Recurring job behavior

The server runs recurring expense and investment jobs every 12 hours. The
older recurring APIs store these rules in legacy tables; they are distinct
from `/api/scheduled-flows`, which is the newer mechanism for forecasting and
confirming ledger occurrences.

## Spending trends

### `GET /api/trends/spending`

Optional `months` query parameter defaults to `6` and must be from 1 to 24.
The query starts at the first day of the oldest included month, joins legacy
expenses to categories, and returns monthly category totals:

```json
[
  {"month": "2026-10", "category": "Groceries", "amount": 125.5},
  {"month": "2026-11", "category": "Groceries", "amount": 42.5}
]
```

`amount` is a decimal JSON number. Success: `200 OK`; errors: `400` if
`months` is invalid; `500` if trend data cannot be read or encoded.

## Health and readiness

These routes are outside the `/api` prefix.

### `GET /healthz`

Liveness check. Returns `200 OK` with:

```json
{"status":"ok"}
```

### `GET /readyz`

Checks that the database connection can be obtained and pinged. Returns
`200 OK` with `{"status":"ready"}` when available, or `503 Service
Unavailable` with plain-text body `database unavailable` when not.

## Related documentation

- [UI preview and project overview](../README.md)
- [Backend development notes](../backend/README-backend.md)
