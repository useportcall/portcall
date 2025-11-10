import { AppContext } from "@/contexts/app-context";
import { useContext } from "react";

export const useApp = () => {
  const context = useContext(AppContext);
  if (!context) throw new Error("useApp must be used within a AppProvider");
  return context;
};
