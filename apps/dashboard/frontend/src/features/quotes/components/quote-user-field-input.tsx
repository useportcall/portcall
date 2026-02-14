import { useSyncedInput } from "./use-synced-input";

export function QuoteUserFieldInput({
  icon,
  value,
  placeholder,
  onSave,
}: {
  icon: React.ReactNode;
  value: string | null | undefined;
  placeholder: string;
  onSave: (value: string) => Promise<unknown>;
}) {
  const { currentValue, setCurrentValue } = useSyncedInput(value);

  return (
    <div className="flex gap-1.5 px-2.5 py-1 items-center text-sm font-medium">
      {icon}
      <input
        type="text"
        className="outline-none w-full"
        value={currentValue}
        onChange={(e) => setCurrentValue(e.target.value)}
        onBlur={async () => {
          if (currentValue !== (value || "")) {
            await onSave(currentValue);
          }
        }}
        placeholder={placeholder}
      />
    </div>
  );
}
