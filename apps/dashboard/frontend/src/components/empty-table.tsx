import { Card } from "@/components/ui/card";
import { ReactNode } from "react";

export default function EmptyTable({
  message,
  button,
}: {
  message: string;
  button?: ReactNode;
}) {
  return (
    <Card className="w-full flex flex-col justify-center items-center p-10 gap-4">
      <p className="text-sm text-muted-foreground">{message}</p>
      {button}
    </Card>
  );
}
