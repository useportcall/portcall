import { AppContext } from "@/contexts/app-context";
import { useListApps } from "@/hooks";
import { useAuth } from "@/lib/keycloak/auth";
import { useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { user } = useAuth();
  const { data } = useListApps();

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
    if (!user) {
      setSelectedProject(null);
      queryClient.clear();
      return;
    } else if (data?.data?.length) {
      setId(data.data[0].id.toString());
      setSelectedProject(data.data[0].id.toString());
    }
  }, [user, data]);

  return (
    <AppContext.Provider value={{ id: p ?? id, setId: setSelectedProject }}>
      {children}
    </AppContext.Provider>
  );
};
