import { Card } from "@/components/ui/card";

interface HomeStatCardProps {
  title: string;
  value: number;
  meta: string;
  testId: string;
}

export function HomeStatCard({ title, value, meta, testId }: HomeStatCardProps) {
  return (
    <Card className="p-4 space-y-1">
      <p className="text-sm text-muted-foreground">{title}</p>
      <p data-testid={testId} className="text-2xl font-semibold">
        {value}
      </p>
      <p className="text-xs text-muted-foreground">{meta}</p>
    </Card>
  );
}
