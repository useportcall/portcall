import { useGetApp } from "@/hooks";
import { BarChart3 } from "lucide-react";
import { useTranslation } from "react-i18next";
import { SidebarItem } from "./sidebar-item";

export function UsageSidebarItem({ isCollapsed }: { isCollapsed: boolean }) {
  const { data: app } = useGetApp();
  const { t } = useTranslation();

  if (app?.data?.billing_exempt) return null;

  return (
    <SidebarItem
      name={t("sidebar.billing", "Billing")}
      to="/usage"
      icon={<BarChart3 className="h-4 w-4" />}
      isCollapsed={isCollapsed}
    />
  );
}
