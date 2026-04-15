const translation = {
  title: "VPN Admin — Enterprise OpenVPN Management",
  subtitle: "A modern, open-source OpenVPN administration platform. Manage users, certificates, departments, and server status through a clean web interface — no CLI required.",
  getStarted: "Sign In",
  signUp: "Create Account",
  goDashboard: "Go to Dashboard",
  featuresSection: {
    title: "Everything You Need",
    subtitle: "Built for teams that need reliable, auditable, and scalable VPN management",
    auth: {
      title: "User & Role Management",
      description: "Multi-role access control (Superadmin / Admin / Manager / User) with department-level isolation and per-user VPN certificate lifecycle management."
    },
    ui: {
      title: "Real-Time Server Monitoring",
      description: "Live connection status, VPN IP assignments, traffic counters, and OpenVPN server health — all refreshed automatically."
    },
    docker: {
      title: "One-Command Deployment",
      description: "Ship the entire stack — Go API, Next.js frontend, and PostgreSQL — with a single `docker compose up`. No external dependencies."
    },
    viewDocsButton: "View Documentation"
  },
  techStackSection: {
    title: "Modern, Lightweight Tech Stack",
    subtitle: "Designed for production from day one: statically exported frontend served by nginx, a single Go binary for the API, and PostgreSQL for reliability.",
    viewOnGithubButton: "View on GitHub",
    downloadButton: "Download Release",
    frontend: {
      title: "Frontend",
      description: "Next.js 14 · TypeScript · Tailwind CSS · Radix UI"
    },
    backend: {
      title: "Backend",
      description: "Go · Gin · PostgreSQL · Goose Migrations · JWT"
    }
  },
  githubSection: {
    title: "Open Source & Community Driven",
    subtitle: "VPN Admin is MIT-licensed and welcomes contributions. Star the repo, open issues, and help shape the roadmap.",
    openSource: {
      title: "MIT License",
      description: "Free to use, modify, and self-host. No vendor lock-in."
    },
    community: {
      title: "Community",
      description: "Join discussions on GitHub Issues and share your deployment stories."
    },
    documentation: {
      title: "Documentation",
      description: "Installation guides, API references, and Docker deployment walkthroughs."
    },
    releases: {
      title: "Active Development",
      description: "Regular releases with new features, security patches, and bug fixes."
    },
    viewSourceButton: "View Source on GitHub",
    reportIssueButton: "Report an Issue"
  },
  contactSection: {
    title: "Get in Touch",
    subtitle: "Have questions, feature requests, or found a bug? Reach out through any of the channels below.",
    github: {
      title: "GitHub Issues",
      description: "The best place for bug reports and feature requests."
    },
    email: {
      title: "Email",
      description: "For private inquiries or security disclosures."
    },
    telegram: {
      title: "Telegram",
      description: "Join the community channel for updates and support."
    },
    githubButton: "Open an Issue",
    emailButton: "Send Email",
    telegramButton: "Join on Telegram"
  }
}

export default translation
