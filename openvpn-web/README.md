# Next.js Template

A modern Next.js template with authentication, UI components, and Docker support.

## Features

- Next.js 14 with TypeScript
- Authentication system
- Modern UI components with Tailwind CSS
- Toast notifications
- Docker support for easy deployment
- ESLint and TypeScript configuration

## Prerequisites

- Node.js 18.x or higher
- npm or yarn
- Docker (optional, for containerized deployment)

## Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/nextjs-template.git
cd nextjs-template
```

2. Install dependencies

```bash
npm install
# or
yarn install
```

3. Create a `.env` file based on `.env.example`

4. Run the development server

```bash
npm run dev
# or
yarn dev
```

5. Open [http://localhost:3000](http://localhost:3000) in your browser

## Docker Deployment

1. Build the Docker image

```bash
docker-compose build
```

2. Start the containers

```bash
docker-compose up -d
```

## Project Structure

- `src/app` - Next.js app router pages
- `src/components` - Reusable UI components
- `src/lib` - Utility functions and configurations
- `src/store` - State management
- `src/assets` - Static assets

## License

MIT
