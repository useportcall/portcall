import { Separator } from "@/components/ui/separator";
import { useTranslation } from "react-i18next";
import { AppSelector } from "./app-selector";
import { NAV_ITEMS, SETTINGS_ITEMS } from "./navigation-items";
import { SidebarItem } from "./sidebar-item";
import { UsageSidebarItem } from "./usage-sidebar-item";

export function SidebarNav({ isCollapsed }: { isCollapsed: boolean }) {
  const { t } = useTranslation();

  return (
    <nav className={`grid items-start px-2 text-sm font-medium ${isCollapsed ? "px-2" : "lg:px-4"}`}>
      <div className="hidden lg:block">
        <AppSelector iconTrigger={isCollapsed} createOptionOnlyWhenEmpty />
      </div>
      <div className="lg:hidden">
        <AppSelector iconTrigger createOptionOnlyWhenEmpty />
      </div>
      {NAV_ITEMS.map(({ key, fallback, to, icon: Icon }) => (
        <SidebarItem key={key} name={t(key, fallback)} to={to} icon={<Icon className="h-4 w-4" />} isCollapsed={isCollapsed} />
      ))}
      <Separator className={`mb-2 ${isCollapsed ? "hidden" : "hidden lg:block"}`} />
      {SETTINGS_ITEMS.map(({ key, fallback, to, icon: Icon }) => (
        <SidebarItem key={key} name={t(key, fallback)} to={to} icon={<Icon className="h-4 w-4" />} isCollapsed={isCollapsed} />
      ))}
      <UsageSidebarItem isCollapsed={isCollapsed} />
    </nav>
  );
}
