import React, { useState } from "react";
import { Input } from "./ui/input";

type ExpiryInputProps = {
  fullYear?: boolean;
  name?: string;
  disabled?: boolean;
  value?: string;
  onValueChange?: (value: string) => void;
};

export default function ExpiryInput({
  fullYear = false,
  name = "expiry",
  disabled,
  value,
  onValueChange,
}: ExpiryInputProps) {
  const [internalValue, setInternalValue] = useState("");
  const inputValue = value ?? internalValue;

  const formatExpiry = (raw: string): string => {
    const digits = raw.replace(/\D/g, "");
    const maxDigits = fullYear ? 6 : 4;
    const limited = digits.slice(0, maxDigits);
    const month = limited.slice(0, 2);
    const year = limited.slice(2);
    let formatted = month;
    if (month.length === 2) {
      const m = Number(month);
      if (m < 1) formatted = "01";
      else if (m > 12) formatted = "12";
    }
    if (limited.length > 2) {
      formatted += "/" + (fullYear ? year.padEnd(4, "0") : year);
    }
    return formatted;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const formatted = formatExpiry(e.target.value);
    if (onValueChange) {
      onValueChange(formatted);
      return;
    }
    setInternalValue(formatted);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    const input = e.currentTarget;
    const cursor = input.selectionStart || 0;
    if (e.key >= "0" && e.key <= "9" && cursor === 3 && inputValue.length === 2) {
      setTimeout(() => input.setSelectionRange(4, 4), 0);
    }
  };

  const maxLen = fullYear ? 7 : 5;

  return (
    <Input
      data-testid="checkout-expiry-input"
      id={name}
      name={name}
      type="text"
      inputMode="numeric"
      autoComplete="cc-exp"
      placeholder={fullYear ? "MM/YYYY" : "MM/YY"}
      value={inputValue}
      onChange={handleChange}
      onKeyDown={handleKeyDown}
      maxLength={maxLen}
      className="bg-white text-sm"
      disabled={disabled}
    />
  );
}
