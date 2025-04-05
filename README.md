# SSD Assignment API

## Description
A Go-based API for managing configurations with authentication using JWT tokens. It includes routes for registering, logging in, and managing configurations.

## Setup

### Prerequisites
- Go 1.16 or higher
- Git
- A terminal (Linux, macOS, or Windows with WSL)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/merveinan/ssd-assignment-api.git
   cd ssd-assignment-api

2. Install Dependencies
    ```bash
    go mod tidy

3. Run The Application
    ```bash
    go run main.go

### Swagger Documentation
    http://localhost:8000/swagger/index.html
    
### Deployment

1. Build the Go application:
   ```bash
        go build

2. Deploy using Docker (optional):

    Create a Dockerfile.

3. Build the Docker image:

    ```bash
        docker build -t ssd-assignment-api .

4. Run the container:

    ```bash
        docker run -p 8000:8000 ssd-assignment-api

5. Alternatively, if you're using Docker Compose, you can run:
    
    ```bash
        docker-compose up
    