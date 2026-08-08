"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TorrentBlockerReportSchema = void 0;
const zod_1 = require("zod");
const extended_users_schema_1 = require("./extended-users.schema");
const nodes_schema_1 = require("./nodes.schema");
exports.TorrentBlockerReportSchema = zod_1.z.object({
    id: zod_1.z.number(),
    userId: zod_1.z.number(),
    nodeId: zod_1.z.number(),
    user: extended_users_schema_1.ExtendedUsersSchema.pick({
        username: true,
    }),
    node: nodes_schema_1.NodesSchema.pick({
        uuid: true,
        name: true,
        countryCode: true,
    }),
    report: zod_1.z.object({
        actionReport: zod_1.z.object({
            blocked: zod_1.z.boolean(),
            ip: zod_1.z.string(),
            blockDuration: zod_1.z.number(),
            willUnblockAt: zod_1.z.iso
                .datetime({ offset: true, local: true })
                .transform((str) => new Date(str)),
            userId: zod_1.z.string(),
            processedAt: zod_1.z.iso
                .datetime({ offset: true, local: true })
                .transform((str) => new Date(str)),
        }),
        xrayReport: zod_1.z.object({
            email: zod_1.z.string().nullable(),
            level: zod_1.z.number().nullable(),
            protocol: zod_1.z.string().nullable(),
            network: zod_1.z.string(),
            source: zod_1.z.string().nullable(),
            destination: zod_1.z.string(),
            routeTarget: zod_1.z.string().nullable(),
            originalTarget: zod_1.z.string().nullable(),
            inboundTag: zod_1.z.string().nullable(),
            inboundName: zod_1.z.string().nullable(),
            inboundLocal: zod_1.z.string().nullable(),
            outboundTag: zod_1.z.string().nullable(),
            ts: zod_1.z.number(),
        }),
    }),
    createdAt: zod_1.z.iso.datetime().transform((str) => new Date(str)),
});
