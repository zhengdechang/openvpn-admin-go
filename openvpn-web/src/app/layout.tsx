import React from 'react';
import './globals.css';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { AuthProvider } from '@/lib/auth-context';
import { Toaster } from 'sonner';
import { useTranslation, getLocaleOnServer } from '@/i18n/server'; // Import server-side translation utils
import type { Locale } from '@/i18n';

const inter = Inter({ subsets: ['latin'] });

export async function generateMetadata({ params }: { params: { lang: Locale } }): Promise<Metadata> {
  // Determine locale - this might need adjustment based on how lang is passed or detected
  // For now, assuming getLocaleOnServer can be used or lang is available via params
  const locale = getLocaleOnServer(); // Or params.lang if available and reliable
  const { t } = await useTranslation(locale, 'layout');
  return {
    title: t('metadataTitle'),
    description: t('metadataDescription'),
  };
}

export default async function RootLayout({ // Make it an async function
  children,
}: {
  children: React.ReactNode;
}) {
  const locale = getLocaleOnServer(); // Get current locale

  return (
    <html lang={locale}> {/* Use determined locale */}
      <body className={inter.className}>
        <AuthProvider>
          {children}
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
                marginBottom: '0.5rem'
              }
            }}
          />
        </AuthProvider>
      </body>
    </html>
  );
} 