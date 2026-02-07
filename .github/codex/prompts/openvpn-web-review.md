# Frontend Review (openvpn-web)

Focus areas:
- Next.js 14 App Router behavior and server/client boundaries.
- React hooks usage, state management (Zustand), and effect cleanup.
- API calls, error handling, and loading states.
- i18n keys, fallback behavior, and translation consistency.
- Accessibility, form validation, and user feedback.

Checklist:
- Avoids hydration mismatches and server-only module usage in client code.
- Uses `use client` only where needed.
- Handles empty/error states in UI flows.
- Ensures user input is validated before API calls.

Suggested tests:
- `cd openvpn-web && npm run lint`
- `cd openvpn-web && npm run build`
- `cd openvpn-web && npm run check-i18n` (when translations change)
