import { useEffect, useState } from "react";

export function useSyncedInput(value: string | null | undefined) {
  const [currentValue, setCurrentValue] = useState(value || "");

  useEffect(() => {
    setCurrentValue(value || "");
  }, [value]);

  return { currentValue, setCurrentValue };
}
