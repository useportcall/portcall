import { Check, Copy } from "lucide-react";
import { useState } from "react";

export default function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);

  const onClick = () => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 1000); // Reset after 2 seconds
  };

  if (copied) {
    return <Check className="w-4 h-4" />;
  }

  return (
    <button onClick={onClick} className="flex justify-center items-center">
      <Copy className="w-4 h-4 hover:text-muted-foreground active:text-muted-foreground/70" />
    </button>
  );
}
