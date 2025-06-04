"use client";

import React from "react";
import Navbar from "./navbar";

interface MainLayoutProps {
  children: React.ReactNode;
  showFooter?: boolean;
  className?: string;
}

export default function MainLayout({
  children,
  showFooter = true,
  className,
}: MainLayoutProps) {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />

      <main className={`flex-grow ${className}`}>{children}</main>

      {showFooter && (
        <footer className="bg-white border-t border-gray-200 py-6">
          <div className="container mx-auto px-4">
            <div className="text-center">
              <p className="mb-2">
                Next.js Template © {new Date().getFullYear()} All Rights Reserved
              </p>
              <p className="text-sm text-gray-600">
                Built with Next.js, TypeScript, and Tailwind CSS
              </p>
            </div>
          </div>
        </footer>
      )}
    </div>
  );
}
