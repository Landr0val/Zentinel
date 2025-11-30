"use client";

import { Geist, Geist_Mono } from "next/font/google";
import { usePathname } from "next/navigation";
import "./globals.css";
import Providers from "./providers";
import Sidebar from "../components/Sidebar";
import AuthGuard from "../components/auth/AuthGuard";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const pathname = usePathname();
  const isLoginPage = pathname === "/login";

  return (
    <html lang="en">
      <head>
        <title>Zentinel - Fraud Monitoring Dashboard</title>
        <meta
          name="description"
          content="Real-time transaction monitoring and fraud detection system"
        />
      </head>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased bg-background text-foreground`}
      >
        <Providers>
          <AuthGuard>
            <div className="flex min-h-screen">
              {!isLoginPage && <Sidebar />}
              <main
                className={`flex-1 transition-all duration-300 ease-in-out ${
                  !isLoginPage ? "ml-64 p-10" : "w-full"
                }`}
              >
                <div className={!isLoginPage ? "max-w-7xl mx-auto" : ""}>
                  {children}
                </div>
              </main>
            </div>
          </AuthGuard>
        </Providers>
      </body>
    </html>
  );
}
