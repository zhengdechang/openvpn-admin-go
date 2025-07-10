"use client";

import React, { useEffect } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth-context";
import { useTranslation } from "react-i18next";
import MainLayout from "@/components/layout/main-layout";

export default function Home() {
  const { user, loading, refreshToken } = useAuth();
  const { t } = useTranslation();

  useEffect(() => {
    refreshToken();
  }, []);

  return (
    <MainLayout>
      {/* Hero Section */}
      <div className="hero-pattern py-24">
        <div className="container mx-auto px-4">
          <div className="text-center max-w-3xl mx-auto">
            <h1 className="text-4xl md:text-5xl font-bold mb-4 text-primary">
              {t("home.title")}
            </h1>
            <div className="h-1 w-24 bg-accent mx-auto my-6"></div>
            <p className="text-lg text-gray-700 mb-10">{t("home.subtitle")}</p>
            {!user ? (
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-4">
                <Button asChild size="lg" className="w-full sm:w-auto">
                  <Link href="/auth/login">{t("home.getStarted")}</Link>
                </Button>
                <Button
                  asChild
                  variant="outline"
                  size="lg"
                  className="w-full sm:w-auto"
                >
                  <Link href="/auth/register">{t("home.signUp")}</Link>
                </Button>
              </div>
            ) : (
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-4">
                <Button asChild size="lg" className="w-full sm:w-auto">
                  <Link href="/dashboard">{t("home.goDashboard")}</Link>
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="bg-white py-24">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-3xl font-semibold mb-6">
              {t("home.featuresSection.title")}
            </h2>
            <div className="h-1 w-16 bg-primary mx-auto mb-8"></div>
            <p className="text-gray-700 mb-10">
              {t("home.featuresSection.subtitle")}
            </p>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
              <div className="card p-6">
                <div className="text-primary text-4xl mb-4 flex justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-16 w-16"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"
                    />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.featuresSection.auth.title")}
                </h3>
                <p className="text-gray-600">
                  {t("home.featuresSection.auth.description")}
                </p>
              </div>
              <div className="card p-6">
                <div className="text-primary text-4xl mb-4 flex justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-16 w-16"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z"
                    />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.featuresSection.ui.title")}
                </h3>
                <p className="text-gray-600">
                  {t("home.featuresSection.ui.description")}
                </p>
              </div>
              <div className="card p-6">
                <div className="text-primary text-4xl mb-4 flex justify-center">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-16 w-16"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1.5}
                      d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                    />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.featuresSection.docker.title")}
                </h3>
                <p className="text-gray-600">
                  {t("home.featuresSection.docker.description")}
                </p>
              </div>
            </div>
            <Button asChild size="lg" className="px-8">
              <Link href="https://github.com/zhengdechang/openvpn-admin-go#readme" target="_blank" rel="noopener noreferrer">
                {t("home.featuresSection.viewDocsButton")}
              </Link>
            </Button>
          </div>
        </div>
      </div>

      {/* Tech Stack Section */}
      <div className="bg-gradient-to-r from-secondary to-secondary/50 py-24">
        <div className="container mx-auto px-4">
          <div className="max-w-5xl mx-auto text-center">
            <h2 className="text-3xl font-semibold mb-6">
              {t("home.techStackSection.title")}
            </h2>
            <div className="h-1 w-16 bg-primary mx-auto mb-8"></div>
            <p className="text-gray-700 mb-12 text-lg">
              {t("home.techStackSection.subtitle")}
            </p>

            {/* Frontend & Backend Tech Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
              {/* Frontend Card */}
              <div className="card p-8 bg-white/80 backdrop-blur-sm">
                <div className="text-primary text-4xl mb-4 flex justify-center">
                  <svg className="h-16 w-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
                <h3 className="text-xl font-semibold mb-3">
                  {t("home.techStackSection.frontend.title")}
                </h3>
                <p className="text-gray-600 mb-4">
                  {t("home.techStackSection.frontend.description")}
                </p>
                <div className="flex flex-wrap justify-center gap-2">
                  <span className="px-3 py-1 bg-black text-white text-sm rounded-full">Next.js 14</span>
                  <span className="px-3 py-1 bg-blue-600 text-white text-sm rounded-full">TypeScript</span>
                  <span className="px-3 py-1 bg-cyan-500 text-white text-sm rounded-full">Tailwind CSS</span>
                  <span className="px-3 py-1 bg-purple-600 text-white text-sm rounded-full">Radix UI</span>
                </div>
              </div>

              {/* Backend Card */}
              <div className="card p-8 bg-white/80 backdrop-blur-sm">
                <div className="text-primary text-4xl mb-4 flex justify-center">
                  <svg className="h-16 w-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
                  </svg>
                </div>
                <h3 className="text-xl font-semibold mb-3">
                  {t("home.techStackSection.backend.title")}
                </h3>
                <p className="text-gray-600 mb-4">
                  {t("home.techStackSection.backend.description")}
                </p>
                <div className="flex flex-wrap justify-center gap-2">
                  <span className="px-3 py-1 bg-cyan-600 text-white text-sm rounded-full">Go</span>
                  <span className="px-3 py-1 bg-green-600 text-white text-sm rounded-full">Gin</span>
                  <span className="px-3 py-1 bg-blue-800 text-white text-sm rounded-full">SQLite</span>
                  <span className="px-3 py-1 bg-orange-600 text-white text-sm rounded-full">OpenVPN</span>
                </div>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Button asChild  size="lg" className="px-8">
                <Link href="https://github.com/zhengdechang/openvpn-admin-go/releases" target="_blank" rel="noopener noreferrer">
                  <svg className="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
                  </svg>
                  {t("home.techStackSection.downloadButton")}
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* GitHub & Community Section */}
      <div className="bg-white py-24">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center">
            <h2 className="text-3xl font-semibold mb-6">
              {t("home.githubSection.title")}
            </h2>
            <div className="h-1 w-16 bg-primary mx-auto mb-8"></div>
            <p className="text-gray-700 mb-10 text-lg">
              {t("home.githubSection.subtitle")}
            </p>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
              <div className="card p-6 text-center">
                <div className="text-primary text-3xl mb-4 flex justify-center">
                  <svg className="h-12 w-12" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.githubSection.openSource.title")}
                </h3>
                <p className="text-gray-600 text-sm">
                  {t("home.githubSection.openSource.description")}
                </p>
              </div>

              <div className="card p-6 text-center">
                <div className="text-primary text-3xl mb-4 flex justify-center">
                  <svg className="h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.githubSection.community.title")}
                </h3>
                <p className="text-gray-600 text-sm">
                  {t("home.githubSection.community.description")}
                </p>
              </div>

              <div className="card p-6 text-center">
                <div className="text-primary text-3xl mb-4 flex justify-center">
                  <svg className="h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.githubSection.documentation.title")}
                </h3>
                <p className="text-gray-600 text-sm">
                  {t("home.githubSection.documentation.description")}
                </p>
              </div>

              <div className="card p-6 text-center">
                <div className="text-primary text-3xl mb-4 flex justify-center">
                  <svg className="h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold mb-2">
                  {t("home.githubSection.releases.title")}
                </h3>
                <p className="text-gray-600 text-sm">
                  {t("home.githubSection.releases.description")}
                </p>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Button asChild size="lg" className="px-8">
                <Link href="https://github.com/zhengdechang/openvpn-admin-go" target="_blank" rel="noopener noreferrer">
                  <svg className="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                  {t("home.githubSection.viewSourceButton")}
                </Link>
              </Button>
              <Button asChild variant="outline" size="lg" className="px-8">
                <Link href="https://github.com/zhengdechang/openvpn-admin-go/issues" target="_blank" rel="noopener noreferrer">
                  <svg className="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  {t("home.githubSection.reportIssueButton")}
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
