import { Button } from "@/components/ui/button";
import { Calendar } from "lucide-react";

export function BillInAdvanceSelect() {
  return (
    <Button
      disabled
      variant="ghost"
      size="sm"
      className="w-full text-left flex justify-start"
    >
      <Calendar className="h-4 w-4" />
      Bill in advance (coming soon)
    </Button>
  );
}
