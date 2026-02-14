export function toSnakeCase(str: string) {
  return str
    .trim()
    .replace(/([a-z0-9])([A-Z])/g, "$1_$2")
    .replace(/[\s-]+/g, "_")
    .replace(/__+/g, "_")
    .replace(/[/\\]+/g, "_")
    .toLowerCase();
}
