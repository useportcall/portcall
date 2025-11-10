import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Separator } from "./ui/separator";

export function BaseView({
  title,
  description,
  children,
  actions,
}: {
  title: string;
  description: string;
  children: React.ReactNode;
  actions?: React.ReactNode;
}) {
  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8">
          <div className="flex flex-col justify-start space-y-2">
            <p className="text-lg md:text-xl font-semibold">{title}</p>
            <p className="text-sm text-slate-400">{description}</p>
          </div>
        </div>
        <div>{actions}</div>
      </div>
      <div className="flex flex-col gap-4">{children}</div>
    </div>
  );
}

export function BaseSection({
  title,
  children,
  actions,
}: {
  title: string;
  children: React.ReactNode;
  actions?: React.ReactNode;
}) {
  return (
    <>
      <Separator />
      <div className="flex justify-between">
        <h4 className="text-slate-400 text-sm">{title}</h4>
        {actions}
      </div>
      <div className="flex flex-col gap-4">{children}</div>
    </>
  );
}

export function LabeledInput({
  id,
  label,
  helperText,
  ...props
}: {
  id: string;
  label: string;
  helperText: string;
} & React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <div className="w-full flex flex-col gap-2 justify-start items-start">
      <Label htmlFor={id}>{label}</Label>
      <p className="text-xs text-slate-400">{helperText}</p>
      <Input id={id} {...props} />
    </div>
  );
}
