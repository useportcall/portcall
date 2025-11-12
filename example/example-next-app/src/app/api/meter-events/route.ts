import { callPortcallApi } from "@/lib/api";

export async function POST(request: Request) {
  const { featureId, userId, usage } = await request.json();

  return callPortcallApi("POST", "/v1/meter-events", {
    feature_id: featureId,
    user_id: userId,
    usage: usage,
  });
}
