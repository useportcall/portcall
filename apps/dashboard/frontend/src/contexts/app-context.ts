import { createContext } from "react";

type AppContextType = {
  id: string | null;
  setId: (id: string) => void;
};

export const AppContext = createContext<AppContextType | undefined>(undefined);
