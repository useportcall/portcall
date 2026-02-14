import { ReactNode } from "react";
import { Check, X } from "lucide-react";
import { getUserEntitlement } from "./api/get-user-entitlement";
import { Badge } from "./components/ui/badge";

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
      <Badge className="bg-slate-100 text-slate-600 text-xs">
        <X />
        {children}
      </Badge>
    );
  }

  return (
    <Badge variant={"outline"} className="text-xs bg-white">
      <Check />
      {children}
    </Badge>
  );
}
