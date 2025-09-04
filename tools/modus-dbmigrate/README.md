# DB Migrate Tool

This is a small utility that will apply database migrations for the Modus runtime's PostgreSQL DB.

## Usage

You must have Go v1.25 or newer installed locally. We currently do not release this tool as a
binary.

Before running, set the `MODUS_DB` environment variable to the connection string URI for the
PostgreSQL database you want to target.

You can run the tool in one of two ways:

- Clone this repo, then from the `/tools/modus-dbmigrate` directory, run the utility via: `go run .`
- Install it with `go install github.com/hypermodeinc/modus/tools/modus-dbmigrate@latest` then run
  `modus-dbmigrate`

The tool will connect to the database and apply any migrations required to bring the schema up to
date.

## Note

The migration tool embeds the appropriate migration scripts. You do not need to copy them separately
or install any other tool.
