import type { ApiCategory } from "../lib/openapi";

interface SidebarProps {
  endpoints: ApiCategory[];
}

export function Sidebar({ endpoints }: SidebarProps) {
  const scrollToSection = (categoryId: string, endpointId?: string) => {
    const id = endpointId ? `${categoryId}-${endpointId}` : categoryId;
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  };

  return (
    <nav className="px-6 space-y-8">
      {endpoints.map((category) => (
        <div key={category.id}>
          <button
            onClick={() => scrollToSection(category.id)}
            className="font-semibold text-sm text-foreground hover:text-primary mb-3 block w-full text-left transition-colors"
          >
            {category.title}
          </button>
          <ul className="space-y-2 ml-0 border-l-2 border-border pl-4">
            {category.endpoints.map((endpoint) => (
              <li key={endpoint.id}>
                <button
                  onClick={() => scrollToSection(category.id, endpoint.id)}
                  className="text-sm text-muted-foreground hover:text-foreground block w-full text-left py-1 transition-colors"
                >
                  {endpoint.title}
                </button>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </nav>
  );
}
