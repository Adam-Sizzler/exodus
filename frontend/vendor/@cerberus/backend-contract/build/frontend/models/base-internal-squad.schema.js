"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BaseInternalSquadSchema = void 0;
const zod_1 = require("zod");
exports.BaseInternalSquadSchema = zod_1.z.object({
    uuid: zod_1.z.string().uuid(),
    name: zod_1.z.string(),
});
