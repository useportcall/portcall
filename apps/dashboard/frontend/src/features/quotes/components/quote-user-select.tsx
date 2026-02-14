import { Button } from "@/components/ui/button";
import {
  Command,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { useCreateUser, useListUsers } from "@/hooks";
import { useUpdateQuote } from "@/hooks/api/quotes";
import { Quote } from "@/models/quote";
import { PopoverPortal } from "@radix-ui/react-popover";
import { UserIcon } from "lucide-react";
import { startTransition, useState } from "react";

export function QuoteUserSelect({ quote }: { quote: Quote }) {
  const [open, setOpen] = useState(false);
  const [searchValue, setSearchValue] = useState("");
  const { data: users } = useListUsers({ email: searchValue });
  const { mutateAsync: createUser } = useCreateUser();
  const { mutateAsync: updateQuote } = useUpdateQuote(quote.id);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="ghost" size="sm" className="w-full justify-start text-left">
          <UserIcon className="h-4 w-4" />
          {quote.user?.email || "Select user"}
        </Button>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent className="w-[200px] p-0" align="start" side="left">
          <Command>
            <CommandInput
              placeholder="Search users..."
              value={open ? searchValue : ""}
              onValueChange={(text) => startTransition(() => setSearchValue(text))}
              onKeyDown={async (e) => {
                if (
                  e.key === "Enter" &&
                  !users.data.find((item) => item.id === e.currentTarget.value)
                ) {
                  await createUser({ email: e.currentTarget.value });
                  setOpen(false);
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
