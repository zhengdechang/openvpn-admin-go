# openvpn-web

This is the web frontend for the `openvpn-admin-go` project. It provides a user interface for managing OpenVPN servers and users.

## Features

- User authentication
- Dashboard for managing OpenVPN servers
- User management interface
- Client configuration download

## Prerequisites

- Node.js 18.x or higher
- npm or yarn
- A running instance of the `openvpn-admin-go` backend.

## Installation

- `git clone <repository-url>`
- `cd openvpn-web`
- `npm install` (or `yarn install`)
- Create a `.env` file based on `.env.example` and configure the backend API URL.
- `npm run dev` (or `yarn dev`)

## Project Structure

- `src/app` - Next.js app router pages and layouts
- `src/components` - Reusable UI components (buttons, forms, etc.)
- `src/services` - API communication and data fetching logic
- `src/lib` - Utility functions, configurations, and helper scripts
- `src/store` - State management (e.g., Zustand, Redux)
- `src/assets` - Static assets like images and fonts
- `public` - Publicly accessible static files
