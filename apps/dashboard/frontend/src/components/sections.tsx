export function Section(props: { title: string; children: React.ReactNode }) {
  return (
    <div className="space-y-2">
      <h2 className="text-sm text-muted-foreground px-2">{props.title}</h2>
      <div className="space-y-2">{props.children}</div>
    </div>
  );
}
