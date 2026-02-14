import { DropdownMenuItem } from "@/components/ui/dropdown-menu";
import { useAppConfig } from "@/hooks/api/app-configs";
import { Check, Loader2, Star } from "lucide-react";

export function MakeDefaultItem({ id, checked }: { id: string; checked: boolean }) {
  const { mutateAsync, isPending } = useAppConfig();

  if (checked) {
    return (
      <DropdownMenuItem disabled className="text-muted-foreground">
        <Check className="w-4 h-4 mr-2" />
        Current default
      </DropdownMenuItem>
    );
  }

  return (
    <DropdownMenuItem disabled={isPending} onClick={() => mutateAsync({ connection_id: id })}>
      {isPending ? (
        <span className="flex items-center gap-2"><Loader2 className="animate-spin w-4 h-4" /> Setting...</span>
      ) : (
        <span className="flex items-center gap-2"><Star className="w-4 h-4" /> Set as default</span>
      )}
    </DropdownMenuItem>
  );
}
