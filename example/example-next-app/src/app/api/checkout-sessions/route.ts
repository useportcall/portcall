// Runs on the server
import { callPortcallApi } from "@/lib/api";

export async function POST(request: Request) {
  const { planId, userId } = await request.json();

  return callPortcallApi("POST", "/v1/checkout-sessions", {
    plan_id: planId,
    user_id: userId,
    cancel_url: "http://localhost:3000?status=canceled", // TODO: fix
    redirect_url: "http://localhost:3000?status=success", // TODO: fix
  });
}
