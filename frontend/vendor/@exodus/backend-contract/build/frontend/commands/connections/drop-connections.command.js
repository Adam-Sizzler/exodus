"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DropConnectionsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DropConnectionsCommand;
(function (DropConnectionsCommand) {
    DropConnectionsCommand.url = api_1.REST_API.CONNECTIONS.DROP_CONNECTIONS;
    DropConnectionsCommand.TSQ_url = DropConnectionsCommand.url;
    DropConnectionsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONNECTIONS_ROUTES.DROP_CONNECTIONS, 'post', 'Drop Connections for Users or IPs', { scope: 'drop', kind: 'write' });
    DropConnectionsCommand.DropBySchema = zod_1.z.discriminatedUnion('by', [
        zod_1.z
            .object({
            by: zod_1.z.literal('userIds'),
            userIds: zod_1.z.array(zod_1.z.number()).min(1),
        })
            .describe('Drop by user IDs'),
        zod_1.z
            .object({
            by: zod_1.z.literal('ipAddresses'),
            ipAddresses: zod_1.z.array(zod_1.z.union([zod_1.z.ipv4(), zod_1.z.ipv6()])).min(1),
        })
            .describe('Drop by IP addresses'),
    ]);
    DropConnectionsCommand.TargetNodesSchema = zod_1.z.discriminatedUnion('target', [
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
    DropConnectionsCommand.RequestBodySchema = zod_1.z.object({
        dropBy: DropConnectionsCommand.DropBySchema,
        targetNodes: DropConnectionsCommand.TargetNodesSchema,
    });

})(DropConnectionsCommand || (exports.DropConnectionsCommand = DropConnectionsCommand = {}));
