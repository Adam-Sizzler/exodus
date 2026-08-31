"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SshServerMessageSchema = exports.SshClientMessageSchema = exports.SshClientErrorSchema = exports.SshHostKeyReplySchema = exports.SshSignReplySchema = exports.SshIdentitiesReplySchema = exports.SshResizeMessageSchema = exports.SshOpenMessageSchema = exports.SSH_TERMINAL_WS_PROTOCOL = exports.SSH_TERMINAL_WS_PATH = void 0;
const zod_1 = require("zod");
exports.SSH_TERMINAL_WS_PATH = '/api/node-ssh/ws';
exports.SSH_TERMINAL_WS_PROTOCOL = 'ex';
const requestId = zod_1.z.number().int().nonnegative();
exports.SshOpenMessageSchema = zod_1.z.object({
    t: zod_1.z.literal('open'),
    host: zod_1.z.string().min(1).max(253),
    port: zod_1.z.number().int().min(1).max(65535),
    username: zod_1.z.string().min(1).max(64),
    cols: zod_1.z.number().int().min(1).max(1000),
    rows: zod_1.z.number().int().min(1).max(1000),
});
exports.SshResizeMessageSchema = zod_1.z.object({
    t: zod_1.z.literal('resize'),
    cols: zod_1.z.number().int().min(1).max(1000),
    rows: zod_1.z.number().int().min(1).max(1000),
});
exports.SshIdentitiesReplySchema = zod_1.z.object({
    t: zod_1.z.literal('identities'),
    id: requestId,
    keys: zod_1.z.array(zod_1.z.string().min(1).max(4096)).max(16),
});
exports.SshSignReplySchema = zod_1.z.object({
    t: zod_1.z.literal('sign'),
    id: requestId,
    signature: zod_1.z.base64().max(2048),
});
exports.SshHostKeyReplySchema = zod_1.z.object({
    t: zod_1.z.literal('hostkey'),
    id: requestId,
    accept: zod_1.z.boolean(),
});
exports.SshClientErrorSchema = zod_1.z.object({
    t: zod_1.z.literal('error'),
    id: requestId,
    message: zod_1.z.string().max(500),
});
/** Browser -> panel. */
exports.SshClientMessageSchema = zod_1.z.discriminatedUnion('t', [
    exports.SshOpenMessageSchema,
    exports.SshResizeMessageSchema,
    exports.SshIdentitiesReplySchema,
    exports.SshSignReplySchema,
    exports.SshHostKeyReplySchema,
    exports.SshClientErrorSchema,
]);
/** Panel -> browser. */
exports.SshServerMessageSchema = zod_1.z.discriminatedUnion('t', [
    zod_1.z.object({ t: zod_1.z.literal('agent-identities'), id: requestId }),
    zod_1.z.object({
        t: zod_1.z.literal('agent-sign'),
        id: requestId,
        publicKey: zod_1.z.base64().max(4096),
        data: zod_1.z.base64().max(8192),
        hash: zod_1.z.string().max(16).nullable(),
    }),
    zod_1.z.object({
        t: zod_1.z.literal('hostkey'),
        id: requestId,
        algo: zod_1.z.string().max(64),
        fingerprint: zod_1.z.string().max(128),
        key: zod_1.z.base64().max(4096),
    }),
    zod_1.z.object({ t: zod_1.z.literal('ready') }),
    zod_1.z.object({
        t: zod_1.z.literal('exit'),
        code: zod_1.z.number().int().nullable(),
        signal: zod_1.z.string().max(32).nullable(),
    }),
    zod_1.z.object({ t: zod_1.z.literal('error'), message: zod_1.z.string().max(500) }),
]);
