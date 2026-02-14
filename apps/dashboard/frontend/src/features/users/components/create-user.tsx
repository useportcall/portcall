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
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router-dom";
import z from "zod";

export const schema = z.object({
  email: z.string().email("Invalid email address"),
});

export function CreateUser() {
  const { t } = useTranslation();
  const navigate = useNavigate();
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
        <Button data-testid="add-user-button" size={"sm"} variant={"outline"}>
          <Plus />
          {t("views.users.create.action")}
        </Button>
      </DialogTrigger>
      <DialogContent className="flex flex-col gap-4">
        <DialogHeader>
          <DialogTitle>{t("views.users.create.title")}</DialogTitle>
          <DialogDescription>
            {t("views.users.create.description")}
          </DialogDescription>
        </DialogHeader>
        <form
          onSubmit={form.handleSubmit(async (fields) => {
            const result = await mutateAsync({ email: fields.email });
            // Navigate immediately without closing the dialog first
            navigate(`/users/${result.data.id}`);
          })}
          className="w-full flex flex-col gap-4"
        >
          <BaseSection title="">
            <LabeledInput
              data-testid="create-user-email-input"
              autoFocus
              id="email"
              label={t("views.users.create.email")}
              placeholder={t("views.users.create.email_placeholder")}
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
              {t("common.cancel")}
            </Button>
            <Button
              data-testid="create-user-submit"
              size="sm"
              disabled={isPending}
            >
              {t("views.users.create.action")}
              {isPending && <Loader2 className="size-4 animate-spin" />}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
