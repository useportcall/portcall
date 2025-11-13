import { useAuth } from "@/lib/keycloak/auth";
import {
  Boxes,
  BriefcaseBusiness,
  Code,
  File,
  Home,
  Link2,
  Menu,
  Repeat,
  User,
} from "lucide-react";
import { Link } from "react-router-dom";
import logoUrl from "../../assets/portcall-logo.svg";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { Separator } from "../ui/separator";
import { Sheet, SheetContent, SheetTrigger } from "../ui/sheet";
import { Avatar, AvatarFallback } from "../ui/avatar";
import { useMemo } from "react";

export default function Navbar() {
  const { logout, user } = useAuth();

  const firstLetter = useMemo(() => {
    if (user?.name && user.name.length > 0) {
      return user.name.charAt(0).toUpperCase();
    }
    return "U";
  }, [user]);

  return (
    <header className="w-full flex h-14 items-center space-mono-regular justify-between md:justify-end gap-4 border-b bg-white px-4 lg:h-[60px] lg:px-6">
      <Sheet>
        <SheetTrigger asChild>
          <Button variant="outline" size="icon" className="shrink-0 md:hidden">
            <Menu className="h-5 w-5" />
            <span className="sr-only">Toggle navigation menu</span>
          </Button>
        </SheetTrigger>
        <SheetContent side="left" className="flex flex-col bg-white">
          <nav className="grid gap-2 text-lg font-medium">
            <Link
              to="/"
              className="flex flex-row  justify-center items-center gap-2 text-lg font-semibold"
            >
              <img
                src={logoUrl}
                alt="Portcall logo"
                className="w-32 self-center my-2"
              />
            </Link>
            <Separator />
            <Link
              to="/"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <Home className="h-4 w-4" />
              <span className="">Home</span>
            </Link>
            <Link
              to="/users"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <User className="h-4 w-4" />
              <span className="">User</span>
            </Link>
            <Link
              to="/plans"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <Boxes className="h-4 w-4" />
              <span className="">Plans</span>
            </Link>
            <Link
              to="/subscriptions"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <Repeat className="h-4 w-4" />
              <span className="">Subscriptions</span>
            </Link>
            <Link
              to="/invoices"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <File className="h-4 w-4" />
              <span className="">Invoices</span>
            </Link>
            <Link
              to="/company"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <BriefcaseBusiness className="h-4 w-4" />
              <span className="">Company details</span>
            </Link>
            <Link
              to="/integrations"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <Link2 className="h-4 w-4" />
              <span className="">Payment integrations</span>
            </Link>
            <Link
              to="/developer"
              className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-slate-50`}
            >
              <Code className="h-4 w-4" />
              <span className="">Developer</span>
            </Link>
          </nav>
        </SheetContent>
      </Sheet>
      <DropdownMenu>
        <DropdownMenuTrigger>
          <Avatar className="border border-cyan-800/40">
            <AvatarFallback>{firstLetter}</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuItem disabled>{user?.name}</DropdownMenuItem>
          <DropdownMenuItem disabled>{user?.email}</DropdownMenuItem>
          <DropdownMenuItem onClick={() => logout()}>Log out</DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </header>
  );
}
