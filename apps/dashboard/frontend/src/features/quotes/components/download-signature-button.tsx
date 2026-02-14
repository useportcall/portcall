import { Button } from "@/components/ui/button";
import { useAxiosClient } from "@/hooks/api/api";
import { Quote } from "@/models/quote";
import { Download } from "lucide-react";
import { useState } from "react";

export function DownloadSignatureButton({ quote }: { quote: Quote }) {
  const client = useAxiosClient();
  const [loading, setLoading] = useState(false);
  const canDownload = quote.status === "accepted";

  return (
    <Button
      variant="outline"
      size="sm"
      disabled={!canDownload || loading}
      onClick={async () => {
        setLoading(true);
        try {
          const response = await client.get(`/quotes/${quote.id}/signature`, {
            responseType: "blob",
          });
          const href = URL.createObjectURL(response.data);
          const anchor = document.createElement("a");
          anchor.href = href;
          anchor.download = `${quote.id}-signature.png`;
          document.body.appendChild(anchor);
          anchor.click();
          anchor.remove();
          URL.revokeObjectURL(href);
        } finally {
          setLoading(false);
        }
      }}
    >
      <Download className="size-4" />
      {loading ? "Downloading..." : "Signed copy"}
    </Button>
  );
}
