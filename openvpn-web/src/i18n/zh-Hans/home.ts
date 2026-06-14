const translation = {
  title: "VPN Admin — 企业级 OpenVPN 管理平台",
  subtitle: "现代化开源 OpenVPN 管理系统，通过简洁的 Web 界面管理用户、证书、部门和服务器状态，无需命令行操作。",
  getStarted: "立即登录",
  signUp: "注册账号",
  goDashboard: "进入控制台",
  featuresSection: {
    title: "一站式管理能力",
    subtitle: "专为需要可靠、可审计、可扩展 VPN 管理的团队而构建",
    auth: {
      title: "用户与角色管理",
      description: "多角色权限控制（超级管理员 / 管理员 / 经理 / 普通用户），支持部门级隔离和每用户 VPN 证书生命周期管理。"
    },
    ui: {
      title: "实时服务器监控",
      description: "在线状态、VPN IP 分配、流量统计和 OpenVPN 服务器健康状态实时自动刷新。"
    },
    docker: {
      title: "一键部署",
      description: "通过单条 `docker compose up` 命令部署整个技术栈——Go API、Next.js 前端、PostgreSQL 数据库，无需任何外部依赖。"
    },
    viewDocsButton: "查看文档"
  },
  techStackSection: {
    title: "现代轻量技术栈",
    subtitle: "从第一天起就为生产环境设计：nginx 提供静态前端服务，单一 Go 二进制文件处理 API，PostgreSQL 保障数据可靠性。",
    viewOnGithubButton: "查看 GitHub 源码",
    downloadButton: "下载发布版本",
    frontend: {
      title: "前端技术",
      description: "Next.js 16 · React 19 · TypeScript · Tailwind CSS · Radix UI · MUI v5"
    },
    backend: {
      title: "后端技术",
      description: "Go 1.21+ · Gin v1.10 · PostgreSQL 16 · GORM · Goose 迁移 · JWT"
    }
  },
  githubSection: {
    title: "开源项目 & 社区驱动",
    subtitle: "VPN Admin 采用 MIT 许可证，欢迎贡献代码。给仓库加星、提交 Issue，共同参与路线图规划。",
    openSource: {
      title: "MIT 许可证",
      description: "可自由使用、修改和自部署，无供应商锁定。"
    },
    community: {
      title: "社区",
      description: "在 GitHub Issues 参与讨论，分享您的部署经验。"
    },
    documentation: {
      title: "文档",
      description: "安装指南、API 参考和 Docker 部署教程。"
    },
    releases: {
      title: "持续迭代",
      description: "定期发布新功能、安全补丁和 Bug 修复版本。"
    },
    viewSourceButton: "查看 GitHub 源码",
    reportIssueButton: "报告问题"
  },
  contactSection: {
    title: "联系我们",
    subtitle: "有问题、功能建议或发现了 Bug？通过以下任意渠道联系我们。",
    github: {
      title: "GitHub Issues",
      description: "Bug 报告和功能请求的最佳渠道。"
    },
    email: {
      title: "电子邮件",
      description: "私密问询或安全漏洞披露请通过邮件联系。"
    },
    telegram: {
      title: "Telegram",
      description: "加入社区频道获取更新和支持。"
    },
    githubButton: "提交 Issue",
    emailButton: "发送邮件",
    telegramButton: "加入 Telegram"
  }
}

export default translation
