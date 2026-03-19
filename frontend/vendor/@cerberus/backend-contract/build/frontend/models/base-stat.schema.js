"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BaseStatSchema = void 0;
const zod_1 = require("zod");
exports.BaseStatSchema = zod_1.z.object({
    current: zod_1.z.string(),
    previous: zod_1.z.string(),
    difference: zod_1.z.string(),
});
