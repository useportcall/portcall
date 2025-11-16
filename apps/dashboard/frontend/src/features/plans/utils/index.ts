export function currencyCodetoSymbol(code: string): string {
  switch (code.toLowerCase()) {
    case "usd":
      return "$";
    case "eur":
      return "€";
    case "jpy":
      return "¥";
    case "gbp":
      return "£";
    case "aud":
      return "A$";
    case "cad":
      return "C$";
    case "krw":
      return "₩";
    case "inr":
      return "₹";
    case "sgd":
      return "S$";
    case "cny":
      return "¥";
    case "try":
      return "₺";
    default:
      return code;
  }
}

export function toSnakeCase(str: string): string {
  return (
    str
      .replace(/([a-z0-9])([A-Z])/g, "$1_$2") // camelCase or PascalCase → snake_case
      .replace(/[\s\-]+/g, "_") // spaces and hyphens → underscores
      .replace(/__+/g, "_") // collapse multiple underscores
      // replace / with underscore
      .replace(/[\/\\]+/g, "_")
      .trim()
      .toLowerCase()
  );
}

export function toIntervalText(interval: string) {
  switch (interval) {
    case "month":
      return "Monthly";
    case "year":
      return "Yearly";
    case "no_reset":
      return "One-off";
    default:
      return interval;
  }
}
