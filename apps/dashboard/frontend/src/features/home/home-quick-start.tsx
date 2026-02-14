import { Button } from "@/components/ui/button";
import { BriefcaseBusiness, Boxes, Search, User } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router";

export function HomeQuickStart() {
  const { t } = useTranslation();

  return (
    <div className="w-full flex flex-col gap-4 justify-center items-center h-full">
      <p className="text-lg font-medium">{t("home.quick_start.title")}</p>
      <div
        data-testid="home-quick-start"
        className="h-fit w-full lg:max-w-lg justify-center items-center my-auto mx-auto grid grid-cols-2 gap-4 bg-muted p-4 md:p-10 rounded-md"
      >
        <Link to="/plans" className="w-full">
          <Button
            variant="outline"
            className="flex flex-col w-full justify-center h-full gap-2"
          >
            <Boxes />
            {t("home.quick_start.plans")}
          </Button>
        </Link>
        <Link to="/users" className="w-full">
          <Button
            variant="outline"
            className="flex flex-col justify-center h-full w-full gap-2"
          >
            <User />
            {t("home.quick_start.users")}
          </Button>
        </Link>
        <Link to="/company" className="w-full">
          <Button
            variant="outline"
            className="flex flex-col justify-center h-full w-full gap-2"
          >
            <BriefcaseBusiness />
            {t("home.quick_start.company")}
          </Button>
        </Link>
        <Link to="https://useportcall.com/docs" target="_blank" className="w-full">
          <Button
            variant="outline"
            className="flex flex-col justify-center h-full w-full gap-2"
          >
            <Search />
            {t("home.quick_start.docs")}
          </Button>
        </Link>
      </div>
    </div>
  );
}
