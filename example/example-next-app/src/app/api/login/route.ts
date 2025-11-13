import { API_KEY, PUBLIC_API } from "@/constants";
import { NextResponse } from "next/server";

export async function POST(req: Request) {
  try {
    const { email } = await req.json();

    const userRes = await fetch(
      `${PUBLIC_API}/v1/users?email=${encodeURIComponent(email)}`,
      {
        headers: {
          "X-API-KEY": process.env.PC_API_TOKEN!,
          "Content-Type": "application/json",
        },
      }
    );
    const result = await userRes.json();

    const user = result?.data?.[0] ?? null;

    if (!user) {
      return NextResponse.json({ error: "user not found" }, { status: 404 });
    }

    return NextResponse.json({ user });
  } catch (err) {
    return NextResponse.json({ error: "invalid request" }, { status: 400 });
  }
}
