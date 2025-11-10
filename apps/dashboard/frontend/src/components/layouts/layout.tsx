// import { use } from "@/hooks/api";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";
import Navbar from "./navbar";
import Sidebar from "./sidebar";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <main className="min-h-screen w-screen">
      <div className="grid h-full w-full md:grid-cols-[56px_auto] lg:grid-cols-[240px_auto]">
        <Sidebar />
        <div className={cn("flex flex-col h-screen w-full")}>
          <Navbar />
          <div className="flex flex-1 flex-col lg:gap-6 lg:p-6 h-full w-full overflow-auto">
            {children}
          </div>
        </div>
      </div>
    </main>
  );
}
