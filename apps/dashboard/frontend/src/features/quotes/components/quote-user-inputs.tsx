import { useUpdateQuote } from "@/hooks/api/quotes";
import { useUpdateUser } from "@/hooks";
import { Quote } from "@/models/quote";
import { User } from "@/models/user";
import { Briefcase, Building2, Signature } from "lucide-react";
import { QuoteUserFieldInput } from "./quote-user-field-input";

export function QuoteUserNameInput({ user }: { user: User }) {
  const { mutateAsync } = useUpdateUser(user.id);
  return (
    <QuoteUserFieldInput
      icon={<Signature className="size-4" />}
      value={user.name}
      placeholder="Enter user's name"
      onSave={(name) => mutateAsync({ name })}
    />
  );
}

export function QuoteUserTitleInput({ user }: { user: User }) {
  const { mutateAsync } = useUpdateUser(user.id);
  return (
    <QuoteUserFieldInput
      icon={<Briefcase className="size-4" />}
      value={user.company_title}
      placeholder="Enter user's company title"
      onSave={(company_title) => mutateAsync({ company_title })}
    />
  );
}

export function QuoteUserCompanyNameInput({ quote }: { quote: Quote }) {
  const { mutateAsync } = useUpdateQuote(quote.id);
  return (
    <QuoteUserFieldInput
      icon={<Building2 className="size-4" />}
      value={quote.company_name}
      placeholder="Enter company name"
      onSave={(company_name) => mutateAsync({ company_name })}
    />
  );
}
