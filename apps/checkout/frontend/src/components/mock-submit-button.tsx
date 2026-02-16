import { cn } from "@/lib/utils";
import { Check, Loader2, Lock } from "lucide-react";
import { MouseEvent } from "react";

export function MockSubmitButton(props: {
  total: string;
  state: "idle" | "loading" | "success";
  consentMode: "save" | "charge";
  disabled?: boolean;
  onTrySubmit?: (event: MouseEvent<HTMLButtonElement>) => void;
  t: (key: string, options?: Record<string, unknown>) => string;
}) {
  return (
    <button
      data-testid="checkout-submit-button"
      type="submit"
      disabled={props.state !== "idle" || props.disabled}
      onClick={props.onTrySubmit}
      className={cn(
        "w-full p-3 flex justify-center items-center gap-2 text-white font-semibold mt-2 rounded-md transition-all duration-300",
        props.state === "idle" && "bg-slate-900 hover:bg-slate-700 cursor-pointer",
        props.state === "loading" && "bg-slate-700 cursor-wait",
        props.state === "success" && "bg-emerald-500 cursor-default",
        props.disabled && "opacity-50 cursor-not-allowed",
      )}
    >
      {props.state === "idle" ? (
        <>
          <span>
            {props.consentMode === "save"
              ? props.t("submit.authorize_and_continue")
              : props.t("submit.pay", { total: props.total })}
          </span>
          <Lock className="w-4 h-4" />
        </>
      ) : null}
      {props.state === "loading" ? (
        <>
          <Loader2 className="w-5 h-5 animate-spin" />
          <span>{props.t("submit.processing")}</span>
        </>
      ) : null}
      {props.state === "success" ? (
        <>
          <Check className="w-5 h-5" />
          <span>{props.t("submit.payment_successful")}</span>
        </>
      ) : null}
    </button>
  );
}
