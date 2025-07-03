/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-05 13:07:03
 */
import React from "react";
import "./globals.css";
import { AuthProvider } from "@/lib/auth-context";
import { Toaster } from "sonner";
import { LoadingScreen } from "@/components/ui/loading";
import { getLocaleOnServer } from "@/i18n/server";
import { I18nProvider } from "@/i18n/i18n-provider";

export default async function RootLayout({
  // Make it an async function
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = getLocaleOnServer();

  return (
    <html lang={locale} suppressHydrationWarning>
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="OpenVPN Management System" />
        <title>OpenVPN Management System</title>
      </head>
      <body className="font-sans" suppressHydrationWarning>
        <I18nProvider locale={locale}>
          <LoadingScreen />
          <AuthProvider>
            <div>{children}</div>
            <Toaster
              position="top-right"
              expand={true}
              visibleToasts={6}
              closeButton={true}
              richColors={true}
              toastOptions={{
                duration: 5000,
                className: "toast-message",
                style: {
                  marginBottom: "0.5rem",
                },
              }}
            />
          </AuthProvider>
        </I18nProvider>
      </body>
    </html>
  );
}
