# Gator - a RSS blog aggregator CLI
To use this program you'll need to have Postgres and Go installed.

Create a `~/.gatorconfig.json` containing `{"db_url": "postgres://your-information", "current_user_name": "your-name"}`

Then run `go install github.com/nhmosko/gator@latest`.

## How to use:
`gator help` will give you a full list of commands.

`gator register` and `gator login` take a username as an argument and create/select a user.

`gator feeds` and `gator users` display all the feeds/users from the database.

Once you are logged in you may use:

`gator addfeed <feed-name> <feed-url>` to add a new feed to the database.

`gator follow <feed-url>` to follow an existing feed.

`gator unfollow <feed-url>` to unfollow a followed feed.

`gator following` to see your followed feeds.

To actually see the posts:

run `gator agg <interval>` for a while to fetch the posts (I usually use 1s as the interval)

Finally run `gator browse [optional <limit>]` to see the posts from your followed feeds.

!!! `gator reset` deletes everything from the database (careful).
