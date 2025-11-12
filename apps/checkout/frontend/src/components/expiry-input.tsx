"use client"; // <-- remove if you are not using Next.js App Router

import React, { useState } from "react";
import { Input } from "./ui/input";

interface ExpiryInputProps {
  /** Show full year (MM/YYYY) instead of MM/YY */
  fullYear?: boolean;
  /** Optional label */
  label?: string;
  /** Optional name for form submission */
  name?: string;
  /** Optional className */
  className?: string;
}

/**
 * Credit-card expiration date input
 * Formats as MM/YY or MM/YYYY
 */
export default function ExpiryInput({
  fullYear = false,
  label = "Expiration date",
  name = "expiry",
  className = "",
}: ExpiryInputProps) {
  const [value, setValue] = useState("");

  const formatExpiry = (raw: string): string => {
    // 1. Keep only digits
    const digits = raw.replace(/\D/g, "");

    // 2. Max length depends on fullYear flag
    const maxDigits = fullYear ? 6 : 4; // MMYY vs MMYYYY
    const limited = digits.slice(0, maxDigits);

    // 3. Split into month / year
    const month = limited.slice(0, 2);
    const year = limited.slice(2);

    // 4. Validate month (01-12)
    let formatted = month;
    if (month.length === 2) {
      const m = Number(month);
      if (m < 1) formatted = "01";
      else if (m > 12) formatted = "12";
    }

    // 5. Add slash + year when month is complete
    if (limited.length > 2) {
      formatted += "/" + (fullYear ? year.padEnd(4, "0") : year);
    }

    return formatted;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const formatted = formatExpiry(e.target.value);
    setValue(formatted);
  };

  // Optional: prevent cursor jump when inserting slash
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    const input = e.currentTarget;
    const cursor = input.selectionStart || 0;

    // If user types a digit after the slash, keep cursor after slash
    if (e.key >= "0" && e.key <= "9" && cursor === 3 && value.length === 2) {
      // let the change happen first, then move cursor
      setTimeout(() => input.setSelectionRange(4, 4), 0);
    }
  };

  const maxLen = fullYear ? 7 : 5; // MM/YY → 5 chars, MM/YYYY → 7 chars

  return (
    <Input
      id={name}
      name={name}
      type="text"
      inputMode="numeric"
      autoComplete="cc-exp"
      placeholder={fullYear ? "MM/YYYY" : "MM/YY"}
      value={value}
      onChange={handleChange}
      onKeyDown={handleKeyDown}
      maxLength={maxLen}
      className="bg-white text-sm"
    />
  );
}
