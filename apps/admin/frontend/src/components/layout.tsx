import { ReactNode } from "react";
import { Link, useLocation } from "react-router-dom";
import { cn } from "@/lib/utils";
import {
  LayoutDashboard,
  AppWindow,
  ListTodo,
  LogOut,
  Shield,
  ClipboardList,
  Dog,
  CreditCard,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/lib/keycloak/auth";

interface LayoutProps {
  children: ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const { logout, user } = useAuth();
  const location = useLocation();

  const navItems = [
    { name: "Dashboard", href: "/", icon: LayoutDashboard },
    { name: "Apps", href: "/apps", icon: AppWindow },
    { name: "Connections", href: "/connections", icon: CreditCard },
    { name: "Dogfood", href: "/dogfood", icon: Dog },
    { name: "Queue Jobs", href: "/jobs", icon: ClipboardList },
    { name: "Trigger Tasks", href: "/queues", icon: ListTodo },
  ];

  return (
    <div className="flex h-screen w-screen overflow-hidden">
      {/* Sidebar */}
      <div className="w-64 flex-shrink-0 bg-sidebar text-sidebar-foreground flex flex-col">
        <div className="p-4 flex items-center gap-2">
          <Shield className="h-8 w-8" />
          <span className="text-xl font-bold">Admin Panel</span>
        </div>
        <Separator className="bg-sidebar-border" />
        <ScrollArea className="flex-1">
          <nav className="p-4 space-y-2">
            {navItems.map((item) => {
              const isActive =
                item.href === "/"
                  ? location.pathname === "/"
                  : location.pathname.startsWith(item.href);
              return (
                <Link
                  key={item.href}
                  to={item.href}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-lg transition-colors",
                    isActive
                      ? "bg-sidebar-accent text-sidebar-accent-foreground"
                      : "hover:bg-sidebar-accent/50",
                  )}
                >
                  <item.icon className="h-5 w-5" />
                  <span>{item.name}</span>
                </Link>
              );
            })}
          </nav>
        </ScrollArea>
        <Separator className="bg-sidebar-border" />
        <div className="p-4 space-y-2">
          {user?.email && (
            <div className="text-sm text-sidebar-foreground/70 truncate">
              {user.email}
            </div>
          )}
          <Button
            variant="ghost"
            className="w-full justify-start text-sidebar-foreground hover:bg-sidebar-accent/50"
            onClick={logout}
          >
            <LogOut className="h-4 w-4 mr-2" />
            Logout
          </Button>
        </div>
      </div>

      {/* Main Content */}
      <main className="flex-1 flex flex-col overflow-hidden bg-background">
        <ScrollArea className="flex-1">
          <div className="p-8">{children}</div>
        </ScrollArea>
      </main>
    </div>
  );
}
