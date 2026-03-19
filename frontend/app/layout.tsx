import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner"; // 1. ADD THIS IMPORT

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Jhinkx | Real-time Matchmaking",
  description: "Find your squad instantly.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        {children}
        {/* 2. ADD THE TOASTER HERE (richColors makes the success toast green!) */}
        <Toaster richColors position="top-center" /> 
      </body>
    </html>
  );
}