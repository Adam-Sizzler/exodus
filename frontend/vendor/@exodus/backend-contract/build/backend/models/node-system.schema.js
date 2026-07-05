"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodeSystemSchema = exports.NodeSystemStatsSchema = exports.NodeSystemInfoSchema = exports.NetworkInterfaceSchema = void 0;
const zod_1 = require("zod");
exports.NetworkInterfaceSchema = zod_1.z.object({
    interface: zod_1.z.string(),
    rxBytesPerSec: zod_1.z.number(),
    txBytesPerSec: zod_1.z.number(),
    rxTotal: zod_1.z.number(),
    txTotal: zod_1.z.number(),
});
exports.NodeSystemInfoSchema = zod_1.z.object({
    arch: zod_1.z.string(),
    cpus: zod_1.z.number().int(),
    cpuModel: zod_1.z.string(),
    memoryTotal: zod_1.z.number(),
    hostname: zod_1.z.string(),
    platform: zod_1.z.string(),
    release: zod_1.z.string(),
    type: zod_1.z.string(),
    version: zod_1.z.string(),
    networkInterfaces: zod_1.z.array(zod_1.z.string()),
});
exports.NodeSystemStatsSchema = zod_1.z.object({
    memoryFree: zod_1.z.number(),
    memoryUsed: zod_1.z.number(),
    uptime: zod_1.z.number(),
    loadAvg: zod_1.z.array(zod_1.z.number()),
    interface: zod_1.z.nullable(exports.NetworkInterfaceSchema),
});
exports.NodeSystemSchema = zod_1.z.object({
    info: exports.NodeSystemInfoSchema,
    stats: exports.NodeSystemStatsSchema,
});
