# openvpn-admin-go

This project is the backend for a web-based OpenVPN management panel. It is built with Go and provides a RESTful API for the frontend (openvpn-web) to consume. The primary goal of this project is to simplify the management of OpenVPN servers, users, and their configurations through a user-friendly web interface.

## Features

- **User Authentication and Management:** Secure user registration and login for administrators and regular users.
- **Department Management:** Organize users into departments for better access control and management.
- **OpenVPN Server Configuration Management:** Easily create, update, and delete OpenVPN server configurations.
- **OpenVPN Client Configuration Generation and Management:** Generate client configuration files (e.g., .ovpn) and manage client access.
- **Server Status Monitoring:** Monitor the status of OpenVPN servers, including active connections and traffic.
- **Client Connection Monitoring:** View currently connected clients and their session details.
- **Log Viewing:** Access and view server and client logs for troubleshooting and monitoring.

## Installation and Setup

1.  **Prerequisites:**
    *   Go 1.18 or higher.

2.  **Clone the Repository:**
    ```bash
    git clone https://github.com/your-username/openvpn-admin-go.git
    # (Replace with the actual repository URL if different)
    cd openvpn-admin-go
    ```

3.  **Install Dependencies:**
    ```bash
    go mod tidy
    ```

4.  **Configure the Application:**
    *   Copy the example environment file:
        ```bash
        cp .env.example .env
        ```
    *   Edit the `.env` file and set the necessary variables, including:
        *   `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` for database connection.
        *   `JWT_SECRET` for signing JWT tokens.
        *   Other relevant settings for your environment.

5.  **Run the Application:**
    ```bash
    go run main.go
    ```
    Alternatively, if your main package is in a `cmd` directory:
    ```bash
    go run cmd/main.go
    ```

## Project Structure

Here's an overview of the key directories and their purpose:

*   `cmd/`: Contains the main application entry point(s).
*   `common/`: Utility functions and helpers shared across the project.
*   `constants/`: Project-wide constants, such as configuration keys or default values.
*   `controller/`: Handles incoming HTTP requests, processes them, and interacts with services.
*   `database/`: Database connection management, migrations, and query builders.
*   `middleware/`: HTTP middleware for tasks like authentication, logging, and CORS.
*   `model/`: Data structures and types representing application entities (e.g., User, OpenVPN Server).
*   `openvpn/`: Core logic for interacting with OpenVPN servers (e.g., managing configurations, monitoring status).
*   `router/`: Defines API routes and maps them to controller handlers.
*   `services/`: Business logic and coordination between controllers and data layers (database, OpenVPN).
*   `template/`: HTML templates or other template files, if any (e.g., for client configuration files).

## API Endpoints

The RESTful API allows the frontend (openvpn-web) and other clients to interact with the backend. The API routes are defined in the `router/` directory. Endpoints are generally grouped by functionality:

*   **Authentication:** Endpoints for user login, registration, and token management.
*   **Users:** Managing user accounts and profiles.
*   **Departments:** Managing departments for user organization.
*   **Servers:** CRUD operations for OpenVPN server configurations.
*   **Clients:** Generating and managing client VPN configurations.
*   **Status & Monitoring:** Endpoints for checking server status and connected clients.
*   **Logs:** Accessing server and client logs.

For detailed information on specific endpoints, please refer to the code within the `router/` directory and the corresponding controller handlers.
