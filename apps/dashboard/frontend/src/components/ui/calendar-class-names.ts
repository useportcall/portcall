import { buttonVariants, Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { DayPicker, getDefaultClassNames } from "react-day-picker";
import type React from "react";

export function buildCalendarClassNames({
  classNames,
  buttonVariant,
  captionLayout,
}: {
  classNames?: React.ComponentProps<typeof DayPicker>["classNames"];
  buttonVariant: React.ComponentProps<typeof Button>["variant"];
  captionLayout: React.ComponentProps<typeof DayPicker>["captionLayout"];
}) {
  const d = getDefaultClassNames();
  return {
    root: cn("w-fit", d.root), months: cn("relative flex flex-col gap-4 md:flex-row", d.months), month: cn("flex w-full flex-col gap-4", d.month), nav: cn("absolute inset-x-0 top-0 flex w-full items-center justify-between gap-1", d.nav),
    button_previous: cn(buttonVariants({ variant: buttonVariant, size: "icon" }), "h-[--cell-size] w-[--cell-size] select-none p-0 aria-disabled:opacity-50", d.button_previous),
    button_next: cn(buttonVariants({ variant: buttonVariant, size: "icon" }), "h-[--cell-size] w-[--cell-size] select-none p-0 aria-disabled:opacity-50", d.button_next),
    month_caption: cn("flex h-[--cell-size] w-full items-center justify-center px-[--cell-size]", d.month_caption), dropdowns: cn("flex h-[--cell-size] w-full items-center justify-center gap-1.5 text-sm font-medium", d.dropdowns),
    dropdown_root: cn("has-focus:border-ring border-input shadow-xs has-focus:ring-ring/50 has-focus:ring-[3px] relative rounded-md border", d.dropdown_root), dropdown: cn("bg-popover absolute inset-0 opacity-0", d.dropdown),
    caption_label: cn("select-none font-medium", captionLayout === "label" ? "text-sm" : "[&>svg]:text-muted-foreground flex h-8 items-center gap-1 rounded-md pl-2 pr-1 text-sm [&>svg]:size-3.5", d.caption_label),
    table: "w-full border-collapse", weekdays: cn("flex", d.weekdays), weekday: cn("text-muted-foreground flex-1 select-none rounded-md text-[0.8rem] font-normal", d.weekday), week: cn("mt-2 flex w-full", d.week),
    week_number_header: cn("w-[--cell-size] select-none", d.week_number_header), week_number: cn("text-muted-foreground select-none text-[0.8rem]", d.week_number), day: cn("group/day relative aspect-square h-full w-full select-none p-0 text-center [&:first-child[data-selected=true]_button]:rounded-l-md [&:last-child[data-selected=true]_button]:rounded-r-md", d.day),
    range_start: cn("bg-accent rounded-l-md", d.range_start), range_middle: cn("rounded-none", d.range_middle), range_end: cn("bg-accent rounded-r-md", d.range_end),
    today: cn("bg-accent text-accent-foreground rounded-md data-[selected=true]:rounded-none", d.today), outside: cn("text-muted-foreground aria-selected:text-muted-foreground", d.outside), disabled: cn("text-muted-foreground opacity-50", d.disabled), hidden: cn("invisible", d.hidden),
    ...classNames,
  };
}
