import React from 'react';
import './globals.css';
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import { AuthProvider } from '@/lib/auth-context';
import { Toaster } from 'sonner';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Next.js Template',
  description: 'A modern Next.js template with authentication and UI components',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
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