import { createContext, useContext } from "react";

export type SubmitState = "idle" | "loading" | "success";

export const SubmitStateContext = createContext<{
  state: SubmitState;
  setState: (state: SubmitState) => void;
}>({ state: "idle", setState: () => {} });

export function useSubmitState() {
  return useContext(SubmitStateContext);
}
