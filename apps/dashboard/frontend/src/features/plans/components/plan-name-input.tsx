import { useRetrievePlan, useUpdatePlan } from "@/hooks";
import { cn } from "@/lib/utils";

export function PlanNameInput(props: { id: string }) {
  const { data: plan } = useRetrievePlan(props.id);
  const { mutateAsync } = useUpdatePlan(props.id);
  if (!plan?.data) return <input placeholder="Plan Name" className="font-semibold text-2xl outline-none w-full text-muted-foreground" disabled />;

  return (
    <input
      data-testid="plan-name-input"
      defaultValue={plan.data.name}
      placeholder="Plan Name"
      className={cn(
        "font-semibold text-2xl outline-none w-full",
        plan.data.name.length ? "text-black" : "text-muted-foreground",
      )}
      onBlur={(e) => {
        if (e.target.value.length > 0) {
          mutateAsync({ name: e.target.value });
        }
      }}
    />
  );
}
