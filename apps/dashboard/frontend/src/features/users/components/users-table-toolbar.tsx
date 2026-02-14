import { Input } from "@/components/ui/input";

export function UsersTableToolbar(props: {
  value: string;
  onChange: (v: string) => void;
  placeholder: string;
}) {
  return (
    <div className="relative flex items-center gap-2 w-full max-w-xs">
      <Input
        type="text"
        placeholder={props.placeholder}
        value={props.value}
        onChange={(e) => props.onChange(e.target.value)}
      />
    </div>
  );
}
