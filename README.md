# Windy Webcams Platform Importer

This project is a Go application that fetches asset platform data from the [Windy Webcams API](https://api.windy.com/webcams/docs) and stores it in a PostgreSQL database. The application is structured following Clean Architecture principles to ensure it is modular, maintainable, and testable.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation and Setup](#installation-and-setup)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

- Fetches platform data from the Windy API.
- Stores data in a PostgreSQL database using GORM.
- Uses GORM's auto-migration to create the database schema.
- Configuration is managed via environment variables.

## Architecture

This project follows **Clean Architecture** principles to separate concerns and create a maintainable and scalable system. The main layers are:

- **Domain**: Contains the core business logic and entities (e.g., `Platform`).
- **Application/Service**: Orchestrates the data flow and implements use cases (e.g., `FetchAndStorePlatforms`).
- **Infrastructure/Platform**: Handles external concerns like database access (`repository`), API clients (`windy`), and configuration.

Dependencies point inwards, from the outer layers (infrastructure) to the inner layers (domain), ensuring that the core business logic is independent of external frameworks and tools.

## Getting Started

Follow these instructions to set up and run the project locally.

### Prerequisites

Ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (Version 1.18 or later)
- [PostgreSQL](https://www.postgresql.org/download/)

### Installation and Setup

1.  **Clone the Repository**

    ```sh
    git clone <your-repository-url>
    cd <repository-directory>
    ```

2.  **Set Up the Database**

    Connect to your PostgreSQL instance and create a new database.

    ```sql
    CREATE DATABASE webcams;
    ```

3.  **Configure Environment Variables**

    The application uses a `.env` file for configuration. Copy the example file and update it with your credentials.

    ```sh
    cp .env.example .env
    ```

    Edit the `.env` file with your Windy API key and PostgreSQL connection details:

    ```env
    # Windy API Key
    WINDY_API_KEY=your_actual_api_key_here

    # PostgreSQL connection details
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_postgres_user
    DB_PASSWORD=your_postgres_password
    DB_NAME=webcams
    DB_SSLMODE=disable
    ```

4.  **Install Dependencies**

    Install the required Go modules:

    ```sh
    go mod tidy
    ```

## Usage

To run the application, execute the `main.go` file. This will trigger a one-time operation to connect to the database, run migrations, fetch platform data from the Windy API, and store it.

```sh
go run cmd/app/main.go
```

You should see output in your console indicating the progress:

```
Database migration completed.
Fetching and storing platforms...
Successfully saved <N> platforms.
Process finished successfully.
```

## Project Structure

The codebase is organized to separate concerns, making it easier to maintain and test.

-   `cmd/app`: Contains the main application entrypoint.
-   `configs`: Handles loading configuration from environment variables.
-   `internal/`: Contains the core application logic.
    -   `domain`: Defines the primary data structures (e.g., Platform).
    -   `platform`: Implements the business logic (service) and data access (repository) for platforms.
-   `pkg/`: Contains shared, reusable packages.
    -   `database`: Manages the PostgreSQL database connection and migrations.
    -   `windy`: Provides a client for interacting with the Windy Webcams API.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
