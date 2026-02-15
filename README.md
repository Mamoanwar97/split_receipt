# Split Receipt

A Go REST API server for splitting restaurant receipts among friends.

## Prerequisites

- Go 1.24+
- PostgreSQL
- [goose](https://github.com/pressly/goose) (migrations)
- [sqlc](https://sqlc.dev) (code generation)

## Setup

```bash
# Create database
createdb split_receipt

# Run migrations
goose -dir sql/migrations postgres "postgres://postgres:postgres@localhost:5432/split_receipt?sslmode=disable" up

# Build and run
go build -o server ./cmd/server
DATABASE_URL="postgres://postgres:postgres@localhost:5432/split_receipt?sslmode=disable" ./server
```

The server starts on port 8080 by default. Set the `PORT` environment variable to change it.

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /friends | Create a friend |
| GET | /friends | List all friends |
| GET | /friends/{id} | Get a friend |
| DELETE | /friends/{id} | Delete a friend |
| POST | /receipts | Create a receipt |
| GET | /receipts | List all receipts |
| GET | /receipts/{id} | Get a receipt |
| DELETE | /receipts/{id} | Delete a receipt |
| POST | /receipts/{id}/meals | Add a meal to a receipt |
| GET | /receipts/{id}/meals | List meals on a receipt |
| DELETE | /receipts/{id}/meals/{mealId} | Remove a meal |
| POST | /receipts/{id}/meals/{mealId}/friends | Assign a friend to a meal |
| DELETE | /receipts/{id}/meals/{mealId}/friends/{friendId} | Remove friend from meal |
| POST | /receipts/{id}/fixed-fees | Add a fixed fee |
| GET | /receipts/{id}/fixed-fees | List fixed fees |
| DELETE | /receipts/{id}/fixed-fees/{feeId} | Remove a fixed fee |
| POST | /receipts/{id}/percentage-fees | Add a percentage fee |
| GET | /receipts/{id}/percentage-fees | List percentage fees |
| DELETE | /receipts/{id}/percentage-fees/{feeId} | Remove a percentage fee |
| GET | /receipts/{id}/friends/{friendId}/settlement | Calculate friend's share |

## Settlement Calculation

1. For each meal the friend participates in: `meal.total_price / count(friends_in_meal)`
2. Sum all meal shares = `friend_meal_total`
3. Each fixed fee: `fee.amount / count(unique_friends_in_receipt)`
4. Each percentage fee: `fee.percentage * friend_meal_total` (clamped to `cap_amount` if set)
5. Total = friend_meal_total + sum(fixed_fee_shares) + sum(percentage_fee_shares)

## Development

```bash
# Regenerate sqlc code after changing queries
sqlc generate

# Run migrations
goose -dir sql/migrations postgres "$DATABASE_URL" up

# Build
go build ./...
```
