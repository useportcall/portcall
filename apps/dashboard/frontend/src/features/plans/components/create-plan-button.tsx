import { Button } from "@/components/ui/button";
import { useCreatePlan } from "@/hooks";
import { Loader2, Plus } from "lucide-react";

export function CreatePlanButton() {
  const { mutate, isPending } = useCreatePlan();

  return (
    <Button
      data-testid="create-plan-button"
      size="sm"
      variant="outline"
      onClick={() => {
        mutate({});
      }}
    >
      {isPending && <Loader2 className="size-4 animate-spin" />}
      {!isPending && <Plus className="size-4" />}
      Add plan
    </Button>
  );
}
