import { Suspense } from "react";

export function AnimatedSuspense({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={<div className="h-8 w-48 bg-muted rounded animate-pulse" />}>
      {children}
    </Suspense>
  );
}
