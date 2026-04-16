"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkNodesProfileModificationCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const api_1 = require("../../../api");
var BulkNodesProfileModificationCommand;
(function (BulkNodesProfileModificationCommand) {
    BulkNodesProfileModificationCommand.url = api_1.REST_API.NODES.BULK_ACTIONS.PROFILE_MODIFICATION;
    BulkNodesProfileModificationCommand.TSQ_url = BulkNodesProfileModificationCommand.url;
    BulkNodesProfileModificationCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODES_ROUTES.BULK_ACTIONS.PROFILE_MODIFICATION, 'post', 'Modify Inbounds & Profile for many nodes');
    BulkNodesProfileModificationCommand.RequestSchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.string().uuid()).min(1, 'Must be at least 1 Node UUID'),
        configProfile: zod_1.z.object({
            activeConfigProfileUuid: zod_1.z.string().uuid(),
            activeInbounds: zod_1.z
                .array(zod_1.z.string().uuid(), {
                invalid_type_error: 'Must be an array of UUIDs',
            })
                .min(1, 'Must be at least 1 inbound UUID'),
        }),
    });
    BulkNodesProfileModificationCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(BulkNodesProfileModificationCommand || (exports.BulkNodesProfileModificationCommand = BulkNodesProfileModificationCommand = {}));
