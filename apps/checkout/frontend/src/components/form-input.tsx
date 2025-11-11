import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { ReactNode } from "react";
import { Control, FieldValues, Path } from "react-hook-form";

export function FormInput<T extends FieldValues>(props: {
  control: Control<T>;
  name: Path<T>;
  title: string;
  description?: string;
  placeholder?: string;
  disabled?: boolean;
  hidden?: boolean;
  suffix?: ReactNode;
  type?: "text" | "datetime-local";
  className?: string;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
}) {
  return (
    <FormField
      control={props.control}
      name={props.name}
      render={({ field }) => (
        <FormItem
          className={cn(
            "w-full",
            props.hidden ? "hidden" : "",
            props.className
          )}
        >
          <div className="flex items-center justify-between">
            <FormLabel className="text-sm font-medium text-slate-800">
              {props.title}
            </FormLabel>
            <FormMessage className="text-xs" />
          </div>
          <FormControl className="relative">
            <div>
              {!!props.suffix && props.suffix}
              <Input
                type={props.type ?? "text"}
                placeholder={props.placeholder}
                disabled={props.disabled}
                onChange={field.onChange}
                onBlur={props.onBlur}
                defaultValue={
                  props.type === "datetime-local" && field.value
                    ? new Date(field.value).toISOString().slice(0, 16)
                    : field.value
                }
                className="text-sm text-slate-800 bg-white"
              />
            </div>
          </FormControl>
        </FormItem>
      )}
    />
  );
}
