import { Toaster } from "@/components/ui/sonner";
import { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Example Clay billing app",
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
      <body>
        {children}
        <Toaster />
      </body>
    </html>
  );
}
