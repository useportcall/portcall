import axios from "axios";

export function getSessionErrorMessage(error: unknown): string | null {
  if (!error) {
    return null;
  }
  if (axios.isAxiosError<{ error?: string }>(error)) {
    const message = error.response?.data?.error;
    return typeof message === "string" && message.length > 0 ? message : null;
  }
  if (error instanceof Error && error.message.length > 0) {
    return error.message;
  }
  return null;
}
