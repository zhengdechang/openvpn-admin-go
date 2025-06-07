# openvpn-web

This project is the frontend for the `openvpn-admin-go` backend. It is built with Next.js, TypeScript, and Tailwind CSS, providing a modern and user-friendly web interface to manage OpenVPN servers, users, and configurations.

## Features

- **Secure User Login:** Robust authentication and session management to protect access.
- **Dashboard Overview:** A comprehensive dashboard displaying the status of OpenVPN servers and key metrics.
- **Server Configuration Management:** Interface for viewing, and potentially creating, editing, and deleting OpenVPN server configurations (depending on backend capabilities).
- **User Management:** Tools for administrators to manage user accounts, assign permissions, and control VPN access.
- **Department Management:** Organize users into departments for streamlined administration.
- **Client Configuration Download:** Easily generate and download OpenVPN client configuration files (e.g., `.ovpn`) for users.
- **Log Viewing:** Access and view server logs and client connection logs for monitoring and troubleshooting.
- **User Profile Management:** Allows users to view and manage their profile settings.

## Prerequisites

- Node.js 18.x or higher (as specified in `package.json` or project documentation).
- `npm` (Node Package Manager) or `yarn`.
- A running instance of the `openvpn-admin-go` backend.

## Installation

1.  **Clone the Repository:**
    ```bash
    git clone https://github.com/your-username/openvpn-web.git
    cd openvpn-web
    ```
    (Replace `https://github.com/your-username/openvpn-web.git` with the actual repository URL)

2.  **Install Dependencies:**
    Using npm:
    ```bash
    npm install
    ```
    Or using yarn:
    ```bash
    yarn install
    ```

3.  **Configure Environment Variables:**
    *   Copy the example environment file:
        ```bash
        cp .env.example .env
        ```
    *   Edit the `.env` file and set the following variables:
        *   `NEXT_PUBLIC_API_BASE_URL`: This is crucial. Set it to the URL where your `openvpn-admin-go` backend is running (e.g., `http://localhost:8080/api`).

4.  **Run the Development Server:**
    Using npm:
    ```bash
    npm run dev
    ```
    Or using yarn:
    ```bash
    yarn dev
    ```
    The application should now be accessible at `http://localhost:3000` (or another port if specified).

## Project Structure

Here's an overview of the key directories and files in the project:

-   `src/app/`: Core of the Next.js application using the App Router. Contains pages, layouts, and route definitions.
    -   `dashboard/`: Contains the different pages and layouts for the main application dashboard after login.
-   `src/components/`: Reusable React components.
    -   `ui/`: Basic UI elements like buttons, inputs, cards, etc. (often from a UI library like Shadcn/ui).
    -   `layout/`: Components responsible for structuring the overall layout of the application (e.g., header, sidebar, footer).
-   `src/services/`: Handles API communication with the `openvpn-admin-go` backend. Contains functions for fetching and sending data.
    -   `api.ts`: Likely the main file for configuring the API client (e.g., Axios instance) and defining API call functions.
-   `src/lib/`: Utility functions, helper scripts, and shared modules.
    -   `auth-context.tsx`: React Context for managing authentication state (e.g., user session, tokens).
    -   `auth-utils.ts`: Utility functions related to authentication (e.g., token storage, validation).
    -   `utils.ts`: General utility functions.
-   `src/store/`: Client-side state management.
    -   `useUserStore.ts`: Zustand store for managing global user-related state.
-   `src/i18n/`: Internationalization and localization setup. Contains configuration and translation files for supporting multiple languages.
-   `src/types/`: TypeScript type definitions and interfaces used throughout the project.
-   `public/`: Static assets that are served directly by the web server (e.g., images, favicons). Assets in this directory are not processed by the build pipeline.
-   `src/assets/`: Static assets like images and fonts that are imported into your components. These assets are processed by Next.js's build system and can be optimized or bundled.
-   `tailwind.config.js`: Configuration file for Tailwind CSS, allowing customization of design tokens, plugins, etc.
-   `next.config.js`: Configuration file for Next.js, used to customize its behavior (e.g., redirects, environment variables, webpack modifications).
-   `.env.example`: Example environment variable file. Copy this to `.env` to configure your local setup.
-   `.env`: Local environment variable file (should not be committed to Git). Contains sensitive or environment-specific settings like API keys or backend URLs.
-   `package.json`: Lists project dependencies, scripts (for building, developing, linting), and project metadata.
-   `yarn.lock` (or `package-lock.json`): Locks down the exact versions of dependencies.
-   `tsconfig.json`: TypeScript compiler configuration.

## Environment Variables

Environment variables are used to configure the application without hardcoding values into the source code. This is particularly important for sensitive information or settings that vary between environments (development, production, etc.).

-   **`.env.example`**: This file serves as a template for the actual environment variables file. It lists all the environment variables the application expects and provides example values. **Do not store sensitive information in this file.**
-   **`.env`**: This is the actual file where you will define your environment-specific variables.
    -   Create it by copying `.env.example`: `cp .env.example .env`
    -   This file is listed in `.gitignore` and **should not be committed to your version control system (Git)** to protect sensitive data.

### Key Environment Variables

-   `NEXT_PUBLIC_API_BASE_URL`:
    -   **Purpose**: Specifies the base URL for the `openvpn-admin-go` backend API. The frontend will make API requests to this URL.
    -   **Example**: `NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api` or `NEXT_PUBLIC_API_BASE_URL=https://your-backend-domain.com/api`
    -   **Note**: The `NEXT_PUBLIC_` prefix makes this variable accessible in the browser-side JavaScript code.

Make sure to set all required variables in your `.env` file after copying it from `.env.example` for the application to function correctly.
