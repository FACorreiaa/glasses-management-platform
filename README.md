# Glasses Management Platform

This project is a web-based platform for managing glasses, including inventory and shipping. It's built with a Go backend, a PostgreSQL database, and a modern frontend stack using HTMX, Tailwind CSS, and DaisyUI.

## Technologies Used

*   **Backend:**
    *   [Go](https://golang.org/)
    *   [Gorilla/Mux](https://github.com/gorilla/mux) for routing
    *   [pgx](https://github.com/jackc/pgx) for PostgreSQL driver and toolkit
    *   [Templ](https://templ.dev/) for Go templating
    *   [go-playground/validator](https://github.com/go-playground/validator) for data validation
    *   [godotenv](https://github.com/joho/godotenv) for managing environment variables
    *   [OpenTelemetry](https://opentelemetry.io/) for observability
    *   [Prometheus](https://prometheus.io/) for monitoring

*   **Frontend:**
    *   [HTMX](https://htmx.org/)
    *   [Tailwind CSS](https://tailwindcss.com/)
    *   [DaisyUI](https://daisyui.com/)
    *   [PostCSS](https://postcss.org/)
    *   [esbuild](https://esbuild.github.io/)

*   **Database:**
    *   [PostgreSQL](https://www.postgresql.org/)

*   **Containerization:**
    *   [Docker](https://www.docker.com/)

## Getting Started

### Prerequisites

*   [Go](https://golang.org/doc/install)
*   [Docker](https://docs.docker.com/get-docker/)
*   [Node.js and npm](https://nodejs.org/en/download/)
*   [Make](https://www.gnu.org/software/make/)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/FACorreiaa/glasses-management-platform.git
    cd glasses-management-platform
    ```

2.  **Set up environment variables:**

    Copy the `.env.sample` file to `.env` and update the variables with your local configuration.

    ```bash
    cp .env.sample .env
    ```

3.  **Install frontend dependencies:**

    ```bash
    npm install
    ```

4.  **Build and run the application with Docker:**

    ```bash
    make compose-dev
    ```

    This command will start the application in development mode with live reloading.

## Available Commands

This project uses a `Makefile` to streamline common tasks. Here are some of the most useful commands:

*   `make compose-dev`: Start the application in development mode.
*   `make compose-prod`: Start the application in production mode.
*   `make compose-down`: Stop the application and remove the containers.
*   `make go-test`: Run the Go tests.
*   `make lint`: Lint the Go code.
*   `make build-tailwind`: Build the Tailwind CSS for production.
*   `make watch-tailwind`: Watch for changes in Tailwind CSS files and rebuild automatically.

## Database Migrations

The database migrations are located in the `db/migrations` directory. To apply the migrations, you can use the `migrate` tool.

## Contributing

Contributions are welcome! Please feel free to open an issue or submit a pull request.

## License

This project is licensed under the ISC License. See the `LICENSE` file for more details.