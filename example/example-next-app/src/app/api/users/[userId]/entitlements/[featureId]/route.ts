import { callPortcallApi } from "@/lib/api";
import { NextRequest } from "next/server";

export async function GET(
  _: NextRequest,
  { params }: { params: { userId: string; featureId: string } }
) {
  const { userId, featureId } = params;
  return callPortcallApi(
    "GET",
    `/v1/users/${userId}/entitlements/${featureId}`
  );
}
