"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodeResponseSchema = void 0;
const zod_1 = require("zod");
const models_1 = require("../../models");
exports.NodeResponseSchema = zod_1.z.object({
    response: models_1.NodesSchema,
});
