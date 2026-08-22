"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HostResponseSchema = void 0;
const zod_1 = require("zod");
const hosts_schema_1 = require("../../models/hosts.schema");
exports.HostResponseSchema = zod_1.z.object({
    response: hosts_schema_1.HostsSchema,
});
