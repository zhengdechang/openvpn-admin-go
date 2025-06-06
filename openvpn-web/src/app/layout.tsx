/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-06-05 13:07:03
 */
import React from "react";
import "./globals.css";
import { Inter } from "next/font/google";
import { AuthProvider } from "@/lib/auth-context";
import { Toaster } from "sonner";
import { LoadingScreen } from "@/components/ui/loading";
import { cookies } from "next/headers";
import { i18n } from "@/i18n";

const inter = Inter({ subsets: ["latin"] });

export default async function RootLayout({
  // Make it an async function
  children,
}: {
  children: React.ReactNode;
}) {
  const cookieStore = cookies();
  const locale = cookieStore.get("locale")?.value || i18n.defaultLocale;

  return (
    <html lang={locale} suppressHydrationWarning>
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta name="description" content="OpenVPN Management System" />
        <title>OpenVPN Management System</title>
      </head>
      <body className={inter.className} suppressHydrationWarning>
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
      </body>
    </html>
  );
}
