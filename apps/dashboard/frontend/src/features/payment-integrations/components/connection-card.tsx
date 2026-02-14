import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Connection } from "@/hooks/api/connections";
import { Star } from "lucide-react";
import { ConnectionActionsMenu } from "./connection-actions-menu";
import { getProviderDescription, getProviderIcon } from "./provider-meta";

export function ConnectionCard({ connection, isDefault }: { connection: Connection; isDefault: boolean }) {
  return (
    <Card className={`relative transition-all hover:shadow-md ${isDefault ? "ring-2 ring-primary/20 bg-primary/5" : ""}`}>
      <CardContent className="p-4">
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3">
            <div className="mt-0.5">{getProviderIcon(connection.source, "w-10 h-10")}</div>
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-2">
                <h3 className="font-semibold text-sm">{connection.name}</h3>
                {isDefault && <Badge variant="secondary" className="text-xs px-1.5 py-0 h-5 bg-primary/10 text-primary border-0"><Star className="w-3 h-3 mr-1 fill-current" />Default</Badge>}
              </div>
              <p className="text-xs text-muted-foreground">{getProviderDescription(connection.source)}</p>
              {connection.source !== "local" && connection.public_key && (
                <code className="text-xs text-muted-foreground font-mono mt-1 bg-muted px-1.5 py-0.5 rounded w-fit">
                  {connection.public_key.slice(0, 12)}...
                </code>
              )}
              {connection.source === "braintree" && connection.webhook_url && (
                <p className="text-xs text-muted-foreground mt-1 truncate max-w-[220px]" title={connection.webhook_url}>
                  Webhook: {connection.webhook_url}
                </p>
              )}
            </div>
          </div>
          <ConnectionActionsMenu connection={connection} isDefault={isDefault} />
        </div>
      </CardContent>
    </Card>
  );
}
