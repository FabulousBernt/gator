# Gator

A multi-user CLI RSS feed aggregator built with Go and PostgreSQL. Add feeds, follow them, and browse the latest posts right in your terminal.

## Prerequisites

- [Go](https://golang.org/dl/) (1.21 or later)
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

Install the CLI using `go install`:
```bash
go install github.com/FabulousBernt/gator@latest
```

## Configuration

Create a config file at `~/.gatorconfig.json`:
```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable"
}
```

Replace `username` and `password` with your PostgreSQL credentials.

## Database Setup

Install [Goose](https://github.com/pressly/goose) and run the migrations:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir sql/schema postgres "postgres://username:@localhost:5432/gator?sslmode=disable" up
```

## Commands

| Command | Description |
|---|---|
| `gator register <name>` | Register a new user |
| `gator login <name>` | Login as an existing user |
| `gator users` | List all users |
| `gator addfeed <name> <url>` | Add a new feed and follow it |
| `gator feeds` | List all feeds |
| `gator follow <url>` | Follow an existing feed |
| `gator unfollow <url>` | Unfollow a feed |
| `gator following` | List feeds you are following |
| `gator agg <interval>` | Start the aggregator (e.g. `gator agg 1m`) |
| `gator browse [limit]` | Browse latest posts (default limit: 2) |
| `gator reset` | Reset the database |

## Usage Example
```bash
# Register and login
gator register FabulousBernt

# Add some feeds
gator addfeed "Boot.dev Blog" "https://www.wagslane.dev/index.xml"
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"

# Start the aggregator in one terminal
gator agg 1m

# Browse posts in another terminal
gator browse 10
```