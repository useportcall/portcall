import type { ApiEndpoint } from "../lib/openapi";

interface EndpointSectionProps {
  endpoint: ApiEndpoint;
  categoryId: string;
}

const methodColors: Record<string, string> = {
  GET: "bg-blue-100 text-blue-800 border-blue-200 dark:bg-blue-900/30 dark:text-blue-300 dark:border-blue-800",
  POST: "bg-green-100 text-green-800 border-green-200 dark:bg-green-900/30 dark:text-green-300 dark:border-green-800",
  PUT: "bg-yellow-100 text-yellow-800 border-yellow-200 dark:bg-yellow-900/30 dark:text-yellow-300 dark:border-yellow-800",
  PATCH:
    "bg-orange-100 text-orange-800 border-orange-200 dark:bg-orange-900/30 dark:text-orange-300 dark:border-orange-800",
  DELETE:
    "bg-red-100 text-red-800 border-red-200 dark:bg-red-900/30 dark:text-red-300 dark:border-red-800",
};

export function EndpointSection({
  endpoint,
  categoryId,
}: EndpointSectionProps) {
  const pathParams =
    endpoint.parameters?.filter((p) => p.location === "path") || [];
  const queryParams =
    endpoint.parameters?.filter((p) => p.location === "query") || [];
  const bodyParams =
    endpoint.parameters?.filter((p) => p.location === "body") || [];

  return (
    <div
      id={`${categoryId}-${endpoint.id}`}
      data-endpoint-id={`${categoryId}-${endpoint.id}`}
      className="scroll-mt-24"
    >
      <div className="mb-6">
        <div className="flex items-center gap-3 mb-4">
          <span
            className={`px-3 py-1.5 rounded-md text-xs font-bold border ${
              methodColors[endpoint.method] ||
              "bg-muted text-foreground dark:bg-muted dark:text-foreground"
            }`}
          >
            {endpoint.method}
          </span>
          <code className="text-sm font-mono bg-muted px-4 py-1.5 rounded-md border border-border text-foreground">
            {endpoint.path}
          </code>
        </div>
        <h3 className="text-2xl font-semibold mb-3 text-foreground">
          {endpoint.title}
        </h3>
        <p className="text-muted-foreground leading-relaxed">
          {endpoint.description}
        </p>
        {endpoint.authentication && (
          <div className="mt-4 inline-flex items-center gap-2 text-sm text-amber-700 dark:text-amber-300 bg-amber-50 dark:bg-amber-950/30 px-3 py-1.5 rounded-md border border-amber-200 dark:border-amber-800">
            🔒 Authentication required
          </div>
        )}
      </div>

      {/* Parameters */}
      {endpoint.parameters && endpoint.parameters.length > 0 && (
        <div className="mb-6">
          <h4 className="text-lg font-semibold mb-4 text-foreground">
            Parameters
          </h4>

          {pathParams.length > 0 && (
            <div className="mb-6">
              <h5 className="text-sm font-semibold text-muted-foreground mb-3 uppercase tracking-wide">
                Path Parameters
              </h5>
              <div className="overflow-x-auto rounded-lg border border-border">
                <table className="min-w-full divide-y divide-border">
                  <thead className="bg-muted/50">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Name
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Type
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Required
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Description
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-background divide-y divide-border">
                    {pathParams.map((param) => (
                      <tr
                        key={param.name}
                        className="hover:bg-muted/30 transition-colors"
                      >
                        <td className="px-4 py-3 text-sm font-mono font-semibold text-foreground">
                          {param.name}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground">
                          {param.type}
                        </td>
                        <td className="px-4 py-3 text-sm">
                          {param.required ? (
                            <span className="text-red-600 dark:text-red-400 font-medium">
                              Yes
                            </span>
                          ) : (
                            <span className="text-muted-foreground">No</span>
                          )}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground leading-relaxed">
                          {param.description}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {queryParams.length > 0 && (
            <div className="mb-6">
              <h5 className="text-sm font-semibold text-muted-foreground mb-3 uppercase tracking-wide">
                Query Parameters
              </h5>
              <div className="overflow-x-auto rounded-lg border border-border">
                <table className="min-w-full divide-y divide-border">
                  <thead className="bg-muted/50">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Name
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Type
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Required
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Description
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-background divide-y divide-border">
                    {queryParams.map((param) => (
                      <tr
                        key={param.name}
                        className="hover:bg-muted/30 transition-colors"
                      >
                        <td className="px-4 py-3 text-sm font-mono font-semibold text-foreground">
                          {param.name}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground">
                          {param.type}
                        </td>
                        <td className="px-4 py-3 text-sm">
                          {param.required ? (
                            <span className="text-red-600 dark:text-red-400 font-medium">
                              Yes
                            </span>
                          ) : (
                            <span className="text-muted-foreground">No</span>
                          )}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground leading-relaxed">
                          {param.description}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {bodyParams.length > 0 && (
            <div className="mb-6">
              <h5 className="text-sm font-semibold text-muted-foreground mb-3 uppercase tracking-wide">
                Body Parameters
              </h5>
              <div className="overflow-x-auto rounded-lg border border-border">
                <table className="min-w-full divide-y divide-border">
                  <thead className="bg-muted/50">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Name
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Type
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Required
                      </th>
                      <th className="px-4 py-3 text-left text-xs font-semibold text-foreground uppercase tracking-wide">
                        Description
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-background divide-y divide-border">
                    {bodyParams.map((param) => (
                      <tr
                        key={param.name}
                        className="hover:bg-muted/30 transition-colors"
                      >
                        <td className="px-4 py-3 text-sm font-mono font-semibold text-foreground">
                          {param.name}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground">
                          {param.type}
                        </td>
                        <td className="px-4 py-3 text-sm">
                          {param.required ? (
                            <span className="text-red-600 dark:text-red-400 font-medium">
                              Yes
                            </span>
                          ) : (
                            <span className="text-muted-foreground">No</span>
                          )}
                        </td>
                        <td className="px-4 py-3 text-sm text-muted-foreground leading-relaxed">
                          {param.description}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
