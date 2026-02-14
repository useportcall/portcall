// import { use } from "@/hooks/api";
import { cn } from "@/lib/utils";
import {
  ReactNode,
  useState,
  useEffect,
  createContext,
  useContext,
} from "react";
import { useLocation } from "react-router-dom";
import Navbar from "./navbar";
import Sidebar from "./sidebar";

interface SidebarContextType {
  isCollapsed: boolean;
  setIsCollapsed: (collapsed: boolean) => void;
}

export const SidebarContext = createContext<SidebarContextType>({
  isCollapsed: false,
  setIsCollapsed: () => {},
});

export const useSidebar = () => useContext(SidebarContext);

export default function Layout({ children }: { children: ReactNode }) {
  const location = useLocation();
  const [isCollapsed, setIsCollapsed] = useState(() => {
    const stored = localStorage.getItem("sidebar-collapsed");
    return stored === "true";
  });

  useEffect(() => {
    localStorage.setItem("sidebar-collapsed", String(isCollapsed));
  }, [isCollapsed]);

  return (
    <SidebarContext.Provider value={{ isCollapsed, setIsCollapsed }}>
      <main className="min-h-screen w-screen">
        <div
          className={cn(
            "grid h-full w-full transition-[grid-template-columns] duration-200 ease-out md:grid-cols-[56px_auto]",
            isCollapsed
              ? "lg:grid-cols-[56px_auto]"
              : "lg:grid-cols-[240px_auto]",
          )}
        >
          <Sidebar />
          <div className={cn("flex flex-col h-screen w-full")}>
            <Navbar />
            <div className="flex flex-1 flex-col lg:gap-6 lg:p-6 h-full w-full overflow-auto">
              <div key={location.pathname} className="route-fade-in h-full w-full">
                {children}
              </div>
            </div>
          </div>
        </div>
      </main>
    </SidebarContext.Provider>
  );
}
