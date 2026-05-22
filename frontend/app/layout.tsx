import type { Metadata } from "next";
import "./globals.css";

// This file defines the shared HTML shell for all App Router pages.

export const metadata: Metadata = {
  title: "Artist Distribution Study",
  description: "Music distribution management study app",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ko">
      <body>{children}</body>
    </html>
  );
}
