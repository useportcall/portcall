import { Card } from "@/components/ui/card";
import { TablePagination } from "@/components/table/table-pagination";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useNavigate } from "react-router-dom";
import {
  UserActionsTableCell,
  UserEmailTableCell,
  UserNameTableCell,
  UserPaymentMethodCell,
  UserStatusTableCell,
} from "./user-table-cells";
import { User } from "@/models/user";

export function UsersTableCard(props: {
  users: User[];
  page: number;
  totalPages: number;
  total: number;
  limit: number;
  onPrevious: () => void;
  onNext: () => void;
  onLimitChange: (value: number) => void;
  t: (key: string, options?: Record<string, unknown>) => string;
}) {
  const navigate = useNavigate();

  return (
    <Card className="p-2 rounded-xl animate-fade-in space-y-1">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-80">{props.t("views.users.table.email")}</TableHead>
            <TableHead className="w-48">{props.t("views.users.table.name")}</TableHead>
            <TableHead className="w-40">{props.t("views.users.table.status")}</TableHead>
            <TableHead className="w-40">{props.t("views.users.table.payment")}</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {props.users.map((user) => (
            <TableRow key={user.id} onClick={() => navigate(`/users/${user.id}`)} className="rounded hover:bg-accent cursor-pointer">
              <UserEmailTableCell user={user} />
              <UserNameTableCell user={user} />
              <UserStatusTableCell user={user} />
              <UserPaymentMethodCell user={user} />
              <UserActionsTableCell user={user} />
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <TablePagination
        page={props.page}
        totalPages={props.totalPages}
        onPrevious={props.onPrevious}
        onNext={props.onNext}
        previousLabel={props.t("views.users.table.previous")}
        nextLabel={props.t("views.users.table.next")}
        rowsPerPageLabel={props.t("views.users.table.rows_per_page")}
        rowsPerPage={props.limit}
        rowsPerPageOptions={[20, 50, 100]}
        onRowsPerPageChange={props.onLimitChange}
        pageLabel={props.t("views.users.table.page", {
          page: props.page,
          totalPages: props.totalPages,
          total: props.total,
        })}
      />
    </Card>
  );
}
