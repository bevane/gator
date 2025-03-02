# Gator

A CLI tool for aggregating RSS feeds

## Boot.dev project

This project was completed as part of the [boot.dev](https://www.boot.dev) back-end developer course curriculum

## Learning Concepts applied

- Integrating a Go app with Postgres database
- Migrating databases with incremental database schema changes
- Using type-safe SQL to store and retrieve data from a Postgres Database
- Writing middleware for user authorization

# Installation and Setup

Requirements:
- Postgres
- Go 1.22+

1. `go install github.com/bevane/gator` to install the CLI tool
2. Create a new postgres db 
3. Create `.gatorconfig.json` in your home directory in the following format:
```json
{
    "db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
    "current_user_name":"bevane"
}
```
note: replace connection string with your own conenction string to the db you created in step 2 and set your own username

# Usage

`gator login [username]` to login as a specific user

`gator register [username]` to register a username

`gator users` to see list of user and current user

`gator addfeed [feed name] [url]` to add and follow a new feed

`gator follow [url]` to follow a feed

`gator unfollow [url]` to unfullow a feed

`gator following` to see a list of feeds the user is following

`gator agg [time_between_requests]` to collect feeds every [time_between_requests]
example: `gator agg 10s`, `gator agg 10m`, `gator agg 1h30m`

`gator browse [number_of_posts]` to see last n recent posts from the aggregated feeds



