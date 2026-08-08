"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.PluginExecutorCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var PluginExecutorCommand;
(function (PluginExecutorCommand) {
    PluginExecutorCommand.url = api_1.REST_API.NODE_PLUGINS.EXECUTOR;
    PluginExecutorCommand.TSQ_url = PluginExecutorCommand.url;
    PluginExecutorCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_PLUGINS_ROUTES.EXECUTOR, 'post', 'Execute command on node plugins', { scope: 'executor', kind: 'write' });
    PluginExecutorCommand.CommandSchema = zod_1.z.discriminatedUnion('command', [
        zod_1.z
            .object({
            command: zod_1.z.literal('blockIps'),
            ips: zod_1.z
                .array(zod_1.z.object({
                ip: zod_1.z.union([zod_1.z.ipv4(), zod_1.z.ipv6()]),
                timeout: zod_1.z.number(),
            }))
                .min(1),
        })
            .describe('Block IPs'),
        zod_1.z
            .object({
            command: zod_1.z.literal('unblockIps'),
            ips: zod_1.z.array(zod_1.z.union([zod_1.z.ipv4(), zod_1.z.ipv6()])).min(1),
        })
            .describe('Unblock IPs'),
        zod_1.z
            .object({
            command: zod_1.z.literal('recreateTables'),
        })
            .describe('Recreate tables'),
    ]);
    PluginExecutorCommand.TargetNodesSchema = zod_1.z.discriminatedUnion('target', [
        zod_1.z
            .object({
            target: zod_1.z.literal('allNodes'),
        })
            .describe('Target all connected nodes'),
        zod_1.z
            .object({
            target: zod_1.z.literal('specificNodes'),
            nodeUuids: zod_1.z.array(zod_1.z.uuid()).min(1),
        })
            .describe('Target specific nodes'),
    ]);
    PluginExecutorCommand.RequestBodySchema = zod_1.z.object({
        command: PluginExecutorCommand.CommandSchema,
        targetNodes: PluginExecutorCommand.TargetNodesSchema,
    });

})(PluginExecutorCommand || (exports.PluginExecutorCommand = PluginExecutorCommand = {}));
