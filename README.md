# gator feed aggregator

A simple, cli, feed aggregator in Go following the boot.dev gator spec.

## Requirements
You will need to have these programs already installed and working on your system before installing gator:

- Go
- Postgres

## Install

1. Create the database in postgres.
    - Login to postgres and create a new database
      ```sql
      CREATE DATABASE gator;
      ```
    - Create a database user for the gator database. This is used by the gator
    client to connect to the database. You only need one database account as all
    gator users will share this account. You can use your admin database account if
    you want to, but it is best practice to create a dedicated user for each application.
    You can name the account as you please; the example uses gator for the username.
      ```sql
      CREATE USER gator WITH PASSWORD 'YOUR_secure_p@ssw0rd?';
      ```
    - Give the gator user (substitute your username if different) to the gator database (substitute your database name if different.)
    ```sql
    GRANT ALL PRIVILEGES ON DATABASE gator TO gator;
    ```
2. Go install the program `go install github.com/theMagicRabbit/gator` This command will download, compile, and install the gatorcli.
**Configure gator before trying to run gator.**

## Configuration

gator reads from a JSON configuration file at `~/.gatorconfig.json`. You will need to create this file before runing gator.
Below is a sample of what the file will look like:

```json
{
  "db_url": "postgres://YOUR_secure_p@ssw0rd?:gator@localhost:5432/gator?sslmode=disable"
}
```

In the db_url, the format is `postgres://<password>:<username>@<hostname>:<db_port>/<datbase_name>?sslmode=disable`
Replace the placeholders with your username, password, and database name as needed. `localhost:5432` should work for your
hostname and database port, unless you have changed the default setup. If this does not work for you, you will need
to check your postgres configuration and see what port and host the database service is listening on.

After you have run gator and created a gator user, gator will add your username to the config file. You do not need
to manually add a user to the config file, this is only a note in case you happen to notice that there is a username
in the config file. That is normal and no need for concern.

### Database schema
When ever you run the gator cli, it will use an embeded [goose](https://github.com/pressly/goose) migration to create
the database schema. This is so you do not have to manually create the database schema. Yay for automation.

## Usage

All examples will use the username "brt". Your username can be whatever you wish it to be, simply replace "brt" with your username.

### Create a new user:

`gator register brt`

### Login the user:

`gator login brt`

### Create a new feed in the system:

For this step, you will need to know the feed you wish to follow.

`gator addfeed "https://example.com/feed.rss"`

### Follow a new feed:

For this step, you will need to know the feed you wish to follow.

`gator follow "https://example.com/feed.rss"`

### List subscribed feeds for your user

`gator following`

### List all feeds in the datbase:


```
gator feeds
```

### Delete *all* data in the database

This will erase *everything* in the database.

```
gator reset
```

### List all users:

```
gator users
```

### Unsubscribe from a feed:

```
gator unfollow "https://example.com/feed.rss"
```

### List recent posts

```
gator browse [count]
```

Count is the number of posts to display; value is optional and default is 2 posts.

### Check for new posts

This is intended to be run as a background process. You may consider making this a scheduled task.
It only updates one feed each run. The time between runs is set by the interval. At the interval,
the most stale feed is checked for new posts. The more feeds you have, the shorter the interval
should be.

The interval argument is requred and should be given something like `1h2m3s` which would check for new
posts once every 1 hour, 2 minutes, and 3 seconds. A more likely setting is something like 30 minutes or
`30m`. 

```
gator agg 30m
```


