import { AppContext } from "@/contexts/app-context";
import { useListApps } from "@/hooks";
import { useAuth } from "@/lib/keycloak/auth";
import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isAuthenticated } = useAuth();
  const { data: apps } = useListApps();
  const queryClient = useQueryClient();

  const [selectedProject, setSelectedProject] = useState<string | null>(() => {
    return localStorage.getItem("selectedProject");
  });

  // Compute valid app ID synchronously - this ensures we never use a stale ID
  const validatedAppId = useMemo(() => {
    const appList = apps?.data ?? [];

    // Check if selected project exists in the fetched apps
    if (selectedProject && appList.length > 0) {
      const exists = appList.some(
        (app) => app.id.toString() === selectedProject,
      );
      if (exists) {
        return selectedProject;
      }
    }

    // If no valid selection, use first app if available
    if (appList.length > 0) {
      return appList[0].id.toString();
    }

    return null;
  }, [selectedProject, apps?.data]);

  // Sync localStorage when validated ID changes
  useEffect(() => {
    if (validatedAppId) {
      localStorage.setItem("selectedProject", validatedAppId);
      // Update state if it was stale
      if (selectedProject !== validatedAppId) {
        setSelectedProject(validatedAppId);
      }
    } else {
      localStorage.removeItem("selectedProject");
      if (selectedProject) {
        setSelectedProject(null);
      }
    }
  }, [validatedAppId]);

  // Clear on logout
  useEffect(() => {
    if (!isAuthenticated) {
      setSelectedProject(null);
      localStorage.removeItem("selectedProject");
      queryClient.clear();
    }
  }, [isAuthenticated, queryClient]);

  // Invalidate app-specific queries when selected app changes
  useEffect(() => {
    if (validatedAppId) {
      queryClient.invalidateQueries({
        predicate: (query) => {
          const key = query.queryKey[0];
          return (
            typeof key === "string" &&
            key !== "/apps" &&
            !key.startsWith("/api/apps")
          );
        },
      });
    }
  }, [validatedAppId, queryClient]);

  return (
    <AppContext.Provider
      value={{ id: validatedAppId, setId: setSelectedProject }}
    >
      {children}
    </AppContext.Provider>
  );
};
