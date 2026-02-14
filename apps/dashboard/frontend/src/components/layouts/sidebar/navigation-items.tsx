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

export const NAV_ITEMS = [
  { key: "sidebar.home", fallback: "Home", to: "/", icon: Home },
  { key: "sidebar.users", fallback: "Users", to: "/users", icon: User },
  { key: "sidebar.plans", fallback: "Plans", to: "/plans", icon: Boxes },
  { key: "sidebar.quotes", fallback: "Quotes", to: "/quotes", icon: ScrollText },
  { key: "sidebar.subscriptions", fallback: "Subscriptions", to: "/subscriptions", icon: Repeat },
  { key: "sidebar.invoices", fallback: "Invoices", to: "/invoices", icon: File },
] as const;

export const SETTINGS_ITEMS = [
  { key: "sidebar.company", fallback: "Company details", to: "/company", icon: BriefcaseBusiness },
  { key: "sidebar.integrations", fallback: "Payment integrations", to: "/integrations", icon: Link2 },
  { key: "sidebar.developer", fallback: "Developer", to: "/developer", icon: Code },
] as const;
