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
                  <Link href="/auth/register">{t("signUp")}</Link>
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
                      d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
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
                      d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
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
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
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
              <Link href="/docs">
                {t("home.featuresSection.viewDocsButton")}
              </Link>
            </Button>
          </div>
        </div>
      </div>

      {/* Tech Stack Section */}
      <div className="bg-gradient-to-r from-secondary to-secondary/50 py-24">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            <h2 className="text-3xl font-semibold mb-6">
              {t("home.techStackSection.title")}
            </h2>
            <div className="h-1 w-16 bg-primary mx-auto mb-8"></div>
            <p className="text-gray-700 mb-10">
              {t("home.techStackSection.subtitle")}
            </p>
            <Button asChild size="lg" className="px-8">
              <Link href="/github" target="_blank" rel="noopener noreferrer">
                {t("home.techStackSection.viewOnGithubButton")}
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </MainLayout>
  );
}
