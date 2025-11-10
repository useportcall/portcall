import { Button } from "@/components/ui/button";
import { Calculator } from "lucide-react";

export function ProrationSelect() {
  return (
    <Button
      disabled
      variant="ghost"
      size="sm"
      className="w-full text-left flex justify-start"
    >
      <Calculator className="h-4 w-4" />
      Proration (coming soon)
    </Button>
  );
}
