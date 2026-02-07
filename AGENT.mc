# Repository Guidelines

## Project Structure
- `cmd/` contains CLI entry points and the interactive menu commands.
- `controller/`, `router/`, `middleware/` implement the HTTP API.
- `model/`, `database/`, `data/` hold persistence logic (SQLite + GORM).
- `openvpn/` contains OpenVPN integration and configuration handling.
- `openvpn-web/` is the Next.js 14 frontend (TypeScript + Tailwind CSS).
- `docker/` contains Dockerfiles, compose files, and deployment docs.
- `config/`, `constants/`, `services/`, `utils/`, `logging/` provide shared logic.

## Build, Test, and Development Commands
- Backend dev: `go run main.go`
- CLI help: `go run cmd/main.go --help`
- Backend build: `go build -o openvpn-go main.go`
- Backend tests: `go test ./...`
- Frontend dev: `cd openvpn-web && npm install && npm run dev`
- Frontend build: `cd openvpn-web && npm run build`
- Frontend lint: `cd openvpn-web && npm run lint`

## Configuration
- Root `.env` (see `.env.example`) for backend settings.
- Docker environment template: `docker/.env.docker.example`.

## Coding Style
- Go code should follow `gofmt` and handle errors explicitly.
- Keep API handlers in `controller/` and wire routes in `router/`.
- Frontend uses TypeScript, Next.js, Tailwind CSS, and Radix UI.

## Testing Notes
- Prefer `go test ./...` for backend changes.
- Frontend has lint and i18n scripts; no dedicated test runner yet.

## PR Guidance
- Include a short summary and note manual testing.
- Call out any OpenVPN or system-level behavior changes in the PR description.
