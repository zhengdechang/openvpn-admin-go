# OpenVPN Web Admin Interface

[![GitHub Repository](https://img.shields.io/badge/GitHub-openvpn--admin--go-blue?style=for-the-badge&logo=github)](https://github.com/zhengdechang/openvpn-admin-go)
[![Next.js](https://img.shields.io/badge/Next.js-14-black?style=for-the-badge&logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.3-blue?style=for-the-badge&logo=typescript)](https://www.typescriptlang.org/)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.3-38B2AC?style=for-the-badge&logo=tailwind-css)](https://tailwindcss.com/)

This project is the modern, responsive frontend interface for the `openvpn-admin-go` backend system. Built with cutting-edge web technologies including Next.js 14, TypeScript, and Tailwind CSS, it provides an intuitive and powerful web interface for comprehensive OpenVPN server management, user administration, and real-time monitoring.

## 🌟 Project Overview

OpenVPN Web Admin Interface is designed to simplify OpenVPN management through a beautiful, user-friendly web dashboard. Whether you're managing a small team or enterprise-level VPN infrastructure, this interface provides all the tools you need in one centralized location.

**🔗 Main Repository:** [openvpn-admin-go](https://github.com/zhengdechang/openvpn-admin-go)

## 🚀 Technology Stack

### Frontend Technologies
- **⚡ Next.js 14** - React framework with App Router for modern web applications
- **📘 TypeScript 5.3** - Type-safe JavaScript for better development experience
- **🎨 Tailwind CSS 3.3** - Utility-first CSS framework for rapid UI development
- **🧩 Radix UI** - Unstyled, accessible UI components
- **🌐 i18next** - Internationalization framework (English/Chinese support)
- **📊 Zustand** - Lightweight state management
- **🔗 Axios** - HTTP client for API communication

### Backend Integration
- **🔧 Go Backend** - High-performance backend built with Go
- **🌐 Gin Framework** - Fast HTTP web framework
- **💾 SQLite Database** - Lightweight, embedded database
- **🔐 JWT Authentication** - Secure token-based authentication
- **🔒 OpenVPN Integration** - Direct system integration with OpenVPN service

### Development Tools
- **📦 npm/yarn** - Package management
- **🔍 ESLint** - Code linting and formatting
- **🐳 Docker** - Containerization support
- **🔄 Hot Reload** - Development server with instant updates

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

## 🌐 Nginx Configuration

This directory includes a production-ready Nginx configuration file (`nginx.conf`) that provides:

### Features
- **Reverse Proxy**: Routes API requests to backend and frontend requests appropriately
- **Rate Limiting**: Protects against abuse with configurable rate limits
- **Security Headers**: Adds essential security headers (X-Frame-Options, X-Content-Type-Options, etc.)
- **Static File Caching**: Optimizes performance with proper cache headers
- **Gzip Compression**: Reduces bandwidth usage
- **Health Checks**: Built-in health check endpoint

### Usage

#### With Docker Compose
The nginx configuration is automatically used when deploying with Docker Compose from the `docker/` directory.

#### Manual Setup
```bash
# Copy to nginx sites directory
sudo cp nginx.conf /etc/nginx/sites-available/openvpn-admin
sudo ln -s /etc/nginx/sites-available/openvpn-admin /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx
```

#### Configuration Details
- **API Routes** (`/api/*`): Proxied to backend with rate limiting (10 req/s)
- **Login Route** (`/api/user/login`): Special rate limiting (1 req/s)
- **Frontend Routes** (`/*`): Proxied to Next.js frontend
- **Static Files**: Cached for 1 year with immutable headers
- **Health Check** (`/health`): Returns simple health status

### SSL/HTTPS Setup
The configuration includes commented SSL sections. To enable HTTPS:

1. Obtain SSL certificates (recommended: Let's Encrypt)
2. Uncomment and configure the HTTPS server block
3. Update certificate paths
4. Reload nginx configuration
