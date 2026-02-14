export function getQuotaTitle(quota: number | null | undefined): string {
  if (quota === -1 || !quota) return "no limit";
  return quota.toLocaleString("en-US", { style: "decimal" });
}

export function sanitizeQuotaInput(value: string): string | null {
  if (value === "-" || value === "") return value;
  const numeric = Number(value);
  if (Number.isNaN(numeric)) return null;
  return numeric.toFixed(0);
}

export function parseQuotaInput(value: string): number {
  const numeric = Number.parseInt(value, 10);
  if (Number.isNaN(numeric)) return -1;
  return numeric;
}
