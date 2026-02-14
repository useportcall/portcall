import { CountryDropdown } from "@/components/country-dropdown";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { cn } from "@/lib/utils";
import { Control, FieldValues, Path } from "react-hook-form";
import { useTranslation } from "react-i18next";

export function CountryInput<T extends FieldValues>(props: {
  control: Control<T>;
  name: Path<T>;
  disabled?: boolean;
  hidden?: boolean;
}) {
  const { t } = useTranslation();

  return (
    <FormField
      control={props.control}
      name={props.name}
      render={({ field }) => (
        <FormItem className={cn("w-full", props.hidden ? "hidden" : "")}>
          <div className="flex items-center justify-between">
            <FormLabel className="text-sm font-medium text-slate-800">
              {t("form.country_territory")}
            </FormLabel>
            <FormMessage className="text-xs" />
          </div>
          <FormControl>
            <CountryDropdown
              placeholder={t("form.select_country")}
              defaultValue={field.value}
              onChange={(c) => field.onChange(c.alpha2)}
            />
          </FormControl>
        </FormItem>
      )}
    />
  );
}
