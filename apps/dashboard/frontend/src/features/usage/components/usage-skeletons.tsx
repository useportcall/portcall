import { Card, CardContent, CardHeader } from "@/components/ui/card";

export function QuotaCardSkeleton() {
  return <CardSkeleton />;
}

export function PlanCardSkeleton() {
  return <CardSkeleton />;
}

export function InvoicesCardSkeleton() {
  return (
    <Card>
      <CardHeader className="space-y-2">
        <div className="h-5 w-32 rounded bg-muted animate-pulse" />
        <div className="h-4 w-60 rounded bg-muted animate-pulse" />
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="h-8 rounded bg-muted animate-pulse" />
        <div className="h-8 rounded bg-muted animate-pulse" />
        <div className="h-8 rounded bg-muted animate-pulse" />
      </CardContent>
    </Card>
  );
}

function CardSkeleton() {
  return (
    <Card>
      <CardHeader className="space-y-2">
        <div className="h-5 w-32 rounded bg-muted animate-pulse" />
        <div className="h-4 w-52 rounded bg-muted animate-pulse" />
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="h-8 rounded bg-muted animate-pulse" />
        <div className="h-8 rounded bg-muted animate-pulse" />
      </CardContent>
    </Card>
  );
}
