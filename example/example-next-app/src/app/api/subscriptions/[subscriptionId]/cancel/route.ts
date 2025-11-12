import { callPortcallApi } from "@/lib/api";
import { NextRequest } from "next/server";

export async function POST(
  _: NextRequest,
  { params }: { params: { subscriptionId: string } }
) {
  const { subscriptionId } = params;

  return callPortcallApi(
    "POST",
    `/v1/subscriptions/${subscriptionId}/cancel`,
    {}
  );
}
