import { useEffect, useState } from "react";
import Prism from "prismjs";
import "prismjs/themes/prism-tomorrow.css";
import "prismjs/components/prism-bash";
import "prismjs/components/prism-json";
import { Copy, Check } from "lucide-react";
import type { ApiEndpoint } from "../lib/openapi";

interface CodePanelProps {
  endpoint: ApiEndpoint | null;
}

export function CodePanel({ endpoint }: CodePanelProps) {
  const [copiedRequest, setCopiedRequest] = useState(false);
  const [copiedResponse, setCopiedResponse] = useState(false);
  const [activeTab, setActiveTab] = useState<"request" | "response">(
    "response",
  );

  useEffect(() => {
    Prism.highlightAll();
  }, [endpoint, activeTab]);

  const copyToClipboard = async (
    text: string,
    type: "request" | "response",
  ) => {
    await navigator.clipboard.writeText(text);
    if (type === "request") {
      setCopiedRequest(true);
      setTimeout(() => setCopiedRequest(false), 2000);
    } else {
      setCopiedResponse(true);
      setTimeout(() => setCopiedResponse(false), 2000);
    }
  };

  if (!endpoint) {
    return (
      <div className="bg-background h-full p-8 flex items-center justify-center">
        <div className="text-center text-muted-foreground">
          <p className="text-sm">Select an endpoint to view code examples</p>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-background h-full overflow-y-auto">
      <div className="sticky top-0 bg-background border-b border-border z-10">
        <div className="p-6">
          <h3 className="font-semibold text-foreground mb-4">Examples</h3>
          {(endpoint.requestBody || endpoint.response) && (
            <div className="flex gap-2">
              {endpoint.requestBody && (
                <button
                  onClick={() => setActiveTab("request")}
                  className={`px-4 py-2 text-sm rounded-md transition-all font-medium ${
                    activeTab === "request"
                      ? "bg-foreground text-background shadow-sm"
                      : "bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground"
                  }`}
                >
                  Request
                </button>
              )}
              {endpoint.response && (
                <button
                  onClick={() => setActiveTab("response")}
                  className={`px-4 py-2 text-sm rounded-md transition-all font-medium ${
                    activeTab === "response"
                      ? "bg-foreground text-background shadow-sm"
                      : "bg-muted text-muted-foreground hover:bg-muted/80 hover:text-foreground"
                  }`}
                >
                  Response
                </button>
              )}
            </div>
          )}
        </div>
      </div>

      <div className="p-6">
        {activeTab === "request" && endpoint.requestBody && (
          <div>
            <div className="flex items-center justify-between mb-3">
              <h4 className="text-sm font-semibold text-foreground">
                Request Body
              </h4>
              <button
                onClick={() =>
                  copyToClipboard(endpoint.requestBody!.example, "request")
                }
                className="flex items-center gap-1.5 px-3 py-1.5 text-xs bg-muted hover:bg-muted/80 rounded-md transition-colors text-foreground font-medium"
              >
                {copiedRequest ? <Check size={14} /> : <Copy size={14} />}
                {copiedRequest ? "Copied!" : "Copy"}
              </button>
            </div>
            <pre className="!mt-0 rounded-lg overflow-x-auto !bg-[#2d2d2d] !p-5 border border-border/50">
              <code className="language-json !text-sm">
                {endpoint.requestBody.example}
              </code>
            </pre>
          </div>
        )}

        {activeTab === "response" && endpoint.response && (
          <div>
            <div className="flex items-center justify-between mb-3">
              <h4 className="text-sm font-semibold text-foreground">
                Response Body
              </h4>
              <button
                onClick={() =>
                  copyToClipboard(endpoint.response!.example, "response")
                }
                className="flex items-center gap-1.5 px-3 py-1.5 text-xs bg-muted hover:bg-muted/80 rounded-md transition-colors text-foreground font-medium"
              >
                {copiedResponse ? <Check size={14} /> : <Copy size={14} />}
                {copiedResponse ? "Copied!" : "Copy"}
              </button>
            </div>
            <pre className="!mt-0 rounded-lg overflow-x-auto !bg-[#2d2d2d] !p-5 border border-border/50">
              <code className="language-json !text-sm">
                {endpoint.response.example}
              </code>
            </pre>
          </div>
        )}

        {!endpoint.requestBody && !endpoint.response && (
          <div className="text-center text-muted-foreground py-12">
            <p className="text-sm">
              No code examples available for this endpoint
            </p>
          </div>
        )}
      </div>
    </div>
  );
}
