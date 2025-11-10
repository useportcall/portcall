import { useUpdatePlan } from "@/hooks";
import { cn } from "@/lib/utils";
import { Plan } from "@/models/plan";

export function PlanNameInput({ plan }: { plan: Plan }) {
  const { mutateAsync } = useUpdatePlan(plan.id);

  return (
    <input
      defaultValue={plan.name}
      placeholder="Plan Name"
      className={cn(
        "font-semibold text-2xl outline-none w-full",
        plan.name.length ? "text-black" : "text-muted-foreground"
      )}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          mutateAsync({ name: e.target.value });
        }
      }}
    />
  );
}
