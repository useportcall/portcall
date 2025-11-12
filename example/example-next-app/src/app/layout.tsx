import { Toaster } from "@/components/ui/sonner";
import { Metadata } from "next";
import "./globals.css";
import Providers from "./providers";

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
        <Providers>{children}</Providers>
        <Toaster />
      </body>
    </html>
  );
}
