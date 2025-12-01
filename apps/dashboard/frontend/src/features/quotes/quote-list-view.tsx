import { Card } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { useCreateQuote, useListQuotes } from "@/hooks/api/quotes";
import { Loader2, Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";

export default function QuoteListView() {
  const navigate = useNavigate();
  const { data: quotes } = useListQuotes();

  return (
    <div className="w-full p-4 lg:p-10 flex flex-col gap-6">
      <div className="flex flex-row justify-between">
        <div className="flex flex-col space-y-8 w-full">
          <div className="flex justify-between w-full">
            <div className="flex flex-col justify-start space-y-2">
              <h1 className="text-lg md:text-xl font-semibold">Quotes</h1>
              <p className="text-sm text-slate-400">
                Manage existing quotes or add new quotes here.
              </p>
            </div>
            <div className="flex flex-row gap-4 items-end">
              <CreateQuoteButton />
            </div>
          </div>
        </div>
      </div>
      <Card className="p-2 rounded-xl animate-fade-in">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-80">Name</TableHead>
              <TableHead className="w-40">Status</TableHead>
              <TableHead className="w-40">Group</TableHead>
              <TableHead className="w-40">Fixed charge</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {quotes.data.map((quote) => (
              <TableRow
                key={quote.id}
                onClick={() => navigate(`/quotes/${quote.id}`)}
                className="rounded hover:bg-slate-50 cursor-pointer"
              >
                <TableCell className="w-80">
                  {quote.id || "Untitled Quote"}
                </TableCell>
                <TableCell className="w-40">{quote.status}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        {!quotes.data.length && (
          <div className="w-full flex flex-col justify-center items-center p-10 gap-4">
            <p className="text-sm text-slate-400">No quotes found.</p>
          </div>
        )}
      </Card>
    </div>
  );
}

export function CreateQuoteButton() {
  const { mutate, isPending } = useCreateQuote();

  return (
    <Button
      size="sm"
      variant="outline"
      onClick={() => {
        mutate({});
      }}
    >
      {isPending && <Loader2 className="size-4 animate-spin" />}
      {!isPending && <Plus className="size-4" />}
      Add quote
    </Button>
  );
}
