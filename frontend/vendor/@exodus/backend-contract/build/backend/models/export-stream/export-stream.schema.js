"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExodusNodeConnectionsStreamMessageSchema = exports.NodeConnectionUserSchema = exports.NODE_CONNECTIONS_STREAM_MESSAGE_VERSION = exports.ExodusSubscriptionRequestStreamMessageSchema = exports.SUBSCRIPTION_REQUEST_STREAM_MESSAGE_VERSION = exports.ExodusUserUsageStreamMessageSchema = exports.UserUsageStreamRecordSchema = exports.USER_USAGE_STREAM_MESSAGE_VERSION = void 0;
const zod_1 = require("zod");
const USER_USAGE_STREAM_KEY = 'ioraw:export:user_usage';
const SUBSCRIPTION_REQUESTS_STREAM_KEY = 'ioraw:export:subscription_requests';
const NODE_CONNECTIONS_STREAM_KEY = 'ioraw:export:node_connections';
exports.USER_USAGE_STREAM_MESSAGE_VERSION = '1';
exports.UserUsageStreamRecordSchema = (zod_1.z || zod_1.default.z || zod_1.default).object({
    userId: (zod_1.z || zod_1.default.z || zod_1.default).string().regex(/^\d+$/).describe('User ID (bigint as string).'),
    totalBytes: zod_1.default
        .string()
        .regex(/^\d+$/)
        .describe('Traffic consumed by the user in this batch, bytes (bigint as string).'),
});
exports.ExodusUserUsageStreamMessageSchema = zod_1.default
    .object({
    v: (zod_1.z || zod_1.default.z || zod_1.default).literal(exports.USER_USAGE_STREAM_MESSAGE_VERSION).describe('Message schema version.'),
    nodeId: zod_1.default
        .string()
        .regex(/^\d+$/)
        .describe('Node ID the batch belongs to (bigint as string).'),
    ts: (zod_1.z || zod_1.default.z || zod_1.default).iso
        .datetime()
        .transform((str) => new Date(str))
        .describe('Time the batch was exported (ISO 8601, UTC).'),
    records: zod_1.default
        .string()
        .regex(/^\d+:\d+(;\d+:\d+)*$/)
        .transform((raw) => raw.split(';').map((pair) => {
        const [userId, totalBytes] = pair.split(':');
        return { userId, totalBytes };
    }))
        .pipe((zod_1.z || zod_1.default.z || zod_1.default).array(exports.UserUsageStreamRecordSchema))
        .describe('User traffic deltas: "userId:totalBytes" pairs separated by ";".'),
})
    .meta({
    description: `A single message of the "${USER_USAGE_STREAM_KEY}" Redis Stream (EXPORT_TO_STREAM_ENABLED).`,
});
exports.SUBSCRIPTION_REQUEST_STREAM_MESSAGE_VERSION = '1';
exports.ExodusSubscriptionRequestStreamMessageSchema = zod_1.default
    .object({
    v: zod_1.default
        .literal(exports.SUBSCRIPTION_REQUEST_STREAM_MESSAGE_VERSION)
        .describe('Message schema version.'),
    userId: zod_1.default
        .string()
        .regex(/^\d+$/)
        .describe('ID of the user who requested the subscription (bigint as string).'),
    requestAt: (zod_1.z || zod_1.default.z || zod_1.default).iso
        .datetime()
        .transform((str) => new Date(str))
        .describe('Time of the subscription request (ISO 8601, UTC).'),
    requestIp: (zod_1.z || zod_1.default.z || zod_1.default).string().optional().describe('Client IP address, omitted if unknown.'),
    userAgent: (zod_1.z || zod_1.default.z || zod_1.default).string().optional().describe('Client User-Agent, omitted if unknown.'),
})
    .meta({
    description: `A single message of the "${SUBSCRIPTION_REQUESTS_STREAM_KEY}" Redis Stream (EXPORT_TO_STREAM_ENABLED).`,
});
exports.NODE_CONNECTIONS_STREAM_MESSAGE_VERSION = '1';
exports.NodeConnectionUserSchema = (zod_1.z || zod_1.default.z || zod_1.default).object({
    userId: (zod_1.z || zod_1.default.z || zod_1.default).string().regex(/^\d+$/).describe('User ID (bigint as string).'),
    ips: zod_1.default
        .array((zod_1.z || zod_1.default.z || zod_1.default).object({
        ip: (zod_1.z || zod_1.default.z || zod_1.default).string().describe('Client IP address.'),
        lastSeen: (zod_1.z || zod_1.default.z || zod_1.default).iso
            .datetime()
            .transform((str) => new Date(str))
            .describe('Last time this IP was seen on the node (ISO 8601, UTC).'),
    }))
        .describe('IP addresses the user is connected from.'),
});
exports.ExodusNodeConnectionsStreamMessageSchema = zod_1.default
    .object({
    v: (zod_1.z || zod_1.default.z || zod_1.default).literal(exports.NODE_CONNECTIONS_STREAM_MESSAGE_VERSION).describe('Message schema version.'),
    nodeId: zod_1.default
        .string()
        .regex(/^\d+$/)
        .describe('Node ID the snapshot belongs to (bigint as string).'),
    ts: (zod_1.z || zod_1.default.z || zod_1.default).iso
        .datetime()
        .transform((str) => new Date(str))
        .describe('Time the snapshot was exported (ISO 8601, UTC).'),
    users: zod_1.default
        .string()
        .transform((raw, ctx) => {
        try {
            return JSON.parse(raw);
        }
        catch {
            ctx.addIssue({ code: 'custom', message: 'users must be a valid JSON array' });
            return (zod_1.z || zod_1.default.z || zod_1.default).NEVER;
        }
    })
        .pipe((zod_1.z || zod_1.default.z || zod_1.default).array(exports.NodeConnectionUserSchema))
        .describe([
        'JSON-encoded array of users connected to the node with their IPs.',
        'Example:',
        '```json',
        '[{ "userId": "42", "ips": [{ "ip": "1.2.3.4", "lastSeen": "2026-07-16T11:59:30.000Z" }] }]',
        '```',
    ].join('\n')),
})
    .meta({
    description: `A single message of the "${NODE_CONNECTIONS_STREAM_KEY}" Redis Stream (EXPORT_TO_STREAM_ENABLED). Snapshot of connections on one node; retention is time-based.`,
});
