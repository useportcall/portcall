import { Button } from "./ui/button";
import { useTranslation } from "react-i18next";

type Props = {
  disabled: boolean;
  onClick: () => void;
};

export function MockAutofillButton({ disabled, onClick }: Props) {
  const { t } = useTranslation();
  return (
    <Button
      type="button"
      variant="outline"
      className="w-full"
      data-testid="checkout-autofill-button"
      disabled={disabled}
      onClick={onClick}
    >
      {t("test_mode.autofill_button")}
    </Button>
  );
}
