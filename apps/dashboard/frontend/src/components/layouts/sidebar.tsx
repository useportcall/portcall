import { useListApps } from "@/hooks";
import { useApp } from "@/hooks/use-app";
import {
  Boxes,
  BriefcaseBusiness,
  Code,
  File,
  Home,
  Link2,
  Repeat,
  ScrollText,
  User,
} from "lucide-react";
import { ReactNode } from "react";
import { Link } from "react-router-dom";
import iconUrl from "../../assets/portcall-icon.png";
import logoUrl from "../../assets/portcall-logo.svg";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { Separator } from "../ui/separator";

export default function Sidebar() {
  return (
    <div className="hidden border-r bg-white text-primary md:block">
      <div className="flex h-full max-h-screen flex-col gap-2">
        <div className="w-full h-14 justify-center items-center border-b px-4 lg:h-[60px] lg:px-6 hidden lg:flex">
          <img src={logoUrl} alt="Portcall logo" className="w-32" />
        </div>
        <div className="w-full h-14 justify-center items-center border-b px-3 flex lg:hidden">
          <img src={iconUrl} alt="Portcall logo" className="w-32" />
        </div>
        <div className="flex- mt-15 lg:mt-0">
          <nav className="grid items-start px-2 text-sm font-medium lg:px-4">
            <ProjectSelect />
            <Sidebaritem
              name="Home"
              to="/"
              icon={<Home className="h-4 w-4" />}
            />
            <Sidebaritem
              name="Users"
              to="/users"
              icon={<User className="h-4 w-4" />}
            />
            <Sidebaritem
              name="Plans"
              to="/plans"
              icon={<Boxes className="h-4 w-4" />}
            />
            <Sidebaritem
              name="Quotes"
              to="/quotes"
              icon={<ScrollText className="h-4 w-4" />}
              disabled
            />
            <Sidebaritem
              name="Subscriptions"
              to="/subscriptions"
              icon={<Repeat className="h-4 w-4" />}
            />
            <Sidebaritem
              name="Invoices"
              to="/invoices"
              icon={<File className="h-4 w-4" />}
            />
            <Separator className="mb-2 hidden lg:block" />
            <Sidebaritem
              name="Company details"
              to="/company"
              icon={<BriefcaseBusiness className="h-4 w-4" />}
            />
            <Sidebaritem
              name="Payment integrations"
              to="/integrations"
              icon={<Link2 className="h-4 w-4" />}
            />
            <Sidebaritem
              name="Developer"
              to="/developer"
              icon={<Code className="h-4 w-4" />}
            />
          </nav>
        </div>
      </div>
    </div>
  );
}

function Sidebaritem({
  name,
  to,
  icon,
  disabled,
}: {
  name: string;
  to: string;
  icon: ReactNode;
  disabled?: boolean;
}) {
  return (
    <Link
      to={to}
      className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50 ${
        disabled ? "opacity-50 pointer-events-none" : ""
      }`}
    >
      {icon}
      <span className="hidden lg:block">{name}</span>
      {disabled && (
        <span className="text-xs hidden lg:block">coming soon!</span>
      )}
    </Link>
  );
}

function ProjectSelect() {
  const { data: apps } = useListApps();
  const app = useApp();

  return (
    <Select
      value={app.id || undefined}
      onValueChange={(value) => app.setId(value)}
    >
      <SelectTrigger className="overflow-auto my-2">
        <SelectValue className="" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel className="text-xs text-slate-400 font-medium">
            Apps
          </SelectLabel>
          {apps?.data.map((app) => (
            <SelectItem key={app.id} value={app.id.toString()}>
              {app.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  );
}
