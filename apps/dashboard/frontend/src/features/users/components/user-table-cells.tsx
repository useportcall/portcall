import { Badge } from "@/components/ui/badge";
import { TableCell } from "@/components/ui/table";
import { User } from "@/models/user";
import { useDeleteUser } from "@/hooks";
import { CopyableIdCell } from "@/components/table/copyable-id-cell";
import { ActionsCell, ActionMenuItem } from "@/components/table/actions-cell";

export function UserEmailTableCell(props: { user: User }) {
  return <CopyableIdCell title={props.user.email} id={props.user.id} />;
}

export function UserNameTableCell(props: { user: User }) {
  return (
    <TableCell className="w-48">
      <span>{props.user.name}</span>
    </TableCell>
  );
}

export function UserStatusTableCell(props: { user: User }) {
  return (
    <TableCell className="w-40">
      <Badge variant={props.user.subscribed ? "success" : "outline"}>
        {props.user.subscribed ? "Active" : "No subscription"}
      </Badge>
    </TableCell>
  );
}

export function UserPaymentMethodCell(props: { user: User }) {
  return (
    <TableCell className="w-40">
      <Badge variant={"outline"}>
        {props.user.payment_method_added ? "Added" : "Missing"}
      </Badge>
    </TableCell>
  );
}

export function UserActionsTableCell({ user }: { user: User }) {
  const { mutate: deleteUser, isPending: isDeleting } = useDeleteUser(user.id);

  const actions: ActionMenuItem[] = [
    {
      label: "Delete",
      onClick: (e) => {
        e.stopPropagation();
        deleteUser({});
      },
      loading: isDeleting,
      disabled: isDeleting,
      className: "text-red-600 focus:text-red-700",
    },
  ];

  return <ActionsCell actions={actions} />;
}
