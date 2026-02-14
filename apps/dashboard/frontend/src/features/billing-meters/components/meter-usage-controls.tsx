import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Minus, Plus } from "lucide-react";
import { useTranslation } from "react-i18next";

export function MeterUsageControls(props: {
  usage: number;
  isEditing: boolean;
  editValue: string;
  setEditValue: (value: string) => void;
  onStartEdit: () => void;
  onCancelEdit: () => void;
  onSetUsage: () => void;
  onIncrement: () => void;
  onDecrement: () => void;
  isUpdating: boolean;
}) {
  const { t } = useTranslation();
  if (props.isEditing) {
    return (
      <div className="flex items-center gap-2">
        <Input
          type="number"
          value={props.editValue}
          onChange={(e) => props.setEditValue(e.target.value)}
          className="w-20 h-8"
          min={0}
        />
        <Button size="sm" variant="outline" onClick={props.onSetUsage}>{t("views.user.billing_meters.set")}</Button>
        <Button size="sm" variant="ghost" onClick={props.onCancelEdit}>{t("common.cancel")}</Button>
      </div>
    );
  }

  return (
    <>
      <Button
        size="sm"
        variant="outline"
        className="h-8 w-8 p-0"
        onClick={props.onDecrement}
        disabled={props.isUpdating || props.usage === 0}
      >
        <Minus className="h-3 w-3" />
      </Button>
      <Tooltip>
        <TooltipTrigger asChild>
          <button className="text-lg font-mono px-2 hover:bg-muted rounded" onClick={props.onStartEdit}>
            {props.usage.toLocaleString()}
          </button>
        </TooltipTrigger>
        <TooltipContent>{t("views.user.billing_meters.click_to_edit")}</TooltipContent>
      </Tooltip>
      <Button
        size="sm"
        variant="outline"
        className="h-8 w-8 p-0"
        onClick={props.onIncrement}
        disabled={props.isUpdating}
      >
        <Plus className="h-3 w-3" />
      </Button>
    </>
  );
}
