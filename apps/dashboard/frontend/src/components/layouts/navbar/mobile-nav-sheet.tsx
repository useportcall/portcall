import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Menu } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import logoUrl from "../../../assets/portcall-logo.svg";
import { NAV_ITEMS, SETTINGS_ITEMS } from "../sidebar/navigation-items";

export function MobileNavSheet({ showBilling }: { showBilling: boolean }) {
  const { t } = useTranslation();
  const items = [
    ...NAV_ITEMS,
    ...SETTINGS_ITEMS,
    ...(showBilling ? [{ key: "sidebar.billing", fallback: "Billing", to: "/usage" }] : []),
  ];

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button variant="outline" size="icon" className="shrink-0 md:hidden">
          <Menu className="h-5 w-5" />
          <span className="sr-only">Toggle navigation menu</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="flex flex-col bg-background">
        <nav className="grid gap-2 text-lg font-medium">
          <Link to="/" className="flex flex-row justify-center items-center gap-2 text-lg font-semibold">
            <img src={logoUrl} alt="Portcall logo" className="w-32 self-center my-2 dark:invert" />
          </Link>
          <Separator />
          {items.map((item) => (
            <Link key={item.key} to={item.to} className="flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:bg-accent">
              {"icon" in item ? <item.icon className="h-4 w-4" /> : null}
              <span>{t(item.key, item.fallback)}</span>
            </Link>
          ))}
        </nav>
      </SheetContent>
    </Sheet>
  );
}
