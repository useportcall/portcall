import { Badge } from "@/components/ui/badge";
import { Wrench } from "lucide-react";

export function DogfoodBadge() {
  return (
    <div className="flex items-center gap-2">
      <Badge
        variant="outline"
        className="bg-amber-50 dark:bg-amber-950 border-amber-300 dark:border-amber-700 text-amber-700 dark:text-amber-300 font-medium flex items-center gap-1.5"
      >
        <Wrench className="h-3 w-3" />
        Dogfood App
      </Badge>
    </div>
  );
}
