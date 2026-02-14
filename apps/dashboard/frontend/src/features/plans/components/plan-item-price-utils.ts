export function parsePriceToCents(value: string): number {
  const parsed = Number.parseFloat(value);
  if (Number.isNaN(parsed)) return 0;
  return Math.round(parsed * 100);
}

export function formatPriceInput(value: string): string {
  const parsed = Number.parseFloat(value);
  if (Number.isNaN(parsed)) return "0.00";
  return parsed.toFixed(2);
}
