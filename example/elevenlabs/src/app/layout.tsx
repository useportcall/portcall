import { Metadata } from "next";
import "./globals.css";
import "@repo/ui/styles.css";

export const metadata: Metadata = {
  title: "Portcall example",
  description: "Set up complex billing logic in 10 minutes.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head />
      <body className="">{children}</body>
    </html>
  );
}
