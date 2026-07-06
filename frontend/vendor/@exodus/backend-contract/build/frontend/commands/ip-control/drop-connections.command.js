"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DropConnectionsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DropConnectionsCommand;
(function (DropConnectionsCommand) {
    DropConnectionsCommand.url = api_1.REST_API.IP_CONTROL.DROP_CONNECTIONS;
    DropConnectionsCommand.TSQ_url = DropConnectionsCommand.url;
    DropConnectionsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.IP_CONTROL_ROUTES.DROP_CONNECTIONS, 'post', 'Drop Connections for Users or IPs', { scope: 'drop-connections', kind: 'write' });
    DropConnectionsCommand.DropBySchema = zod_1.z.discriminatedUnion('by', [
        zod_1.z
            .object({
            by: zod_1.z.literal('userUuids'),
            userUuids: zod_1.z.array(zod_1.z.string().uuid()).min(1),
        })
            .describe('Drop by user UUIDs'),
        zod_1.z
            .object({
            by: zod_1.z.literal('ipAddresses'),
            ipAddresses: zod_1.z.array(zod_1.z.string().ip()).min(1),
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
            nodeUuids: zod_1.z.array(zod_1.z.string().uuid()).min(1),
        })
            .describe('Target specific nodes'),
    ]);
    DropConnectionsCommand.RequestSchema = zod_1.z.object({
        dropBy: DropConnectionsCommand.DropBySchema,
        targetNodes: DropConnectionsCommand.TargetNodesSchema,
    });
    DropConnectionsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(DropConnectionsCommand || (exports.DropConnectionsCommand = DropConnectionsCommand = {}));
