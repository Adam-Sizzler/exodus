import { expect, test } from "vitest";
import { resolveImport } from "../index.js";
import { join } from "node:path";

test("Basic", () => {
  const tsCompilerOptions = {
    moduleResolution: "Bundler" as const,
    baseUrl: ".",
    paths: {
      "~/*": ["./src/*"],
    },
  };

  function fileExists(path: string) {
    return path === "src/shared/ui/index.ts";
  }

  expect(
    resolveImport(
      "~/shared/ui",
      "src/pages/home/ui/HomePage.tsx",
      tsCompilerOptions,
      fileExists,
    ),
  ).toBe(join("src", "shared", "ui", "index.ts"));
});

test("With deprecated moduleResolution: node", () => {
  const tsCompilerOptions = {
    moduleResolution: "node" as const,
    baseUrl: ".",
    paths: {
      "~/*": ["./src/*"],
    },
  };

  function fileExists(path: string) {
    return path === "src/shared/ui/index.ts";
  }

  expect(
    resolveImport(
      "~/shared/ui",
      "src/pages/home/ui/HomePage.tsx",
      tsCompilerOptions,
      fileExists,
    ),
  ).toBe(join("src", "shared", "ui", "index.ts"));
});

test("Alias to absolute paths (for whatever reason)", () => {
  const tsCompilerOptions = {
    paths: {
      "@/*": ["/src/*"],
    },
  };

  function fileExists(path: string) {
    return path === "/src/shared/ui/Button.ts";
  }

  expect(
    resolveImport(
      "@/shared/ui/Button",
      "src/pages/home/ui/HomePage.tsx",
      tsCompilerOptions,
      fileExists,
    ),
  ).toBe(join("/", "src", "shared", "ui", "Button.ts"));
});
