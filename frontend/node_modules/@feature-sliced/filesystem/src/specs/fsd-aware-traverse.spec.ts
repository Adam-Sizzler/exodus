import { join } from "node:path";
import { test, expect, describe } from "vitest";

import {
  getAllSlices,
  getIndexes,
  getSlices,
  isSlice,
  isSliced,
  isCrossImportPublicApi,
  type Folder,
  type File,
} from "../index.js";
import { joinFromRoot, parseIntoFolder } from "./prepare-test.js";

test("getSlices", () => {
  const rootFolder = parseIntoFolder(`
    📂 entities
      📂 user
        📂 ui
        📄 index.ts
      📂 posts
        📂 ui
        📄 index.ts
    📂 features
      📂 comments
        📂 ui
        📄 index.ts
    📂 pages
      📂 editor
        📂 ui
      📂 settings
        📂 notifications
          📂 ui
          📄 index.ts
        📂 profile
          📂 ui
          📄 index.ts
  `);

  expect(getSlices(rootFolder.children[0] as Folder)).toEqual({
    user: (rootFolder.children[0] as Folder).children[0],
    posts: (rootFolder.children[0] as Folder).children[1],
  });
  expect(getSlices(rootFolder.children[1] as Folder)).toEqual({
    comments: (rootFolder.children[1] as Folder).children[0],
  });
  expect(getSlices(rootFolder.children[2] as Folder)).toEqual({
    editor: (rootFolder.children[2] as Folder).children[0],
    [join("settings", "notifications")]: (
      (rootFolder.children[2] as Folder).children[1] as Folder
    ).children[0],
    [join("settings", "profile")]: (
      (rootFolder.children[2] as Folder).children[1] as Folder
    ).children[1],
  });
});

test("getAllSlices", () => {
  const rootFolder = parseIntoFolder(`
    📂 shared
      📂 ui
        📄 index.ts
      📂 i18n
        📄 index.ts
    📂 entities
      📂 user
        📂 ui
        📂 model
        📄 index.ts
    📂 features
      📂 comments
        📂 ui
        📄 index.ts
    📂 pages
      📂 home
        📂 ui
        📄 index.ts
  `);

  const allSlices = getAllSlices(rootFolder);
  expect(Object.keys(allSlices).sort((a, b) => a.localeCompare(b))).toEqual([
    "comments",
    "home",
    "user",
  ]);

  expect(allSlices.user).toEqual({
    ...(rootFolder.children[1] as Folder).children[0],
    layerName: "entities",
  });
  expect(allSlices.comments).toEqual({
    ...(rootFolder.children[2] as Folder).children[0],
    layerName: "features",
  });
  expect(allSlices.home).toEqual({
    ...(rootFolder.children[3] as Folder).children[0],
    layerName: "pages",
  });
});

test("isSliced", () => {
  expect(isSliced("shared")).toBe(false);
  expect(isSliced("app")).toBe(false);
  expect(isSliced("entities")).toBe(true);
  expect(isSliced("features")).toBe(true);
  expect(isSliced("pages")).toBe(true);
  expect(isSliced("widgets")).toBe(true);

  expect(
    isSliced({
      type: "folder",
      path: joinFromRoot("project", "src", "shared"),
      children: [],
    }),
  ).toBe(false);
  expect(
    isSliced({
      type: "folder",
      path: joinFromRoot("project", "src", "entities"),
      children: [],
    }),
  ).toBe(true);
});

describe("getIndexes", () => {
  test("basic functionality", () => {
    const indexFile: File = {
      type: "file",
      path: joinFromRoot("project", "src", "shared", "index.ts"),
    };
    const fileSegment: File = {
      type: "file",
      path: joinFromRoot("project", "src", "entities", "user", "ui.ts"),
    };
    const folderSegment = parseIntoFolder(
      `
      📄 Avatar.tsx
      📄 User.tsx
      📄 index.ts
      `,
      joinFromRoot("project", "src", "entities", "user", "ui"),
    );
    expect(getIndexes(indexFile)).toEqual([indexFile]);
    expect(getIndexes(fileSegment)).toEqual([fileSegment]);
    expect(getIndexes(folderSegment)).toEqual([
      {
        type: "file",
        path: joinFromRoot(
          "project",
          "src",
          "entities",
          "user",
          "ui",
          "index.ts",
        ),
      },
    ]);
  });

  test("recognizes index.server.js as index file", () => {
    const indexServerFile: File = {
      type: "file",
      path: joinFromRoot("project", "src", "shared", "index.server.js"),
    };
    const nonIndexFile: File = {
      type: "file",
      path: joinFromRoot("project", "src", "entities", "user", "ui.ts"),
    };
    const folderSegment = parseIntoFolder(
      `
      📄 Avatar.tsx
      📄 User.tsx
      📄 index.server.js
      `,
      joinFromRoot("project", "src", "entities", "user", "ui"),
    );

    expect(getIndexes(indexServerFile)).toEqual([indexServerFile]);
    expect(getIndexes(nonIndexFile)).toEqual([nonIndexFile]);
    expect(getIndexes(folderSegment)).toEqual([
      {
        type: "file",
        path: joinFromRoot(
          "project",
          "src",
          "entities",
          "user",
          "ui",
          "index.server.js",
        ),
      },
    ]);
  });

  test("handles multiple index files", () => {
    const folderSegment = parseIntoFolder(
      `
      📄 Avatar.tsx
      📄 User.tsx
      📄 index.client.js
      📄 index.server.js
      `,
      joinFromRoot("project", "src", "entities", "user", "ui"),
    );

    expect(getIndexes(folderSegment)).toEqual([
      {
        type: "file",
        path: joinFromRoot(
          "project",
          "src",
          "entities",
          "user",
          "ui",
          "index.client.js",
        ),
      },
      {
        type: "file",
        path: joinFromRoot(
          "project",
          "src",
          "entities",
          "user",
          "ui",
          "index.server.js",
        ),
      },
    ]);
  });
});

test("isCrossImportPublicApi", () => {
  const file: File = {
    path: joinFromRoot(
      "project",
      "src",
      "entities",
      "user",
      "@x",
      "product.ts",
    ),
    type: "file",
  };

  expect(
    isCrossImportPublicApi(file, {
      inSlice: "user",
      forSlice: "product",
      layerPath: joinFromRoot("project", "src", "entities"),
    }),
  ).toBe(true);
  expect(
    isCrossImportPublicApi(file, {
      inSlice: "product",
      forSlice: "user",
      layerPath: joinFromRoot("project", "src", "entities"),
    }),
  ).toBe(false);
});

test("isCrossImportPublicApi with index files", () => {
  const indexFile: File = {
    path: joinFromRoot(
      "project",
      "src",
      "entities",
      "user",
      "@x",
      "product",
      "index.ts",
    ),
    type: "file",
  };

  expect(
    isCrossImportPublicApi(indexFile, {
      inSlice: "user",
      forSlice: "product",
      layerPath: joinFromRoot("project", "src", "entities"),
    }),
  ).toBe(true);
  expect(
    isCrossImportPublicApi(indexFile, {
      inSlice: "product",
      forSlice: "user",
      layerPath: joinFromRoot("project", "src", "entities"),
    }),
  ).toBe(false);
});

test("isSlice", () => {
  const sliceFolder = parseIntoFolder(
    `
    📂 ui
      📄 Avatar.tsx
      📄 User.tsx
    📄 index.ts
    `,
    joinFromRoot("project", "src", "entities", "user"),
  );
  expect(isSlice(sliceFolder)).toBe(true);

  const notSliceFolder = parseIntoFolder(
    `
    📂 shared
    📂 pages
    📂 app
    `,
    joinFromRoot("project", "src"),
  );
  expect(isSlice(notSliceFolder)).toBe(false);
});
