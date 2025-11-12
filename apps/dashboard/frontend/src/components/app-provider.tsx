import { AppContext } from "@/contexts/app-context";
import { useListApps } from "@/hooks";
import { useAuth } from "@/lib/keycloak/auth";
import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

// TODO: refactor
export const AppProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { isAuthenticated } = useAuth();
  const { data: apps } = useListApps();

  const [p, setSelectedProject] = useState<string | null>(() => {
    return localStorage.getItem("selectedProject");
  });
  const [id, setId] = useState<string | null>(null);
  const queryClient = useQueryClient();

  useEffect(() => {
    if (!p) return;
    localStorage.setItem("selectedProject", p);
  }, [p]);

  useEffect(() => {
    if (!isAuthenticated) {
      setSelectedProject(null);
      queryClient.clear();
      return;
    } else {
      if (!apps?.data?.length) {
        console.log("No apps available");
        setId(null);
        setSelectedProject(null);
        return;
      } else {
        console.log("Setting default app", apps.data[0].id.toString());
        setId(apps.data[0].id.toString());
        setSelectedProject(apps.data[0].id.toString());
      }
    }
  }, [isAuthenticated, apps?.data?.length]);

  return (
    <AppContext.Provider value={{ id: p ?? id, setId: setSelectedProject }}>
      {children}
    </AppContext.Provider>
  );
};
