# Split Receipt

A small Go REST API for splitting restaurant receipts among friends. It provides simple endpoints to model receipts, meals, friends, and fees, and to calculate each friend's share.

## **Motivation**

- Split shared bills accurately and transparently.
- Support fixed and percentage-based fees (tax, service) and per-meal assignments.
- Provide a minimal, well-tested backend that can be used as a reference or integrated into a larger app.

## **Quick Start**

Prerequisites:

- Go 1.24+
- PostgreSQL
- `goose` (migrations)
- `sqlc` (code generation)

Run the following to create the DB, migrate, build and start the server locally:

```bash
# create database
createdb split_receipt

# run migrations
goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/split_receipt?sslmode=disable" up

# build and run
go build -o server ./cmd/server
DATABASE_URL="postgres://postgres:postgres@localhost:5432/split_receipt?sslmode=disable" ./server
```

The server listens on port `8080` by default. Set the `PORT` environment variable to change it.

## **Usage**

Primary concepts:

- Friend — a person who can be assigned to meals.
- Receipt — a bill containing meals and fees.
- Meal — an item on a receipt with a price and assigned friends.
- Fixed fee — amount shared evenly across participating friends (e.g., delivery).
- Percentage fee — percentage applied to each friend's meal subtotal (e.g., tax), optionally capped.

Common endpoints (HTTP):

- `POST /friends` — create a friend
- `GET /friends` — list friends
- `POST /receipts` — create a receipt
- `GET /receipts/{id}` — get receipt details
- `POST /receipts/{id}/meals` — add a meal to a receipt
- `POST /receipts/{id}/meals/{mealId}/friends` — assign a friend to a meal
- `POST /receipts/{id}/fixed-fees` — add a fixed fee
- `POST /receipts/{id}/percentage-fees` — add a percentage fee
- `GET /receipts/{id}/friends/{friendId}/settlement` — calculate a friend's total share

Example: calculate settlement for friend `f1` on receipt `r1`:

```bash
curl -s "http://localhost:8080/receipts/r1/friends/f1/settlement" | jq
```

Settlement calculation summary:

1. For each meal the friend participates in: $meal_share = meal.total_price / count(friends\_in\_meal)$
2. Sum meal shares → `friend_meal_total`
3. Fixed fee share: `fee.amount / count(unique_friends_in_receipt)` for each fixed fee
4. Percentage fee share: `fee.percentage * friend_meal_total`, clamped to `cap_amount` if provided
5. Total = `friend_meal_total + sum(fixed_fee_shares) + sum(percentage_fee_shares)`

For full API details and request/response shapes, see the OpenAPI spec: [api/openapi.yaml](api/openapi.yaml)

## **Development**

Regenerate sqlc code after changing queries:

```bash
sqlc generate
```

Run migrations:

```bash
goose -dir sql/migrations postgres "$DATABASE_URL" up
```

Build all packages:

```bash
go build ./...
```

## **Contributing**

- Fork the repo and open a PR with small, focused changes.
- Run `go vet ./...` and `go test ./...` before submitting.
- Keep database schema changes in `sql/migrations/` and regenerate `sql/queries/` and `internal/database` artifacts as needed (`sqlc generate`).
- Document API changes in `api/openapi.yaml`.

If you'd like help adding features or tests, open an issue describing the proposed change.
