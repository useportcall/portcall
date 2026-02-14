import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { ChevronsLeft, ChevronsRight } from "lucide-react";
import { useTranslation } from "react-i18next";
import iconUrl from "../../assets/portcall-icon.png";
import logoUrl from "../../assets/portcall-logo.svg";
import { useSidebar } from "./layout";
import { SidebarNav } from "./sidebar/sidebar-nav";

export default function Sidebar() {
  const { isCollapsed, setIsCollapsed } = useSidebar();
  const { t } = useTranslation();

  return (
    <div className="hidden border-r bg-background text-primary md:block transition-all duration-200 ease-out">
      <div className="flex h-full max-h-screen flex-col gap-2">
        <div className={`w-full h-14 lg:h-[60px] justify-center items-center border-b px-4 lg:px-6 transition-all duration-200 ${isCollapsed ? "hidden" : "hidden lg:flex"}`}>
          <img src={logoUrl} alt="Portcall logo" className="w-32 dark:invert" />
        </div>
        <div className={`w-full h-14 lg:h-[60px] justify-center items-center border-b px-3 transition-all duration-200 ${isCollapsed ? "flex" : "flex lg:hidden"}`}>
          <img src={iconUrl} alt="Portcall logo" className="w-8" />
        </div>
        <div className="flex-1"><SidebarNav isCollapsed={isCollapsed} /></div>
        <div className="hidden lg:flex border-t p-2 justify-center">
          <Tooltip>
            <TooltipTrigger asChild>
              <Button variant="ghost" size="icon" onClick={() => setIsCollapsed(!isCollapsed)} className="h-8 w-8">
                {isCollapsed ? <ChevronsRight className="h-4 w-4" /> : <ChevronsLeft className="h-4 w-4" />}
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right">{isCollapsed ? t("sidebar.expand") : t("sidebar.collapse")}</TooltipContent>
          </Tooltip>
        </div>
      </div>
    </div>
  );
}
