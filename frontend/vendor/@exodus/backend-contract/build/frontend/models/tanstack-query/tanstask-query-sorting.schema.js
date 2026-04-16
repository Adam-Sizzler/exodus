"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TanstackQuerySortingSchema = void 0;
const zod_1 = require("zod");
exports.TanstackQuerySortingSchema = zod_1.z.object({
    id: zod_1.z.string(),
    desc: zod_1.z.boolean(),
});
