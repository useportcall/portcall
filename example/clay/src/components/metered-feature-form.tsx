"use client";

import createMeterEvent from "@/app/actions/create-meter-event";
import { FormEventHandler, useTransition } from "react";
import { toast } from "sonner";
import { Button } from "./ui/button";
import { Input } from "./ui/input";

export function MeteredFeatureForm({ featureId }: { featureId: string }) {
  const [isPending, startTransition] = useTransition();

  const handleSubmit: FormEventHandler<HTMLFormElement> = async (event) => {
    event.preventDefault();

    const formData = new FormData(event.currentTarget);

    const usage = Number(formData.get("usage"));

    startTransition(async () => {
      const result = await createMeterEvent({
        featureId,
        usage,
      });

      if (!result.ok) {
        toast.error(result.message);
        return;
      }
    });
  };

  return (
    <form className="flex items-center gap-2" onSubmit={handleSubmit}>
      <Input
        name="usage"
        type="number"
        className="h-8"
        placeholder="100"
        defaultValue={100}
      />
      <Button size="sm" className="w-fit text-xs" disabled={isPending}>
        Submit 🔼
      </Button>
    </form>
  );
}
