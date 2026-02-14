import { Input } from "./ui/input";
import ExpiryInput from "./expiry-input";
import { CreditCardInput } from "./credit-card-input";
import { MockAutofillButton } from "./mock-autofill-button";
import { SubmitState } from "./submit-state";
import { useTranslation } from "react-i18next";

type Props = {
  submitState: SubmitState;
  cardNumber: string;
  setCardNumber: (value: string) => void;
  expiry: string;
  setExpiry: (value: string) => void;
  cvc: string;
  setCvc: (value: string) => void;
  onAutofill: () => void;
};

export function MockCardFields(props: Props) {
  const { t } = useTranslation();
  const disabled = props.submitState !== "idle";
  return (
    <div className="w-full flex flex-col gap-4 mb-2 md:max-w-sm mx-auto">
      <p className="text-sm font-medium">{t("form.card_details")}</p>
      <MockAutofillButton disabled={disabled} onClick={props.onAutofill} />
      <div className="flex flex-col gap-2">
        <CreditCardInput
          disabled={disabled}
          value={props.cardNumber}
          onValueChange={props.setCardNumber}
        />
        <span className="flex w-full justify-between gap-2">
          <ExpiryInput
            disabled={disabled}
            value={props.expiry}
            onValueChange={props.setExpiry}
          />
          <Input
            data-testid="checkout-cvc-input"
            placeholder={t("form.cvc")}
            className="bg-white text-sm"
            maxLength={4}
            value={props.cvc}
            onChange={(event) => props.setCvc(event.target.value.replace(/\D/g, "").slice(0, 4))}
            disabled={disabled}
          />
        </span>
      </div>
    </div>
  );
}
