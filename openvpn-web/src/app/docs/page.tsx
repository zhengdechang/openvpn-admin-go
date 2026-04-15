"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";

const NAV_SECTIONS = ["overview", "quickstart", "envvars", "roles", "api", "faq"] as const;
type Section = (typeof NAV_SECTIONS)[number];

export default function DocsPage() {
  const { t } = useTranslation();
  const [activeSection, setActiveSection] = useState<Section>("overview");

  useEffect(() => {
    const handleScroll = () => {
      for (const id of [...NAV_SECTIONS].reverse()) {
        const el = document.getElementById(id);
        if (el && el.getBoundingClientRect().top <= 100) {
          setActiveSection(id);
          return;
        }
      }
      setActiveSection("overview");
    };
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const scrollTo = (id: Section) => {
    const el = document.getElementById(id);
    if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
    setActiveSection(id);
  };

  // typed data from translation
  const envVars: { name: string; default: string; required: boolean; description: string }[] =
    t("docs.envvars.vars", { returnObjects: true }) as never;
  const roleList: { role: string; permissions: string }[] =
    t("docs.roles.roleList", { returnObjects: true }) as never;
  const endpoints: { method: string; endpoint: string; auth: boolean; description: string }[] =
    t("docs.api.endpoints", { returnObjects: true }) as never;
  const faqItems: { question: string; answer: string }[] =
    t("docs.faq.items", { returnObjects: true }) as never;
  const featureList: string[] =
    t("docs.overview.featureList", { returnObjects: true }) as never;

  const methodColor: Record<string, string> = {
    GET: "bg-blue-100 text-blue-700",
    POST: "bg-green-100 text-green-700",
    PUT: "bg-yellow-100 text-yellow-700",
    DELETE: "bg-red-100 text-red-700",
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Top bar */}
      <div className="border-b bg-white sticky top-0 z-30">
        <div className="container mx-auto px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Link href="/" className="text-primary font-bold text-lg hover:opacity-80">
              VPN Admin
            </Link>
            <span className="text-gray-400">/</span>
            <span className="text-gray-600 font-medium">{t("docs.pageTitle")}</span>
          </div>
          <Button asChild variant="outline" size="sm">
            <Link href="/">&larr; {t("common.back") || "Back"}</Link>
          </Button>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        <div className="flex gap-8">
          {/* Sidebar */}
          <aside className="hidden md:block w-56 shrink-0">
            <div className="sticky top-20">
              <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-4">
                {t("docs.pageTitle")}
              </p>
              <nav className="space-y-1">
                {NAV_SECTIONS.map((s) => (
                  <button
                    key={s}
                    onClick={() => scrollTo(s)}
                    className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${
                      activeSection === s
                        ? "bg-primary text-white font-medium"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                  >
                    {t(`docs.nav.${s}`)}
                  </button>
                ))}
              </nav>
            </div>
          </aside>

          {/* Mobile nav */}
          <div className="md:hidden mb-6 w-full">
            <div className="flex gap-2 overflow-x-auto pb-2">
              {NAV_SECTIONS.map((s) => (
                <button
                  key={s}
                  onClick={() => scrollTo(s)}
                  className={`shrink-0 px-3 py-1.5 rounded-full text-sm border transition-colors ${
                    activeSection === s
                      ? "bg-primary text-white border-primary"
                      : "border-gray-300 text-gray-600 hover:bg-gray-50"
                  }`}
                >
                  {t(`docs.nav.${s}`)}
                </button>
              ))}
            </div>
          </div>

          {/* Main content */}
          <main className="flex-1 min-w-0 space-y-16">
            {/* Page header */}
            <div>
              <h1 className="text-4xl font-bold text-primary mb-3">{t("docs.pageTitle")}</h1>
              <p className="text-lg text-gray-600">{t("docs.pageSubtitle")}</p>
            </div>

            {/* Overview */}
            <section id="overview" className="scroll-mt-24">
              <h2 className="text-2xl font-semibold mb-1">{t("docs.overview.title")}</h2>
              <div className="h-0.5 w-12 bg-primary mb-4" />
              <p className="text-gray-700 mb-6">{t("docs.overview.description")}</p>
              <h3 className="font-semibold mb-3">{t("docs.overview.features")}</h3>
              <ul className="space-y-2 mb-6">
                {Array.isArray(featureList) &&
                  featureList.map((f, i) => (
                    <li key={i} className="flex items-start gap-2 text-gray-700">
                      <span className="mt-1 text-primary shrink-0">&#10003;</span>
                      {f}
                    </li>
                  ))}
              </ul>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div className="card p-4">
                  <p className="text-xs font-semibold uppercase text-gray-400 mb-1">
                    {t("docs.overview.frontend")}
                  </p>
                  <p className="text-sm text-gray-700">Next.js 14 · TypeScript · Tailwind CSS · Radix UI</p>
                </div>
                <div className="card p-4">
                  <p className="text-xs font-semibold uppercase text-gray-400 mb-1">
                    {t("docs.overview.backend")}
                  </p>
                  <p className="text-sm text-gray-700">Go · Gin · PostgreSQL · Goose · JWT</p>
                </div>
              </div>
            </section>

            {/* Quick Start */}
            <section id="quickstart" className="scroll-mt-24">
              <h2 className="text-2xl font-semibold mb-1">{t("docs.quickstart.title")}</h2>
              <div className="h-0.5 w-12 bg-primary mb-6" />

              <div className="space-y-8">
                <div>
                  <h3 className="font-semibold mb-2">{t("docs.quickstart.step1")}</h3>
                  <pre className="bg-gray-900 text-green-400 rounded-lg p-4 text-sm overflow-x-auto">
                    <code>git clone https://github.com/zhengdechang/openvpn-admin-go.git{"\n"}cd openvpn-admin-go</code>
                  </pre>
                </div>

                <div>
                  <h3 className="font-semibold mb-2">{t("docs.quickstart.step2")}</h3>
                  <p className="text-gray-600 text-sm mb-2">{t("docs.quickstart.step2desc")}</p>
                  <pre className="bg-gray-900 text-green-400 rounded-lg p-4 text-sm overflow-x-auto">
                    <code>JWT_SECRET=your-secret-key-here</code>
                  </pre>
                </div>

                <div>
                  <h3 className="font-semibold mb-2">{t("docs.quickstart.step3")}</h3>
                  <p className="text-gray-600 text-sm mb-2">{t("docs.quickstart.step3desc")}</p>
                  <pre className="bg-gray-900 text-green-400 rounded-lg p-4 text-sm overflow-x-auto">
                    <code>docker compose up -d{"\n"}docker logs openvpn-admin | grep -A4 &apos;superadmin&apos;</code>
                  </pre>
                </div>

                <div>
                  <h3 className="font-semibold mb-2">{t("docs.quickstart.step4")}</h3>
                  <p className="text-gray-600 text-sm mb-2">{t("docs.quickstart.step4desc")}</p>
                  <pre className="bg-gray-900 text-green-400 rounded-lg p-4 text-sm overflow-x-auto">
                    <code>http://localhost:80</code>
                  </pre>
                  <p className="text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded-lg p-3 mt-3">
                    {t("docs.quickstart.loginNote")}
                  </p>
                </div>
              </div>
            </section>

            {/* Environment Variables */}
            <section id="envvars" className="scroll-mt-24">
              <h2 className="text-2xl font-semibold mb-1">{t("docs.envvars.title")}</h2>
              <div className="h-0.5 w-12 bg-primary mb-4" />
              <p className="text-gray-700 mb-6">{t("docs.envvars.description")}</p>
              <div className="overflow-x-auto rounded-lg border">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.envvars.name")}</th>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.envvars.default")}</th>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.envvars.required")}</th>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.envvars.description2")}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {Array.isArray(envVars) &&
                      envVars.map((v, i) => (
                        <tr key={i} className="border-b last:border-0 hover:bg-gray-50">
                          <td className="px-4 py-3 font-mono text-xs font-semibold text-primary">{v.name}</td>
                          <td className="px-4 py-3 text-gray-600 font-mono text-xs">{v.default}</td>
                          <td className="px-4 py-3">
                            <span
                              className={`px-2 py-0.5 rounded-full text-xs font-medium ${
                                v.required
                                  ? "bg-red-100 text-red-700"
                                  : "bg-gray-100 text-gray-500"
                              }`}
                            >
                              {v.required ? t("docs.envvars.yes") : t("docs.envvars.no")}
                            </span>
                          </td>
                          <td className="px-4 py-3 text-gray-700">{v.description}</td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>
            </section>

            {/* Roles */}
            <section id="roles" className="scroll-mt-24">
              <h2 className="text-2xl font-semibold mb-1">{t("docs.roles.title")}</h2>
              <div className="h-0.5 w-12 bg-primary mb-4" />
              <p className="text-gray-700 mb-6">{t("docs.roles.description")}</p>
              <div className="overflow-x-auto rounded-lg border">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="text-left px-4 py-3 font-semibold w-40">{t("docs.roles.role")}</th>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.roles.permissions")}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {Array.isArray(roleList) &&
                      roleList.map((r, i) => (
                        <tr key={i} className="border-b last:border-0 hover:bg-gray-50">
                          <td className="px-4 py-3 font-semibold text-primary">{r.role}</td>
                          <td className="px-4 py-3 text-gray-700">{r.permissions}</td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>
            </section>

            {/* API Reference */}
            <section id="api" className="scroll-mt-24">
              <h2 className="text-2xl font-semibold mb-1">{t("docs.api.title")}</h2>
              <div className="h-0.5 w-12 bg-primary mb-4" />
              <p className="text-gray-700 mb-6">{t("docs.api.description")}</p>
              <div className="overflow-x-auto rounded-lg border">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="text-left px-4 py-3 font-semibold w-20">{t("docs.api.method")}</th>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.api.endpoint")}</th>
                      <th className="text-left px-4 py-3 font-semibold w-20">{t("docs.api.auth")}</th>
                      <th className="text-left px-4 py-3 font-semibold">{t("docs.api.description2")}</th>
                    </tr>
                  </thead>
                  <tbody>
                    {Array.isArray(endpoints) &&
                      endpoints.map((ep, i) => (
                        <tr key={i} className="border-b last:border-0 hover:bg-gray-50">
                          <td className="px-4 py-3">
                            <span
                              className={`px-2 py-0.5 rounded text-xs font-bold font-mono ${
                                methodColor[ep.method] || "bg-gray-100 text-gray-700"
                              }`}
                            >
                              {ep.method}
                            </span>
                          </td>
                          <td className="px-4 py-3 font-mono text-xs text-gray-800">{ep.endpoint}</td>
                          <td className="px-4 py-3">
                            {ep.auth ? (
                              <span className="text-xs text-amber-700">
                                {t("docs.api.yes")}
                              </span>
                            ) : (
                              <span className="text-xs text-gray-400">
                                {t("docs.api.no")}
                              </span>
                            )}
                          </td>
                          <td className="px-4 py-3 text-gray-700">{ep.description}</td>
                        </tr>
                      ))}
                  </tbody>
                </table>
              </div>
            </section>

            {/* FAQ */}
            <section id="faq" className="scroll-mt-24 pb-16">
              <h2 className="text-2xl font-semibold mb-1">{t("docs.faq.title")}</h2>
              <div className="h-0.5 w-12 bg-primary mb-6" />
              <div className="space-y-4">
                {Array.isArray(faqItems) &&
                  faqItems.map((item, i) => (
                    <FaqItem key={i} question={item.question} answer={item.answer} />
                  ))}
              </div>
            </section>
          </main>
        </div>
      </div>
    </div>
  );
}

function FaqItem({ question, answer }: { question: string; answer: string }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="border rounded-lg overflow-hidden">
      <button
        className="w-full text-left px-5 py-4 font-medium flex items-center justify-between hover:bg-gray-50 transition-colors"
        onClick={() => setOpen((o) => !o)}
      >
        <span>{question}</span>
        <span className={`text-gray-400 transition-transform ${open ? "rotate-180" : ""}`}>
          &#9660;
        </span>
      </button>
      {open && (
        <div className="px-5 pb-4 pt-1 text-gray-700 text-sm border-t bg-gray-50">
          <code className="whitespace-pre-wrap font-sans">{answer}</code>
        </div>
      )}
    </div>
  );
}
