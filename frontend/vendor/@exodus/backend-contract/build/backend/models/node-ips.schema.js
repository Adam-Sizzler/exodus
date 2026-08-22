"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodeIpsSchema = exports.NodeIpSchema = void 0;
const zod_1 = require("zod");
const ip_statuses_1 = require("../constants/nodes/ip-statuses");
exports.NodeIpSchema = zod_1.z.object({
    ip: zod_1.z.union([zod_1.z.ipv4(), zod_1.z.ipv6()]),
    status: zod_1.z.enum(ip_statuses_1.NODE_IP_STATUSES),
});
exports.NodeIpsSchema = zod_1.z.array(exports.NodeIpSchema).max(64);
