import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { useGetAccount } from "@/hooks";
import { Boxes, BriefcaseBusiness, Search, User } from "lucide-react";
import { Link } from "react-router";

export default function HomeView() {
  const { data: account } = useGetAccount();

  if (!account?.data) return <></>;

  return (
    <>
      <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
        <div className="flex flex-col gap-4 lg:flex-row justify-between items-start">
          <div className="flex flex-col space-y-2 justify-start">
            <h1 className="text-xl md:text-2xl font-bold">
              Hey {account.data.first_name}!
            </h1>
            <p className="text-slate-400 text-sm md:text-base">
              Manage plans, user entitlements, and more using this dashboard.
            </p>
          </div>
        </div>
        <Separator />
        <div className="w-full flex flex-col gap-4 justify-center items-center h-full">
          <p className="text-lg font-medium">Quick start panel</p>
          <div className="h-fit w-full lg:max-w-lg justify-center items-center my-auto mx-auto grid grid-cols-2 gap-4 bg-slate-50 p-4 md:p-10 rounded-md">
            <Link to={"/plans"} className="w-full">
              <Button
                variant={"outline"}
                className="flex flex-col w-full justify-center h-full gap-2"
              >
                <Boxes />
                Add or manage plans
              </Button>
            </Link>
            <Link to={"/users"} className="w-full">
              <Button
                variant={"outline"}
                className="flex flex-col justify-center h-full w-full gap-2"
              >
                <User />
                Add or manage users
              </Button>
            </Link>
            <Link to={"/company"} className="w-full">
              <Button
                variant={"outline"}
                className="flex flex-col justify-center h-full w-full gap-2"
              >
                <BriefcaseBusiness />
                Add your company details
              </Button>
            </Link>
            <Link
              to={"https://useportcall.com/docs"}
              target="_blank"
              className="w-full"
            >
              <Button
                variant={"outline"}
                className="flex flex-col justify-center h-full w-full gap-2"
              >
                <Search />
                Explore our docs
              </Button>
            </Link>
          </div>
        </div>
      </div>
    </>
  );
}
