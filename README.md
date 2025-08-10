# E-commerce Backend Microservices

This project is a sample e-commerce backend built with a microservices architecture using Go. It includes services for users, products, orders, payments, inventory, and notifications.

## Project Structure

```
ecommerce-platform/
├── services/
│   ├── user-service/
│   ├── product-service/
│   ├── order-service/
│   ├── payment-service/
│   ├── inventory-service/
│   └── notification-service/
├── shared/
│   ├── pkg/
│   ├── proto/
│   └── configs/
├── infrastructure/
│   ├── terraform/
│   └── helm/
├── monitoring/
└── docker-compose.yml
```

## Prerequisites

- [Go](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Task](https://taskfile.dev/)

## Setup

1.  **Clone the repository:**

    ```bash
    git clone git@github.com:hwgao/ecommerce-demo.git
    cd ecommerce-demo
    ```

2.  **Install dependencies:**

    ```bash
    task deps
    ```

## Available Commands

-   `task build`: Build all services.
-   `task run`: Run all services using Docker Compose.
-   `task test`: Run all tests.
-   `task deps`: Install dependencies.
-   `task lint`: Run linter.
-   `task clean`: Remove build artifacts.
-   `task docker-build`: Build all Docker images.

## Running the Services

To run all the services, use the following command:

```bash
task run
```

This will start all the services and their dependencies using Docker Compose.

## Microservices

-   **User Service**: Manages user accounts, authentication, and profiles.
-   **Product Service**: Manages products, categories, and product search.
-   **Order Service**: Manages orders, order status, and order history.
-   **Payment Service**: Manages payments, refunds, and payment provider integrations.
-   **Inventory Service**: Manages product inventory and stock levels.
-   **Notification Service**: Manages sending notifications to users (e.g., email, SMS).
