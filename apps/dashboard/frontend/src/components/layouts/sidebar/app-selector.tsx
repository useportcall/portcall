import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useCreateApp, useListApps } from "@/hooks";
import { useApp } from "@/hooks/use-app";
import { Plus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

export function AppSelector({ iconTrigger, createOptionOnlyWhenEmpty }: { iconTrigger: boolean; createOptionOnlyWhenEmpty: boolean }) {
  const { t } = useTranslation();
  const { data: apps } = useListApps();
  const app = useApp();
  const { mutate, isPending } = useCreateApp();
  const [name, setName] = useState("");
  const [open, setOpen] = useState(false);
  const showCreateOption = !createOptionOnlyWhenEmpty || apps.data.length === 0;
  const selectedApp = apps.data.find((appItem) => appItem.id.toString() === app.id);
  const collapsedBadge = selectedApp
    ? { label: selectedApp.is_live ? "L" : "T", title: selectedApp.is_live ? t("app_selector.live") : t("app_selector.test"), variant: selectedApp.is_live ? "default" as const : "secondary" as const }
    : { label: "?", title: t("app_selector.placeholder"), variant: "outline" as const };
  const collapsedTriggerClass = collapsedBadge.variant === "default"
    ? "bg-primary text-primary-foreground border-primary hover:bg-primary/90"
    : collapsedBadge.variant === "secondary"
      ? "bg-muted text-foreground border-border hover:bg-muted/90"
      : "bg-background text-foreground border-input hover:bg-accent";

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <Select value={app.id || undefined} onValueChange={(value) => (value === "__create_new__" ? setOpen(true) : app.setId(value))}>
        <SelectTrigger
          className={iconTrigger ? `my-2 h-9 w-9 !p-0 !justify-center mx-auto [&>svg]:hidden ${collapsedTriggerClass}` : "my-2"}
          aria-label={t("app_selector.aria_label")}
          title={iconTrigger ? collapsedBadge.title : undefined}
        >
          {iconTrigger ? (
            <div className="h-full w-full grid place-items-center text-sm font-semibold pointer-events-none">{collapsedBadge.label}</div>
          ) : (
            <SelectValue placeholder={t("app_selector.placeholder")} />
          )}
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectLabel className="text-xs text-muted-foreground font-medium">{t("app_selector.apps_label")}</SelectLabel>
            {apps.data.map((appItem) => (
              <SelectItem key={appItem.id} value={appItem.id.toString()}>
                <div className="flex items-center gap-2">
                  <span>{appItem.name}</span>
                  <Badge variant={appItem.is_live ? "default" : "secondary"} className="h-5 px-1.5 text-[10px] pointer-events-none">
                    {appItem.is_live ? t("app_selector.live") : t("app_selector.test")}
                  </Badge>
                </div>
              </SelectItem>
            ))}
          </SelectGroup>
          {showCreateOption && <CreateNewAppOption />}
        </SelectContent>
      </Select>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t("app_selector.create_title")}</DialogTitle>
          <DialogDescription>{t("app_selector.create_description")}</DialogDescription>
        </DialogHeader>
        <Input placeholder={t("app_selector.create_placeholder")} value={name} onChange={(e) => setName(e.target.value)} />
        <DialogFooter>
          <Button
            disabled={isPending || !name.trim()}
            onClick={() => mutate({ name }, { onSuccess: (newApp) => { setOpen(false); setName(""); app.setId(newApp.id); } })}
          >
            {isPending ? t("app_selector.creating") : t("app_selector.create_action")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function CreateNewAppOption() {
  const { t } = useTranslation();
  return (
    <>
      <Separator className="my-1" />
      <SelectItem value="__create_new__" className="cursor-pointer">
        <div className="flex items-center gap-2"><Plus className="h-4 w-4" /><span>{t("app_selector.create_option")}</span></div>
      </SelectItem>
    </>
  );
}
