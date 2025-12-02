# Windy Cams Data Fetcher

This Go microservice fetches webcam data from the [Windy Webcams API](https://api.windy.com/webcams/api/v3/webcams) and stores it in a PostgreSQL database. It is designed to be configurable via environment variables and uses GORM for database interactions.

## Prerequisites

Before you begin, ensure you have the following installed:
- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- [PostgreSQL](https://www.postgresql.org/download/)
- A Windy API key, which you can obtain from the [Windy API documentation](https://api.windy.com/).

## Installation

1.  **Clone the repository:**
    ```sh
    git clone <repository-url>
    cd windy-cams
    ```

2.  **Install dependencies:**
    The project uses Go Modules for dependency management. The dependencies will be downloaded automatically when you build or run the project. You can also install them manually:
    ```sh
    go mod tidy
    ```

## Configuration

The application is configured using environment variables. For local development, you can create a `.env` file in the root of the project.

1.  Create a file named `.env` in the project root.

2.  Add the following configuration to the `.env` file, replacing the placeholder values with your own:

    ```env
    # Windy API Configuration
    WINDY_API_KEY=your_secret_key
    API_LIMIT=50
    API_OFFSET=0
    API_SORT_KEY=createdOn
    API_SORT_DIRECTION=desc
    API_CONTINENTS=AF

    # Database Configuration
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    DB_SCHEMA=public
    ```

### Environment Variables

-   `WINDY_API_KEY`: **(Required)** Your secret API key for the Windy API.
-   `API_LIMIT`: The maximum number of webcams to fetch (Default: `50`). Used only by the `fetch` command.
-   `API_OFFSET`: The starting point for the webcam list (Default: `0`). Used only by the `fetch` command.
-   `API_SORT_KEY`: The key to sort the results by (Default: `createdOn`). Used only by the `fetch` command.
-   `API_SORT_DIRECTION`: The direction of the sort (`asc` or `desc`) (Default: `desc`). Used only by the `fetch` command.
-   `API_CONTINENTS`: The continent code to filter webcams by (Default: `AF`). Used only by the `fetch` command.
-   `DB_HOST`: **(Required)** The hostname of your PostgreSQL server.
-   `DB_PORT`: **(Required)** The port for your PostgreSQL server.
-   `DB_USER`: **(Required)** The username for the database.
-   `DB_PASSWORD`: **(Required)** The password for the database user.
-   `DB_NAME`: **(Required)** The name of the database.
-   `DB_SCHEMA`: The database schema to use. (Default: `public`)

## Usage

The application supports three commands, which can be run from the project's root directory: `fetch`, `fetchAll`, and `update`.

### Developer Mode

Both commands accept an optional `--dev` flag that enables developer mode. In this mode, the raw JSON response from the Windy API will be printed to the console, which is useful for debugging.

### `fetch` Command

This command fetches webcams based on the settings in your `.env` file (`API_LIMIT`, `API_OFFSET`, `API_CONTINENTS`, etc.).

```sh
# Basic usage
go run cmd/*.go fetch

# With developer mode
go run cmd/*.go fetch --dev
```

### `fetchAll` Command

This command fetches all webcams from all continents, ignoring the `API_LIMIT`, `API_OFFSET`, and `API_CONTINENTS` settings in your `.env` file. It iterates through all available continents and paginates through the results to retrieve every webcam.

```sh
# Basic usage
go run cmd/*.go fetchAll

# With developer mode
go run cmd/*.go fetchAll --dev
```

### `update` Command

This command fetches detailed information for every webcam already stored in your database. It iterates through your local records, calls the Windy API for each one, and updates the record with the new, detailed data.

```sh
# Basic usage
go run cmd/*.go update

# With developer mode
go run cmd/*.go update --dev
```

When run, the application will start, connect to the database, fetch data from the Windy API according to the specified command, and save any new webcams it finds.

## Project Structure

The project follows a standard Go application layout:

-   `cmd/main.go`: The main entry point for the application.
-   `pkg/`: Contains the core application logic, separated into packages:
    -   `config/`: Handles loading configuration from environment variables.
    -   `storage/`: Manages database connections, models, and migrations.
    -   `windy/`: Contains the client for interacting with the Windy Webcams API.
-   `go.mod`: Defines the project's module and dependencies.
-   `.gitignore`: Specifies files and directories to be ignored by version control.