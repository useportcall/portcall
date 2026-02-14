/** Returns the divisor for converting minor-unit amounts to major-unit. */
export function getCurrencyDivisor(currency: string): number {
  switch (currency.toLowerCase()) {
    case "jpy":
    case "krw":
      return 1;
    default:
      return 100;
  }
}

/** Returns the number of decimal places for the currency. */
export function getCurrencyDecimalPlaces(currency: string): number {
  switch (currency.toLowerCase()) {
    case "jpy":
    case "krw":
      return 0;
    default:
      return 2;
  }
}

/** Converts a currency code to its symbol. */
export function toCurrencySymbol(currency: string): string {
  const symbols: Record<string, string> = {
    usd: "$",
    eur: "€",
    gbp: "£",
    jpy: "¥",
    aud: "A$",
    cad: "C$",
    cny: "¥",
    inr: "₹",
    rub: "₽",
    brl: "R$",
    mxn: "MX$",
    sgd: "S$",
    hkd: "HK$",
    krw: "₩",
    sek: "kr",
    dkk: "kr",
    nok: "kr",
  };
  return symbols[currency] || currency.toUpperCase();
}
