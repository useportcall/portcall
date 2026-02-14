import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Connection } from "@/hooks/api/connections";
import { MoreVertical } from "lucide-react";
import { DeleteConnectionItem } from "./delete-connection-item";
import { MakeDefaultItem } from "./make-default-item";

export function ConnectionActionsMenu({
  connection,
  isDefault,
}: {
  connection: Connection;
  isDefault: boolean;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          size="icon"
          variant="ghost"
          className="h-8 w-8 text-muted-foreground hover:text-foreground"
        >
          <MoreVertical className="w-4 h-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48">
        <MakeDefaultItem id={connection.id} checked={isDefault} />
        <DropdownMenuSeparator />
        <DeleteConnectionItem id={connection.id} />
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
