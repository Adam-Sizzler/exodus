"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.InfraBillingAvailableNodeSchema = void 0;
const nodes_schema_1 = require("./nodes.schema");
exports.InfraBillingAvailableNodeSchema = nodes_schema_1.NodesSchema.pick({
    uuid: true,
    name: true,
    countryCode: true,
});
