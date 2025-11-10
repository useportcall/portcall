import { UsersTable } from "./components/users-table";

export default function UsersView() {
  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <div className="flex flex-col justify-start space-y-2">
            <p className="text-lg md:text-xl font-semibold">Users</p>
            <p className="text-sm text-slate-400">
              Manage all your users' plans & entitlements here.
            </p>
          </div>
        </div>
      </div>
      <UsersTable />
    </div>
  );
}
