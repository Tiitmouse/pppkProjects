# Patient-Manager

This project is a patient management system with a Go backend and a Vue.js frontend. It provides a RESTful API for managing patient data and a web interface for interacting with the system.

## Features

-   **Patient Management:** Create, read, update, and delete patient records.
-   **Medical Records:** Manage a patient's medical history, including illnesses, checkups, and prescriptions.
-   **User Authentication:** Secure access using JWT (JSON Web Tokens).
-   **API:** A RESTful API for all patient and user operations.
-   **Frontend:** A web-based user interface built with Vue.js.

## Technologies

### Backend (Go)

-   **Web Framework:** [Gin Gonic](https://gin-gonic.com/)
-   **Database:** MongoDB via `go.mongodb.org/mongo-driver`
-   **Dependency Injection:** `go.uber.org/dig`
-   **Logging:** `go.uber.org/zap`
-   **Environment Variables:** `github.com/joho/godotenv` and `github.com/spf13/viper`
-   **Authentication:** JWT via `github.com/dgrijalva/jwt-go` (Note: This package is deprecated; consider a modern alternative like `golang.org/x/crypto/bcrypt` for password hashing and a newer JWT library for tokens).

### Frontend (Vue.js)

-   **Framework:** Vue.js
-   **UI Library:** Vuetify
-   **Package Manager:** npm

## Getting Started

### Prerequisites

-   Go (version 1.24 or newer)
-   Node.js and npm
-   Docker and Docker Compose

### Building and Running

The project can be run using Docker Compose, which will set up the backend, frontend, MongoDB, and Mongo Express for you.

1.  **Configure Environment Variables:**
    Create a `.env` file in the root directory and configure the database connection string and other settings as per `env.example`.

2.  **Start the services:**
    Run the following command in the project root to start all services:
    ```sh
    docker-compose up --build
    ```

3.  **Access the application:**
    -   **Frontend:** The web application will be available at `http://localhost:8080`.
    -   **Backend API:** The API will be available at `http://localhost:8000`.
    -   **Mongo Express:** The MongoDB web interface will be available at `http://localhost:8081`.

### Manual Build

You can use the `Task` runner to build the project binaries.

-   **Build the project:**
    ```sh
    task build
    ```
-   **Run tests:**
    ```sh
    task test
    ```
    (Note: The test task is currently not implemented).