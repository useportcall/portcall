import { LanguageSwitcher } from "../language-switcher";
import { useGetApp } from "@/hooks/api/apps";
import { AccountMenu } from "./navbar/account-menu";
import { DogfoodBadge } from "./navbar/dogfood-badge";
import { MobileNavSheet } from "./navbar/mobile-nav-sheet";
import { ThemeToggle } from "./navbar/theme-toggle";

export default function Navbar() {
  const { data: app } = useGetApp();
  const isDogfoodApp = app?.data?.billing_exempt === true;

  return (
    <header className="w-full flex h-14 items-center space-mono-regular justify-between md:justify-end gap-4 border-b bg-background px-4 lg:h-[60px] lg:px-6">
      {isDogfoodApp && <DogfoodBadge />}
      <div className="flex items-center gap-2">
        <LanguageSwitcher />
        <ThemeToggle />
        <MobileNavSheet showBilling={!isDogfoodApp} />
        <AccountMenu />
      </div>
    </header>
  );
}
