import { SaveIndicator } from "@/components/save-indicator";
import { Section } from "@/components/sections";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Separator } from "@/components/ui/separator";
import { useCreateUser, useListUsers, useUpdateUser } from "@/hooks";
import { useGetQuote, useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { User } from "@/models/user";
import { PopoverPortal } from "@radix-ui/react-popover";
import { format } from "date-fns";
import {
  ArrowLeft,
  Banknote,
  BanknoteX,
  Briefcase,
  Building2,
  CalendarOff,
  Signature,
  UserIcon,
} from "lucide-react";
import { startTransition, Suspense, useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { BillInAdvanceSelect } from "../plans/components/bill-in-advance-select";
import { FreeTrialComboBox } from "../plans/components/free-trial-combo-box";
import MutableEntitlements from "../plans/components/mutable-entitlements";
import { PlanCurrencySelect } from "../plans/components/plan-currency-select";
import { PlanFixedUnitAmountInput } from "../plans/components/plan-fixed-unit-amount-input";
import { PlanGroupSelect } from "../plans/components/plan-group-select";
import { PlanIntervalSelect } from "../plans/components/plan-interval-select";
import { PlanItems } from "../plans/components/plan-items";
import { PlanNameInput } from "../plans/components/plan-name-input";
import { ProrationSelect } from "../plans/components/proration-select";
import {
  PreviewQuoteButton,
  SendQuoteButton,
  VoidQuoteButton,
} from "./components/buttons";

export function QuoteView() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { data: quote } = useGetQuote(id!);

  return (
    <div className="w-full h-full flex flex-col gap-4 p-4 overflow-auto">
      <Button
        variant="ghost"
        className="mb-4 p-0 self-start text-sm flex flex-row justify-start items-center gap-2"
        onClick={() => navigate("/quotes")}
      >
        <ArrowLeft className="w-4 h-4" /> Back to quotes
      </Button>
      <div className="flex gap-4 w-full justify-between px-4">
        <AnimatedSuspense>
          <PlanNameInput id={quote.data.plan.id} />
        </AnimatedSuspense>
        <div className="flex flex-row justify-end gap-2">
          <SaveIndicator />
          <VoidQuoteButton id={quote.data.id} />
          <PreviewQuoteButton id={quote.data.id} />
          <SendQuoteButton id={quote.data.id} />
        </div>
      </div>
      <Separator />
      <div className="grid grid-cols-1 md:grid-cols-[20rem_auto_2fr] w-full h-full">
        <div className="flex flex-col gap-6 px-2 animate-fade-in">
          <AnimatedSuspense>
            <Section title="Quote details">
              <QuoteUserSelect quote={quote.data} />
              {quote.data.user && <QuoteUserNameInput user={quote.data.user} />}
              {quote.data.user && (
                <QuoteUserTitleInput user={quote.data.user} />
              )}
              <QuoteUserCompanyNameInput quote={quote.data} />
              <QuoteExpiryPicker quote={quote.data} />
              <CheckoutRequiredSelect quote={quote.data} />
            </Section>
            <Section title="Fixed price">
              <PlanFixedUnitAmountInput id={quote.data.plan.id} />
              <PlanCurrencySelect id={quote.data.plan.id} />
              <PlanIntervalSelect id={quote.data.plan.id} />
              <PlanGroupSelect id={quote.data.plan.id} />
            </Section>
            <MutableEntitlements id={quote.data.plan.id} />
            <Section title="Other properties">
              <FreeTrialComboBox id={quote.data.plan.id} />
              <BillInAdvanceSelect />
              <ProrationSelect />
            </Section>
          </AnimatedSuspense>
        </div>
        <Separator orientation="vertical" className="hidden md:block" />
        <PlanItems id={quote.data.plan.id} />
      </div>
    </div>
  );
}

function AnimatedSuspense(props: { children: React.ReactNode }) {
  return (
    <Suspense
      fallback={<div className="h-8 w-48 bg-gray-200 rounded animate-pulse" />}
    >
      {props.children}
    </Suspense>
  );
}

function QuoteUserSelect(props: { quote: Quote }) {
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");

  const { data: users } = useListUsers({ email: searchValue });
  const { mutateAsync: createUser } = useCreateUser();
  const { mutateAsync: updateQuote } = useUpdateQuote(props.quote.id);

  // Filter commands based on search
  // const filtered = useMemo(() => {
  //   if (!searchValue) return users.data;

  //   return users.data.filter((user) =>
  //     user.email.toLowerCase().includes(searchValue.toLowerCase())
  //   );
  // }, [users, searchValue]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full text-left flex justify-start"
        >
          <UserIcon className="h-4 w-4" />
          {props.quote.user?.email || "Select user"}
        </Button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent className="w-[200px] p-0" align="start" side="left">
          <Command>
            <CommandInput
              placeholder="Search users..."
              value={open ? searchValue : ""}
              onValueChange={(text) => {
                startTransition(() => {
                  setSearchValue(text);
                });
              }}
              onKeyDown={async (e) => {
                if (e.key === "Enter") {
                  if (!users.data.find((f) => f.id === e.currentTarget.value)) {
                    await createUser({ email: e.currentTarget.value });
                    setOpen(false);
                  }
                }
              }}
            />
            <CommandList>
              {users.data.map((user) => (
                <CommandItem
                  key={user.id}
                  onSelect={async () => {
                    await updateQuote({ user_id: user.id });
                    setOpen(false);
                  }}
                >
                  {user.email}
                </CommandItem>
              ))}
            </CommandList>
          </Command>
        </PopoverContent>
      </PopoverPortal>
    </Popover>
  );
}

function QuoteUserNameInput(props: { user: User }) {
  const [name, setName] = useState(props.user.name || "");
  const { mutateAsync: updateUser } = useUpdateUser(props.user.id);

  useEffect(() => {
    setName(props.user.name || "");
  }, [props.user.name]);

  return (
    <div className="flex gap-1.5 px-2.5 py-1 items-center text-sm font-medium">
      <Signature className="size-4" />
      <input
        type="text"
        className="outline-none w-full"
        value={name}
        onChange={(e) => setName(e.target.value)}
        onBlur={async () => {
          if (name !== props.user.name) {
            await updateUser({ name });
          }
        }}
        placeholder="Enter user's name"
      />
    </div>
  );
}

function QuoteUserTitleInput(props: { user: User }) {
  const [title, setTitle] = useState(props.user.company_title || "");
  const { mutateAsync: updateUser } = useUpdateUser(props.user.id);

  useEffect(() => {
    setTitle(props.user.company_title || "");
  }, [props.user.company_title]);

  return (
    <div className="flex gap-1.5 px-2.5 py-1 items-center text-sm font-medium">
      <Briefcase className="size-4" />
      <input
        type="text"
        className="outline-none w-full"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        onBlur={async () => {
          if (title !== props.user.company_title) {
            await updateUser({ company_title: title });
          }
        }}
        placeholder="Enter user's company title"
      />
    </div>
  );
}

function QuoteUserCompanyNameInput(props: { quote: Quote }) {
  const [companyName, setCompanyName] = useState(
    props.quote.company_name || ""
  );
  const { mutateAsync: updateQuote } = useUpdateQuote(props.quote.id);

  useEffect(() => {
    setCompanyName(props.quote.company_name || "");
  }, [props.quote.company_name]);

  return (
    <div className="flex gap-1.5 px-2.5 py-1 items-center text-sm font-medium">
      <Building2 className="size-4" />
      <input
        type="text"
        className="outline-none w-full"
        value={companyName}
        onChange={(e) => setCompanyName(e.target.value)}
        onBlur={async () => {
          if (companyName !== props.quote.company_name) {
            await updateQuote({ company_name: companyName });
          }
        }}
        placeholder="Enter company name"
      />
    </div>
  );
}

function QuoteExpiryPicker(props: { quote: Quote }) {
  const [date, setDate] = useState<Date | null>(
    props.quote.expires_at ? new Date(props.quote.expires_at) : null
  );
  const { mutateAsync: updateQuote } = useUpdateQuote(props.quote.id);

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="w-full text-left flex justify-start font-medium"
        >
          <CalendarOff className="size-4" />
          {date ? (
            format(date, "PPP")
          ) : (
            <span className="text-gray-400">Add expiry date</span>
          )}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0">
        <Calendar
          mode="single"
          captionLayout="dropdown"
          selected={date ?? undefined}
          onSelect={(date) => {
            setDate(date ?? null);
            updateQuote({ expires_at: date ?? null });
          }}
        />
      </PopoverContent>
    </Popover>
  );
}

export function CheckoutRequiredSelect(props: { quote: Quote }) {
  const [open, setOpen] = useState(false);
  const intervals = ["Checkout Required", "No Checkout Required"];

  const { mutate } = useUpdateQuote(props.quote.id);

  const icon = useMemo(() => {
    if (props.quote.direct_checkout_enabled) {
      return <Banknote className="h-4 w-4" />;
    } else {
      return <BanknoteX className="h-4 w-4" />;
    }
  }, [props.quote.direct_checkout_enabled]);

  const label = useMemo(() => {
    if (props.quote.direct_checkout_enabled) {
      return "Checkout Required";
    } else {
      return "No Checkout Required";
    }
  }, [props.quote.direct_checkout_enabled]);

  return (
    <div>
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <div>
            <Button
              variant="ghost"
              size="sm"
              className="w-full text-left flex justify-start"
            >
              {icon}
              {label}
            </Button>
          </div>
        </PopoverTrigger>
        <PopoverContent className="w-[200px] p-0" align="start" side="left">
          <Command>
            <CommandList>
              {intervals.length === 0 && (
                <CommandEmpty>No results found.</CommandEmpty>
              )}
              <CommandGroup>
                {intervals.map((checkoutRequired) => (
                  <CommandItem
                    key={checkoutRequired}
                    onSelect={() => {
                      setOpen(false);
                      mutate({
                        direct_checkout_enabled:
                          checkoutRequired === "Checkout Required",
                      });
                    }}
                  >
                    {checkoutRequired}
                  </CommandItem>
                ))}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>
    </div>
  );
}
