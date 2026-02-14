import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { ReactNode } from "react";
import { Link } from "react-router-dom";

export function SidebarItem({
  name,
  to,
  icon,
  isCollapsed,
  disabled,
}: {
  name: string;
  to: string;
  icon: ReactNode;
  isCollapsed: boolean;
  disabled?: boolean;
}) {
  const content = (
    <Link
      to={to}
      className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-accent ${
        disabled ? "opacity-50 pointer-events-none" : ""
      } ${isCollapsed ? "justify-center" : ""}`}
    >
      {icon}
      <span className={isCollapsed ? "hidden" : "hidden lg:block"}>{name}</span>
      {disabled && !isCollapsed && <span className="text-xs hidden lg:block">coming soon!</span>}
    </Link>
  );

  if (!isCollapsed) return content;

  return (
    <Tooltip>
      <TooltipTrigger asChild>{content}</TooltipTrigger>
      <TooltipContent side="right">{name}</TooltipContent>
    </Tooltip>
  );
}
