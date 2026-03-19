"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExternalSquadHostOverridesSchema = void 0;
const hosts_schema_1 = require("../hosts.schema");
exports.ExternalSquadHostOverridesSchema = hosts_schema_1.HostsSchema.pick({
    serverDescription: true,
    vlessRouteId: true,
}).partial();
