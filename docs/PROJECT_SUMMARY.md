# OpenVPN Admin Go – Repository Summary

## Project Overview
OpenVPN Admin Go delivers a full-stack, self-hosted control plane for OpenVPN deployments. The platform combines a Go REST API, a Next.js dashboard, terminal tooling, and automation scripts so operators can provision servers, onboard clients, monitor sessions, and manage credentials from a single codebase. All assets required for local development and production rollout—configuration templates, Docker artifacts, database migrations, and background services—are committed to this repository to keep the experience reproducible.

## Backend Architecture
### API Surface
* Gin-powered routers in `router/` register grouped REST endpoints for authentication, client and department management, server lifecycle control, audit logging, and health probes. Each route delegates to the corresponding controller in `controller/` which houses the business rules and response shaping.
* Middleware in `middleware/` layers JWT validation and role-based authorization over the handlers, ensuring that sensitive routes (e.g. client revocation or server restarts) are restricted to elevated roles.

### Data & Security
* GORM models in `model/` define the schema for users, departments, OpenVPN server snapshots, status records, and client activity logs. Database bootstrapping lives in `database/`, providing initialization, auto-migration, and default data seeding (including a bootstrap super admin).
* Password hashing, JWT issuance, and audit-friendly structured logging are centralized in `common/` and `logging/` so controllers remain focused on orchestration rather than plumbing.

### OpenVPN Integration
* The `openvpn/` package encapsulates interactions with the OpenVPN runtime: generating configuration files, creating and revoking client certificates, maintaining client-specific CCD files, parsing live status logs, and synchronizing pause/resume actions through the management interface.
* Background synchronization in `services/openvpn_sync.go` tails the server status log, normalizes connection statistics (bytes transferred, uptime, last activity), and persists them for dashboards.

## Frontend Architecture
* The `openvpn-web/` workspace hosts a Next.js 14 App Router application. Authenticated routes under `src/app/dashboard` surface overview cards, client tables, download links, and configuration editors backed by Zustand stores and Axios-based API clients.
* A shared component library (`src/components/ui`) wraps Tailwind CSS primitives for consistent theming, while `src/lib/auth-context.tsx` and `src/i18n` provide session state and multi-language resources.
* Dashboard pages now expose quick-glance metrics, inline filtering (search, department, status), and badge-based status indicators so operators can triage accounts without paging through full tables.
* Build tooling includes ESLint, Prettier, and a Vite-compatible asset pipeline. Production deployments render behind the provided `nginx.conf` reverse proxy template.

## Command-Line & Automation Tooling
* The interactive CLI in `cmd/` (powered by Cobra and PromptUI) exposes menus for server management, client lifecycle operations, and web service control. Helper utilities in `utils/` provide terminal input helpers, random secret generation, and shell command execution wrappers.
* Environment preparation commands verify prerequisite binaries (OpenVPN, OpenSSL, supervisor), scaffold certificate directories, generate TLS materials, and install supervisor jobs so that administrators can bootstrap a node from scratch.

## Background Services & Observability
* `logging/logging.go` configures structured logging with level control and file sinks that are reused by the API, CLI, and background workers.
* Health endpoints and status inspectors provide coarse-grained monitoring hooks, while supervisor integration files under `docker/` and `template/` keep long-running services managed.

## Deployment & Operations
* Dockerfiles support both combined (API + frontend + OpenVPN) and multi-service deployments. Compose files and supervisor configuration demonstrate how to run the stack either as discrete containers or as a single host-managed bundle.
* Configuration scaffolding under `config/` (including server JSON templates and logging profiles) shows the canonical environment variables and file locations expected in production.
* Documentation in `docs/DEPLOYMENT.md` and the README pair describes environment setup, TLS key management, reverse proxy hardening, and CI/CD recommendations.

## Implemented Features
* **User & Access Management** – CRUD for users and departments, password hashing, JWT-based login, and role-based authorization covering super-admin, admin, and standard operator personas.
* **Client Lifecycle Automation** – Creation, deletion, pause/resume, status inspection, fixed IP assignment, subnet routing via CCD files, and automatic synchronization between database users and generated OpenVPN client artifacts.
* **Server Administration** – Start/stop/restart flows, configuration regeneration, TLS asset rotation, port and network reconfiguration, supervisor integration, and CLI/HTTP endpoints that mirror one another.
* **Observability & Reporting** – Background sync worker ingesting OpenVPN status logs, human-readable byte counters, uptime calculations, audit log endpoints, and dashboard visualizations in the frontend.
* **Operations Tooling** – One-command environment checks, certificate generation scripts, Docker deployment examples, and localization-ready frontend assets.
* **Menu-Driven Docker Bootstrap** – The `openvpn-go` CLI doubles as an in-container installer, wiring supervisor jobs for the API (`openvpn-go-api`) and the Nginx frontend so new hosts can install dependencies and selectively launch services from a single menu.
* **Dashboard Experience** – Aggregated user statistics, client status badges, configurable table filters, and streamlined quick actions (pause/resume, download) keep the web console efficient for day-to-day administration.

## Not Yet Implemented / Future Opportunities
* **Certificate Revocation Workflow** – The frontend ships a `revokeClient` API call, but the backend router lacks a matching endpoint and CRL automation, leaving revocation as a manual task.
* **IPv6 Enhancements** – Data models anticipate IPv6 addresses, yet CCD helpers and validation routines remain IPv4-only.
* **Automated Backups & Rotation** – Certificates, keys, and CCD files are generated locally without scheduled backups or rotation policies. Integrating snapshot/backup jobs would mitigate operator error.
* **Test Coverage & QA Pipelines** – Aside from parser unit tests, most packages lack automated tests or linting pipelines. Adding coverage for controllers, CLI workflows, and OpenVPN integrations would improve regression safety.
* **Telemetry & Analytics** – No metrics export or structured telemetry (Prometheus, OpenTelemetry) is configured. Instrumentation would make capacity planning and troubleshooting easier.
* **UI/UX Polish** – Dark mode, per-column sorting, and richer empty states remain on the wish list to further modernize the Next.js dashboard.

## Testing & Quality Gates
* **Backend** – `go test ./...` validates compile-time health across packages, including the CLI, OpenVPN helpers, and service integrations.
* **Frontend** – `yarn lint` enforces ESLint + Next.js conventions, while `yarn build` exercises production bundling before deployment.

## Developer Workflow Notes
* Go modules target Go 1.18+, though the project builds successfully on 1.24 and later. Run `go test ./...` to compile and execute unit tests (building the CLI pulls in CGO-backed sqlite dependencies on the first run).
* The frontend uses Yarn with Zero-Install configuration (`.yarn/` directory). Use `yarn install` followed by `yarn dev` for local development or `yarn build` for production bundles.
* Docker Compose definitions under `docker/` allow spinning up a full stack with a single command once TLS materials and environment variables are in place.
