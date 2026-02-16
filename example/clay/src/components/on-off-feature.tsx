import { ReactNode } from "react";
import { Badge } from "./ui/badge";
import { Check, X } from "lucide-react";
import { getUserEntitlement } from "@/lib/get-user-entitlement";

export async function OnOffFeature({
  featureId,
  children,
}: {
  featureId: string;
  children: ReactNode;
}) {
  const entitlement = await getUserEntitlement(featureId);

  if (!entitlement || !entitlement.enabled) {
    return (
      <Badge className="bg-slate-100 text-slate-600">
        <X />
        {children}
      </Badge>
    );
  }

  return (
    <Badge variant={"outline"}>
      <Check />
      {children}
    </Badge>
  );
}
