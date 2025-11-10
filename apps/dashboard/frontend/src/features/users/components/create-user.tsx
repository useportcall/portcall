import { BaseSection, LabeledInput } from "@/components/base-view";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { useCreateUser } from "@/hooks";
import { zodResolver } from "@hookform/resolvers/zod";

import { Loader2, Plus } from "lucide-react";
import { useState } from "react";
import { useForm } from "react-hook-form";
import z from "zod";

export const schema = z.object({
  email: z.string().email("Invalid email address"),
});

export function CreateUser() {
  const { mutateAsync, isPending } = useCreateUser();

  const [showForm, setShowForm] = useState(false);
  const form = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      email: "",
    },
  });

  return (
    <Dialog open={showForm} onOpenChange={setShowForm}>
      <DialogTrigger asChild>
        <Button size={"sm"} variant={"outline"}>
          <Plus />
          Add user
        </Button>
      </DialogTrigger>
      <DialogContent className="flex flex-col gap-4">
        <DialogHeader>
          <DialogTitle>Add new user</DialogTitle>
          <DialogDescription>
            Connect an existing user using their email address.
          </DialogDescription>
        </DialogHeader>
        <form
          onSubmit={form.handleSubmit(async (fields) => {
            const result = await mutateAsync({ email: fields.email });

            window.location.href = `/users/${result.data.id}`;
          })}
          className="w-full flex flex-col gap-4"
        >
          <BaseSection title="">
            <LabeledInput
              autoFocus
              id="email"
              label="Email"
              placeholder="Add email"
              className="font-semibold outline-none w-full"
              helperText=""
              {...form.register("email")}
            />
          </BaseSection>
          <div className="w-fit mt-4 flex flex-row gap-2 justify-end self-end">
            <Button
              type="button"
              variant="ghost"
              onClick={() => setShowForm(false)}
            >
              Cancel
            </Button>
            <Button size="sm" disabled={isPending}>
              Add user
              {isPending && <Loader2 className="size-4 animate-spin" />}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
