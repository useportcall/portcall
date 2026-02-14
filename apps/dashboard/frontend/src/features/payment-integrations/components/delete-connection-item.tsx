import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { useDeleteConnection } from "@/hooks/api/connections";
import { Loader2, Trash2 } from "lucide-react";

export function DeleteConnectionItem({ id }: { id: string }) {
  const { mutate, isPending } = useDeleteConnection(id);

  return (
    <DropdownMenuItem
      disabled={isPending}
      onClick={() => mutate({})}
      className="text-destructive focus:text-destructive"
    >
      {isPending ? (
        <span className="flex items-center gap-2"><Loader2 className="animate-spin w-4 h-4" /> Deleting...</span>
      ) : (
        <span className="flex items-center gap-2"><Trash2 className="w-4 h-4" /> Delete</span>
      )}
    </DropdownMenuItem>
  );
}
