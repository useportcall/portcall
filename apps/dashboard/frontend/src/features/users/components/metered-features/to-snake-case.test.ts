import { describe, expect, it } from "vitest";
import { toSnakeCase } from "./to-snake-case";

describe("toSnakeCase", () => {
  it("normalizes whitespace and casing", () => {
    expect(toSnakeCase("My Feature Name")).toBe("my_feature_name");
  });

  it("normalizes separators", () => {
    expect(toSnakeCase("feature/path-name")).toBe("feature_path_name");
  });
});
