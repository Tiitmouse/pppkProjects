# Data-Management

This project is a general terminal application template written in Go. It includes several subcommands for data management tasks.

## Features

- **`version`**: Displays the application version and build information.
- **`tsv`**: Parses TSV files.
- **`show`**: Shows specific data.
- **`scrape`**: Scrapes data from `https://xenabrowser.net/datapages/` using a Python script with Selenium.

## Dependencies

The project uses the following Go modules:

- `go.uber.org/zap` for logging
- `fyne.io/fyne/v2` for a GUI (indirectly, based on go.mod)
- `github.com/google/subcommands` for managing subcommands
- `github.com/joho/godotenv` for handling environment variables
- `github.com/minio/minio-go/v7` for interacting with MinIO
- `go.mongodb.org/mongo-driver` for MongoDB database access

## Building and Installation

This project uses `Task` for automation. You can use the following commands:

- `task dev`: Builds a binary for development.
- `task build`: Builds the project binary with version, commit hash, and timestamp information.
- `task install`: Installs the project as a local binary.
- `task test`: Runs project tests (currently not implemented).
- `task coverage`: Creates a project coverage file (currently not implemented).

## Usage

Run the compiled binary with one of the available subcommands:

```
./data-managment <command>
```

For more information on a specific command, use the help flag:

```
./data-managment help <command>
```