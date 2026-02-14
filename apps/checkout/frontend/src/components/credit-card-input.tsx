import { useState } from "react";
import { Input } from "./ui/input";

type Props = {
  disabled?: boolean;
  value?: string;
  onValueChange?: (value: string) => void;
};

export function CreditCardInput({ disabled, value, onValueChange }: Props) {
  const [internalValue, setInternalValue] = useState("");
  const cardNumber = value ?? internalValue;

  const isAmex =
    cardNumber.replace(/\s+/g, "").startsWith("34") ||
    cardNumber.replace(/\s+/g, "").startsWith("37");

  const formatCardNumber = (value: string) => {
    const digitsOnly = value.replace(/\D/g, "").slice(0, 16);

    if (isAmex) {
      const truncated = digitsOnly.slice(0, 15);
      return truncated.replace(/^(\d{4})(\d{6})(\d{5})$/, "$1 $2 $3");
    }

    const truncated = digitsOnly.slice(0, 16);
    return truncated.match(/.{1,4}/g)?.join(" ") || truncated;
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const nextValue = formatCardNumber(e.target.value);
    if (onValueChange) {
      onValueChange(nextValue);
      return;
    }
    setInternalValue(nextValue);
  };

  return (
    <Input
      data-testid="checkout-card-input"
      type="text"
      value={cardNumber}
      onChange={handleChange}
      placeholder="1234 5678 9012 3456"
      maxLength={isAmex ? 17 : 19}
      className="bg-white text-sm"
      disabled={disabled}
    />
  );
}
