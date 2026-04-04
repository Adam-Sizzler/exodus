"use strict";
var __defProp = Object.defineProperty;
var __defProps = Object.defineProperties;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropDescs = Object.getOwnPropertyDescriptors;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __getOwnPropSymbols = Object.getOwnPropertySymbols;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __propIsEnum = Object.prototype.propertyIsEnumerable;
var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
var __spreadValues = (a, b) => {
  for (var prop in b || (b = {}))
    if (__hasOwnProp.call(b, prop))
      __defNormalProp(a, prop, b[prop]);
  if (__getOwnPropSymbols)
    for (var prop of __getOwnPropSymbols(b)) {
      if (__propIsEnum.call(b, prop))
        __defNormalProp(a, prop, b[prop]);
    }
  return a;
};
var __spreadProps = (a, b) => __defProps(a, __getOwnPropDescs(b));
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/index.ts
var src_exports = {};
__export(src_exports, {
  arrayMove: () => arrayMove,
  arraySwap: () => arraySwap,
  move: () => move,
  swap: () => swap
});
module.exports = __toCommonJS(src_exports);

// src/move.ts
function arrayMove(array, from, to) {
  if (from === to) {
    return array;
  }
  const newArray = array.slice();
  newArray.splice(to, 0, newArray.splice(from, 1)[0]);
  return newArray;
}
function arraySwap(array, from, to) {
  if (from === to) {
    return array;
  }
  const newArray = array.slice();
  const item = newArray[from];
  newArray[from] = newArray[to];
  newArray[to] = item;
  return newArray;
}
function mutate(items, event, mutation) {
  var _a, _b;
  const { source, target, canceled } = event.operation;
  if (!source || !target || canceled) {
    if ("preventDefault" in event) event.preventDefault();
    return items;
  }
  const findIndex = (item, id) => item === id || typeof item === "object" && "id" in item && item.id === id;
  if (Array.isArray(items)) {
    const sourceIndex2 = items.findIndex((item) => findIndex(item, source.id));
    const targetIndex2 = items.findIndex((item) => findIndex(item, target.id));
    if (sourceIndex2 === -1 || targetIndex2 === -1) {
      return items;
    }
    if (!canceled && "index" in source && typeof source.index === "number") {
      const projectedSourceIndex = source.index;
      if (projectedSourceIndex !== sourceIndex2) {
        return mutation(items, sourceIndex2, projectedSourceIndex);
      }
    }
    return mutation(items, sourceIndex2, targetIndex2);
  }
  const entries = Object.entries(items);
  let sourceIndex = -1;
  let sourceParent;
  let targetIndex = -1;
  let targetParent;
  for (const [id, children] of entries) {
    if (sourceIndex === -1) {
      sourceIndex = children.findIndex((item) => findIndex(item, source.id));
      if (sourceIndex !== -1) {
        sourceParent = id;
      }
    }
    if (targetIndex === -1) {
      targetIndex = children.findIndex((item) => findIndex(item, target.id));
      if (targetIndex !== -1) {
        targetParent = id;
      }
    }
    if (sourceIndex !== -1 && targetIndex !== -1) {
      break;
    }
  }
  if (!source.manager) return items;
  const { dragOperation } = source.manager;
  const position = (_b = (_a = dragOperation.shape) == null ? void 0 : _a.current.center) != null ? _b : dragOperation.position.current;
  if (targetParent == null) {
    if (target.id in items) {
      const insertionIndex = target.shape && position.y > target.shape.center.y ? items[target.id].length : 0;
      targetParent = target.id;
      targetIndex = insertionIndex;
    }
  }
  if (sourceParent == null || targetParent == null || sourceParent === targetParent && sourceIndex === targetIndex) {
    if ("preventDefault" in event) event.preventDefault();
    return items;
  }
  if (sourceParent === targetParent) {
    return __spreadProps(__spreadValues({}, items), {
      [sourceParent]: mutation(items[sourceParent], sourceIndex, targetIndex)
    });
  }
  const isBelowTarget = target.shape && Math.round(position.y) > Math.round(target.shape.center.y);
  const modifier = isBelowTarget ? 1 : 0;
  const sourceItem = items[sourceParent][sourceIndex];
  return __spreadProps(__spreadValues({}, items), {
    [sourceParent]: [
      ...items[sourceParent].slice(0, sourceIndex),
      ...items[sourceParent].slice(sourceIndex + 1)
    ],
    [targetParent]: [
      ...items[targetParent].slice(0, targetIndex + modifier),
      sourceItem,
      ...items[targetParent].slice(targetIndex + modifier)
    ]
  });
}
function move(items, event) {
  return mutate(items, event, arrayMove);
}
function swap(items, event) {
  return mutate(items, event, arraySwap);
}
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  arrayMove,
  arraySwap,
  move,
  swap
});
