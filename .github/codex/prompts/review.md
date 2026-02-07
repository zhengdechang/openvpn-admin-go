# Backend Review (openvpn-admin-go)

Focus areas:
- Go services, CLI commands, and API handlers.
- OpenVPN integration, file system access, and command execution.
- Database access (GORM), migrations, and data integrity.
- Concurrency, goroutines, context cancellation, and resource leaks.
- Input validation, auth checks, and privilege boundaries.

Checklist:
- Errors are checked and wrapped with context.
- Logging is actionable and avoids leaking secrets.
- Config/env changes are documented when needed.
- Risky operations (file writes, shell execs) are guarded and validated.

Suggested tests:
- `go test ./...`
- Targeted package tests when behavior changes.
