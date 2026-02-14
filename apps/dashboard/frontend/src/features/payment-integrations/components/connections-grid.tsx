import EmptyTable from "@/components/empty-table";
import { useGetAppConfig } from "@/hooks/api/app-configs";
import { useListConnections } from "@/hooks/api/connections";
import { ConnectionCard } from "./connection-card";

export function ConnectionsGrid() {
  const { data: config } = useGetAppConfig();
  const { data: connections } = useListConnections();

  if (!connections || !connections.data.length) {
    return (
      <EmptyTable
        message="No payment providers connected yet. Add one to start accepting payments."
        button=""
      />
    );
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      {connections.data.map((connection) => (
        <ConnectionCard
          key={connection.id}
          connection={connection}
          isDefault={connection.id === config?.data?.default_connection}
        />
      ))}
    </div>
  );
}
