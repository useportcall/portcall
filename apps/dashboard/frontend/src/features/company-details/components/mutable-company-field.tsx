import { LabeledInput } from "@/components/base-view";
import { InputSaveIndicator, useInputSaveIndicator } from "@/components/input-save-indicator";
import { useUpdateCompany } from "@/hooks";
import { cn } from "@/lib/utils";

export function MutableCompanyField({ id, label, defaultValue, placeholder, payloadKey }: { id: string; label: string; defaultValue: string; placeholder: string; payloadKey: string }) {
  const { mutateAsync: updateCompany } = useUpdateCompany();
  const { visible, showSaved } = useInputSaveIndicator();

  return (
    <div className="flex flex-row gap-2"><LabeledInput id={id} helperText="" label={label} defaultValue={defaultValue || ""} placeholder={placeholder} className={cn("text-lg outline-none w-96")} onBlur={async (e) => { if (e.target.value.length > 0) { const formData = new FormData(); formData.append("data", JSON.stringify({ [payloadKey]: e.target.value })); await updateCompany(formData); showSaved(); } }} /><InputSaveIndicator visible={visible} /></div>
  );
}
