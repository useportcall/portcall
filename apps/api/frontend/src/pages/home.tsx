import { useState, useEffect } from "react";
import { Menu, X, Loader2 } from "lucide-react";
import { Sidebar } from "../components/sidebar";
import { EndpointSection } from "../components/endpoint-section";
import { CodePanel } from "../components/code-panel";
import {
  getApiEndpoints,
  type ApiCategory,
  type ApiEndpoint,
} from "../lib/openapi";

export default function HomePage() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [activeEndpoint, setActiveEndpoint] = useState<string | null>(null);
  const [apiEndpoints, setApiEndpoints] = useState<ApiCategory[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Fetch API endpoints from OpenAPI spec
    getApiEndpoints()
      .then((endpoints) => {
        setApiEndpoints(endpoints);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Failed to load API spec:", err);
        setError("Failed to load API documentation");
        setLoading(false);
      });
  }, []);

  useEffect(() => {
    if (apiEndpoints.length === 0) return;

    // Set up intersection observer to track which endpoint is in view
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const id = entry.target.id;
            if (id) {
              setActiveEndpoint(id);
            }
          }
        });
      },
      {
        rootMargin: "-100px 0px -50% 0px",
        threshold: 0,
      },
    );

    // Observe all endpoint sections
    const sections = document.querySelectorAll("[data-endpoint-id]");
    sections.forEach((section) => observer.observe(section));

    return () => observer.disconnect();
  }, [apiEndpoints]);

  // Get current endpoint data
  const getCurrentEndpoint = (): ApiEndpoint | null => {
    if (!activeEndpoint) return null;
    for (const category of apiEndpoints) {
      const endpoint = category.endpoints.find(
        (e) => `${category.id}-${e.id}` === activeEndpoint,
      );
      if (endpoint) return endpoint;
    }
    return null;
  };

  const currentEndpoint = getCurrentEndpoint();

  if (loading) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="flex items-center gap-3 text-muted-foreground">
          <Loader2 className="h-6 w-6 animate-spin" />
          <span>Loading API documentation...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 mb-4">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-primary text-primary-foreground rounded-md"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="flex h-16 items-center px-6">
          <div className="flex items-center gap-3">
            <button
              onClick={() => setSidebarOpen(!sidebarOpen)}
              className="lg:hidden p-2 hover:bg-accent rounded-md transition-colors"
              aria-label="Toggle sidebar"
            >
              {sidebarOpen ? (
                <X size={20} className="text-foreground" />
              ) : (
                <Menu size={20} className="text-foreground" />
              )}
            </button>
            <div className="flex items-center gap-3">
              <img src="/logo.png" alt="Portcall" className="h-8 w-8 logo" />
              <span className="text-lg font-semibold text-foreground">
                API Reference
              </span>
            </div>
          </div>
          <div className="ml-auto flex items-center gap-4">
            <a
              href="https://useportcall.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-sm text-muted-foreground hover:text-foreground transition-colors"
            >
              useportcall.com
            </a>
          </div>
        </div>
      </header>

      <div className="flex">
        {/* Sidebar */}
        <aside
          className={`
            fixed lg:static inset-y-0 left-0 z-40 w-64 bg-background border-r border-border transform transition-transform duration-200 ease-in-out lg:transform-none
            ${sidebarOpen ? "translate-x-0" : "-translate-x-full lg:translate-x-0"}
          `}
        >
          <div className="sticky top-16 max-h-[calc(100vh-4rem)] overflow-y-auto py-6">
            <Sidebar endpoints={apiEndpoints} />
          </div>
        </aside>

        {/* Overlay for mobile */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 bg-black/50 z-30 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
        )}

        {/* Main Content */}
        <main className="flex-1 min-w-0 lg:flex">
          {/* Content Column */}
          <div className="flex-1 px-6 py-10 lg:px-12 lg:py-12 max-w-4xl">
            <div className="mb-16">
              <h1 className="text-4xl lg:text-5xl font-bold mb-6 text-foreground">
                API Reference
              </h1>
              <p className="text-lg text-muted-foreground mb-8 leading-relaxed">
                Welcome to the Portcall API reference. This documentation
                provides detailed information about all available endpoints,
                including request parameters, response formats, and example
                usage.
              </p>
              <div className="bg-muted/50 border border-border rounded-lg p-6">
                <h3 className="text-sm font-semibold mb-3 text-foreground uppercase tracking-wide">
                  Base URL
                </h3>
                <code className="text-sm font-mono bg-background px-4 py-2 rounded-md border border-border inline-block">
                  https://api.portcall.com
                </code>
              </div>
            </div>

            {/* Endpoint Sections */}
            {apiEndpoints.map((category) => (
              <div key={category.id} id={category.id} className="mb-16">
                <h2 className="text-3xl font-bold mb-8 text-foreground">
                  {category.title}
                </h2>
                <div className="space-y-10">
                  {category.endpoints.map((endpoint) => (
                    <EndpointSection
                      key={endpoint.id}
                      endpoint={endpoint}
                      categoryId={category.id}
                    />
                  ))}
                </div>
              </div>
            ))}
          </div>

          {/* Code Panel - Right Column */}
          <div className="hidden lg:block lg:w-[500px] xl:w-[600px] border-l border-border">
            <div className="sticky top-16 h-[calc(100vh-4rem)]">
              <CodePanel endpoint={currentEndpoint} />
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
