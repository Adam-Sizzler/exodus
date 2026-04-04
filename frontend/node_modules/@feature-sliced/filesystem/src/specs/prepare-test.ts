import { join } from "node:path";
import type { Folder, File } from "../definitions.js";

/** Parse a multi-line indented string with emojis for files and folders into an FSD root. */
export function parseIntoFolder(
  fsMarkup: string,
  basePath = joinFromRoot(),
): Folder {
  function parseFolder(lines: Array<string>, path: string): Folder {
    const children: Array<Folder | File> = [];

    lines.forEach((line, index) => {
      if (line.startsWith("📂 ")) {
        let nestedLines = lines.slice(index + 1);
        const nextIndex = nestedLines.findIndex(
          (line) => !line.startsWith("  "),
        );
        nestedLines = nestedLines.slice(
          0,
          nextIndex === -1 ? nestedLines.length : nextIndex,
        );
        const folder = parseFolder(
          nestedLines.map((line) => line.slice("  ".length)),
          join(path, line.slice("📂 ".length)),
        );
        children.push(folder);
      } else if (line.startsWith("📄 ")) {
        children.push({
          type: "file",
          path: join(path, line.slice("📄 ".length)),
        });
      }
    });

    return { type: "folder", path, children };
  }

  const lines = fsMarkup
    .split("\n")
    .filter(Boolean)
    .map((line, _i, lines) => line.slice(lines[0].search(/\S/)))
    .filter(Boolean);

  return parseFolder(lines, basePath);
}

export function joinFromRoot(...segments: Array<string>) {
  return join("/", ...segments);
}

if (import.meta.vitest) {
  const { test, expect } = import.meta.vitest;

  test("parseIntoFolder", () => {
    const root = parseIntoFolder(`
      📂 entities
        📂 users
          📂 ui
          📄 index.ts
        📂 posts
          📂 ui
          📄 index.ts
      📂 shared
        📂 ui
          📄 index.ts
          📄 Button.tsx
    `);

    expect(root).toEqual({
      type: "folder",
      path: joinFromRoot(),
      children: [
        {
          type: "folder",
          path: joinFromRoot("entities"),
          children: [
            {
              type: "folder",
              path: joinFromRoot("entities", "users"),
              children: [
                {
                  type: "folder",
                  path: joinFromRoot("entities", "users", "ui"),
                  children: [],
                },
                {
                  type: "file",
                  path: joinFromRoot("entities", "users", "index.ts"),
                },
              ],
            },
            {
              type: "folder",
              path: joinFromRoot("entities", "posts"),
              children: [
                {
                  type: "folder",
                  path: joinFromRoot("entities", "posts", "ui"),
                  children: [],
                },
                {
                  type: "file",
                  path: joinFromRoot("entities", "posts", "index.ts"),
                },
              ],
            },
          ],
        },
        {
          type: "folder",
          path: joinFromRoot("shared"),
          children: [
            {
              type: "folder",
              path: joinFromRoot("shared", "ui"),
              children: [
                {
                  type: "file",
                  path: joinFromRoot("shared", "ui", "index.ts"),
                },
                {
                  type: "file",
                  path: joinFromRoot("shared", "ui", "Button.tsx"),
                },
              ],
            },
          ],
        },
      ],
    });
  });
}
