import { callPortcallApi } from "@/lib/api";
import { NextRequest } from "next/server";

export async function POST(
  request: NextRequest,
  { params }: { params: { subscriptionId: string } }
) {
  const { subscriptionId } = params;
  const { planId } = await request.json();

  return callPortcallApi("POST", `/v1/subscriptions/${subscriptionId}`, {
    plan_id: planId,
  });
}
