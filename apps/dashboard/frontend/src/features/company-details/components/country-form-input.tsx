import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { cn } from "@/lib/utils";
import { Control, FieldValues, Path } from "react-hook-form";
import { CountryDropdown } from "../../../components/country-dropdown";

export default function CountryFormInput<T extends FieldValues>(props: {
  control: Control<T, any>;
  name: Path<T>;
  title: string;
  description: string;
  placeholder?: string;
  disabled?: boolean;
  hidden?: boolean;
}) {
  return (
    <FormField
      control={props.control}
      name={props.name}
      render={({ field }) => (
        <FormItem
          className={cn(
            "space-y-2 flex flex-col",
            props.hidden ? "hidden" : ""
          )}
        >
          <div className="space-y-2 flex flex-col">
            <FormLabel className="text-sm font-medium text-slate-800">
              {props.title}
            </FormLabel>
            <FormDescription className="text-xs text-slate-400">
              {props.description}
            </FormDescription>
          </div>
          <FormControl>
            <CountryDropdown
              placeholder={props.placeholder}
              disabled={props.disabled}
              onChange={(value) => {
                field.onChange(value.alpha2);
              }}
              defaultValue={field.value}
            />
          </FormControl>

          <FormMessage />
        </FormItem>
      )}
    />
  );
}
