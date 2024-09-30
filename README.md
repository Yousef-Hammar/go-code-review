# Schwarz IT Code Review Repository

This repository contains my solution for the coding challenge completed during the onboarding process as a Go
developer for Schwarz IT. It aims to address various code smells and security issues identified in the original code.

## Prerequisites

- [Go (1.22.7 or higher)](https://go.dev)
- [Docker](https://www.docker.com)
- [Testify (1.90.0)](https://github.com/vektra/mockery)
- [Ginkgo(2.20.2)](https://github.com/onsi/ginkgo)
- [Mockery (2.40.2)](https://github.com/vektra/mockery])

## Repository Structure

The repository is organized into the following subfolders:

- **`cmd`**: Contains the `main.go` file that initializes and starts the server.
- **`internal`**: Contains the core server code, including routing and middleware.
- **`tests`**: Contains integration test files for validating the application functionality.

## Backend Architecture

The backend is structured into three main layers:

1. **API Layer** (`api`):
    - Responsible for starting the server and handling incoming requests.

2. **Service Layer** (`service`):
    - Contains the business logic of the application, ensuring separation of concerns.

3. **Repository Layer** (`repository`):
    - Manages data storage. This implementation currently uses an in-memory storage solution.

## How to Run

To run the application using the Makefile, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/Yousef-Hammar/go-code-review
   cd go-code-review
   ```
2. Create a .env file containing the following line:
   ```
   ADDR=8080
   ```
3. Run the application with Docker:
   ```
   make docker-run
   ```

This command will build the Docker image and run the server, using the environment variables specified in the .env file.

## How to Test

To run tests using the Makefile, you have two options:

**1. Run Unit Tests:**
   ```bash
   make run-unit-tests
   ```
**2. Run Integration Tests:**
   ```bash
   make run-integration-tests
   ```

These commands will execute all unit tests and integration tests, respectively, providing detailed output.

