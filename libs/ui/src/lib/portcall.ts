import { Portcall } from "@portcall/sdk";

/**
 * Server-side Portcall SDK client
 * This should only be used in server components and server actions
 */
export function getPortcallClient() {
  const apiKey = process.env.PC_API_TOKEN;
  const baseURL = process.env.PC_API_BASE_URL;

  if (!apiKey) {
    throw new Error("PC_API_TOKEN environment variable is not set");
  }

  return new Portcall({
    apiKey,
    baseURL,
  });
}
