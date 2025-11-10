import EmptyTable from "@/components/empty-table";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useListUsers } from "@/hooks";
import { cn } from "@/lib/utils";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { CreateUser } from "./create-user";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";

export function UsersTable() {
  const [search, setSearch] = useState("");

  const { data: users } = useListUsers({ email: search });

  const navigate = useNavigate();

  return (
    <>
      <div className="flex justify-between">
        <SearchOrAddUser
          key="searchbox"
          // count={users?.data.length ?? 0}
          value={search}
          onChange={setSearch}
        />
        <CreateUser />
      </div>

      {!!users?.data.length && (
        <Card className="p-2 rounded-xl animate-fade-in">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-80">Email</TableHead>
                <TableHead className="w-48">Name</TableHead>
                <TableHead className="w-40">Status</TableHead>
                {/* <TableHead className="w-40">Billing</TableHead>
              <TableHead className="w-40">Payments</TableHead> */}
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody className="">
              {users?.data.map((user) => (
                <TableRow key={user.id} className="">
                  <TableCell className="w-20 py-1">
                    <div>
                      <h4 className="font-semibold overflow-ellipsis">
                        {user.email}
                      </h4>
                      <span className="text-xs text-slate-500 italic">
                        {user.id}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="py-1">
                    <span>{user.name}</span>
                  </TableCell>
                  <TableCell className="py-1">
                    <Badge
                      variant={"outline"}
                      className={cn(
                        user.subscribed
                          ? "bg-emerald-400/50 text-emerald-800 border-none"
                          : ""
                      )}
                    >
                      {user.subscribed ? "Active" : "No subscription"}
                    </Badge>
                  </TableCell>
                  <TableCell align="right" className="right py-1">
                    <Button
                      onClick={() => navigate(`/users/${user.id}`)}
                      className="text-xs px-3 py-1"
                      variant="outline"
                    >
                      Manage
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}
      {!users?.data.length && (
        <EmptyTable message="No users added yet." button={""} />
      )}
    </>
  );
}

function SearchOrAddUser({
  // count,
  value,
  onChange,
}: {
  // count: number;
  value: string;
  onChange: (v: string) => void;
}) {
  return (
    <div className="relative flex items-center gap-2 w-full max-w-xs">
      <Input
        type="text"
        placeholder="Find user by email..."
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className=""
      />
    </div>
  );
}
