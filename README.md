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
-   `API_LIMIT`: The maximum number of webcams to fetch. (Default: `50`)
-   `API_OFFSET`: The starting point for the webcam list. (Default: `0`)
-   `API_SORT_KEY`: The key to sort the results by. (Default: `createdOn`)
-   `API_SORT_DIRECTION`: The direction of the sort (`asc` or `desc`). (Default: `desc`)
-   `API_CONTINENTS`: The continent code to filter webcams by. (Default: `AF`)
-   `DB_HOST`: **(Required)** The hostname of your PostgreSQL server.
-   `DB_PORT`: **(Required)** The port for your PostgreSQL server.
-   `DB_USER`: **(Required)** The username for the database.
-   `DB_PASSWORD`: **(Required)** The password for the database user.
-   `DB_NAME`: **(Required)** The name of the database.
-   `DB_SCHEMA`: The database schema to use. (Default: `public`)

## Usage

To run the microservice, execute the following command from the project's root directory:

```sh
go run cmd/main.go
```

The application will start, connect to the database, fetch data from the Windy API, and save any new webcams it finds.

## Project Structure

The project follows a standard Go application layout:

-   `cmd/main.go`: The main entry point for the application.
-   `pkg/`: Contains the core application logic, separated into packages:
    -   `config/`: Handles loading configuration from environment variables.
    -   `storage/`: Manages database connections, models, and migrations.
    -   `windy/`: Contains the client for interacting with the Windy Webcams API.
-   `go.mod`: Defines the project's module and dependencies.
-   `.gitignore`: Specifies files and directories to be ignored by version control.